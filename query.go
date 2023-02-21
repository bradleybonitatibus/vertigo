package vertigo

import (
	"fmt"
	"cloud.google.com/go/aiplatform/apiv1beta1/aiplatformpb"
)

// Query represents a query to the Vertex AI Online Feature Store API for
// getting an Entity's Feature Values.
type Query struct {
	EntityType string
	EntityID   string
	Features   []string
}

// BuildRequest translates the Query struct into an AI Platform ReadFeatureValuesRequest, which is submitted
// to the Vertex AI Online Feature Store API to retrieve the Feature Values for an entity.
func (q *Query) BuildRequest(cfg *Config) *aiplatformpb.ReadFeatureValuesRequest {
	return &aiplatformpb.ReadFeatureValuesRequest{
		EntityType:      makeVertexEntityTypePath(cfg, q.EntityType),
		EntityId:        q.EntityID,
		FeatureSelector: &aiplatformpb.FeatureSelector{IdMatcher: &aiplatformpb.IdMatcher{Ids: q.Features}},
	}
}

// makeVertexEntityTypePath builds the resource name for the specific entity being queried.
func makeVertexEntityTypePath(cfg *Config, entityType string) string {
	return fmt.Sprintf(
		"%v/entityTypes/%v",
		cfg.ParentPath(), entityType,
	)
}
