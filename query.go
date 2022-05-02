package vertigo

import (
	"fmt"
	"google.golang.org/genproto/googleapis/cloud/aiplatform/v1"
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
func (q *Query) BuildRequest() *aiplatform.ReadFeatureValuesRequest {
	return &aiplatform.ReadFeatureValuesRequest{
		EntityType:      q.EntityType,
		EntityId:        q.EntityID,
		FeatureSelector: &aiplatform.FeatureSelector{IdMatcher: &aiplatform.IdMatcher{Ids: q.Features}},
	}
}

// makeVertexEntityTypePath builds the resource name for the specific entity being queried.
func makeVertexEntityTypePath(cfg *Config, entityID string) string {
	return fmt.Sprintf(
		"projects/%v/locations/%v/featurestores/%v/entityTypes/%v",
		cfg.ProjectID, cfg.Region, cfg.FeatureStoreName, entityID,
	)
}
