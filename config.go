package vertigo

import "fmt"

// DefaultRegion is the region
const DefaultRegion = "us-central1"

// VertexEndpoint is the URL:PORT for the AI Platform API. This is used in providing
// non-default regional endpoints for Vertex AI.
const VertexEndpoint = "aiplatform.googleapis.com:443"

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
type ConfigBuilder interface {
	WithRegion(region string) ConfigBuilder
	WithProjectID(projectID string) ConfigBuilder
	WithFeatureStoreName(featureStore string) ConfigBuilder
	Apply() *Config
}

type builderFunc func(cfg *Config)

type builder struct {
	actions []builderFunc
}

func (b *builder) Apply() *Config {
	cfg := &Config{}
	for _, a := range b.actions {
		a(cfg)
	}

	return cfg
}

func (b *builder) WithRegion(region string) ConfigBuilder {
	b.actions = append(b.actions, func(cfg *Config) {
		cfg.Region = region
	})
	return b
}

func (b *builder) WithProjectID(projectID string) ConfigBuilder {
	b.actions = append(b.actions, func(cfg *Config) {
		cfg.ProjectID = projectID
	})
	return b
}

func (b *builder) WithFeatureStoreName(featureStore string) ConfigBuilder {
	b.actions = append(b.actions, func(cfg *Config) {
		cfg.FeatureStoreName = featureStore
	})
	return b
}

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
