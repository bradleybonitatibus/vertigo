package vertigo

import (
	aiplatform "cloud.google.com/go/aiplatform/apiv1beta1"
	"context"
	"fmt"
	"google.golang.org/api/option"
	aiplatformpb "google.golang.org/genproto/googleapis/cloud/aiplatform/v1beta1"
)

// Client is the Vertigo client, which uses the aiplatformv1beta1 gRPC API to communicate
// with the FeaturestoreOnlineServingClient.
type Client struct {
	cfg *Config
	v   *aiplatform.FeaturestoreOnlineServingClient
}

// NewVertigoClient creates a Client using the provided Config.
func NewVertigoClient(ctx context.Context, cfg *Config) (*Client, error) {
	c, err := aiplatform.NewFeaturestoreOnlineServingClient(
		ctx,
		option.WithEndpoint(cfg.APIEndpoint()),
	)

	if err != nil {
		return nil, fmt.Errorf("aiplatform.NewFeaturestoreOnlineServingClient: %v", err)
	}

	return &Client{
		cfg: cfg,
		v:   c,
	}, nil
}

// Entity contains the header and data from the aiplatform.ReadFeatureValuesResponse to
// be used to scan the response into a user provided struct.
type Entity struct {
	header *aiplatformpb.ReadFeatureValuesResponse_Header
	data   []*aiplatformpb.ReadFeatureValuesResponse_EntityView_Data
}

// ScanStruct will parse the ReadFeatureValues response from the online serving client
// and load the features into dst.
// DST must be a pointer to a struct and have valid `vertex` tags that map to the
// feature IDs of the entity being parsed.
func (e *Entity) ScanStruct(dst interface{}) error {
	return nil
}

// GetEntity calls the Vertex AI Online Serving API and retrieves the response in the
// form of an Entity and error if one occurs.
func (c *Client) GetEntity(ctx context.Context, query *Query) (*Entity, error) {
	return nil, nil
}
