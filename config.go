package vertigo

import (
	"errors"
	"fmt"
)

// DefaultRegion is the region
const DefaultRegion = "us-central1"

// VertexEndpoint is the URL:PORT for the AI Platform API. This is used in providing
// non-default regional endpoints for Vertex AI.
const VertexEndpoint = "aiplatform.googleapis.com:443"

var ErrInvalidProjectID = errors.New("project id is not valid")
var ErrInvalidFeatureStoreName = errors.New("feature store name is not valid")

// Config is the struct the contains configuration that is used in the Vertex AI API.
type Config struct {
	// ProjectID is the GCP Project ID your Feature Store resides in.
	ProjectID string `json:"project_id" yaml:"project_id"`

	// Region is the GCP Region (sometimes referred to as Location) the Feature Store is running in.
	// Please see `https://cloud.google.com/vertex-ai/docs/general/locations` for a list
	// of supported regions.
	Region string `json:"region" yaml:"region"`

	// FeatureStoreName is the name of the feature store.
	FeatureStoreName string `json:"feature_store_name" yaml:"feature_store_name"`
}

// ConfigBuilder provides a fluent interface for building the Vertigo Config.
// If Region is not set, it will fall back to the DefaultRegion value.
type ConfigBuilder interface {
	WithRegion(region string) ConfigBuilder
	WithProjectID(projectID string) ConfigBuilder
	WithFeatureStoreName(featureStore string) ConfigBuilder
	Apply() (*Config, error)
}

// builderFunc modifies a pointer to Config to set a value. This is to be used by the ConfigBuilder to
// lazily-evaluate the changes to the Config, and apply them during the call to the "Apply" function.
type builderFunc func(cfg *Config)

type builder struct {
	actions []builderFunc
}

// Apply applies all the changes and validates the Config fields.
// If Region is not set, it will fall back to the DefaultRegion value.
func (b *builder) Apply() (*Config, error) {
	cfg := &Config{}
	for _, a := range b.actions {
		a(cfg)
	}

	if cfg.Region == "" {
		cfg.Region = DefaultRegion
	}

	if cfg.ProjectID == "" {
		return nil, ErrInvalidProjectID
	}

	if cfg.FeatureStoreName == "" {
		return nil, ErrInvalidFeatureStoreName
	}

	return cfg, nil
}

// WithRegion sets the GCP Region for the Config.
func (b *builder) WithRegion(region string) ConfigBuilder {
	b.actions = append(b.actions, func(cfg *Config) {
		cfg.Region = region
	})
	return b
}

// WithProjectID sets the ProjectID in the Config struct.
func (b *builder) WithProjectID(projectID string) ConfigBuilder {
	b.actions = append(b.actions, func(cfg *Config) {
		cfg.ProjectID = projectID
	})
	return b
}

// WithFeatureStoreName sets the FeatureStoreName field in the Config struct.
func (b *builder) WithFeatureStoreName(featureStore string) ConfigBuilder {
	b.actions = append(b.actions, func(cfg *Config) {
		cfg.FeatureStoreName = featureStore
	})
	return b
}

// NewConfigBuilder returns a fluent API to build the Config struct using the ConfigBuilder interface.
func NewConfigBuilder() ConfigBuilder {
	return &builder{
		actions: []builderFunc{},
	}
}

// APIEndpoint is the Vertex AI API Endpoint, specific to the Region you have deployed
// your feature store in.
func (c *Config) APIEndpoint() string {
	if c.Region == "" {
		return buildEndpoint(DefaultRegion)
	}
	return buildEndpoint(c.Region)
}

// ParentPath is the resource hierarchy for the feature store that we are interacting with.
func (c *Config) ParentPath() string {
	return fmt.Sprintf(
		"projects/%v/locations/%v/featurestores/%v",
		c.ProjectID,
		c.Region,
		c.FeatureStoreName,
	)
}

// buildEndpoint returns the regional endpoint for the Vertex AI Platform API.
func buildEndpoint(region string) string {
	return fmt.Sprintf("%v-%v", region, VertexEndpoint)
}
