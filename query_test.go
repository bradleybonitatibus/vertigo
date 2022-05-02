package vertigo

import "testing"

func TestQuery_BuildRequest(t *testing.T) {
	q := &Query{
		EntityType: "test_entity",
		EntityID:   "123abc",
		Features:   []string{"*"},
	}
	cfg := &Config{
		ProjectID:        "my-project",
		Region:           nane,
		FeatureStoreName: "my_featurestore",
	}
	req := q.BuildRequest(cfg)
	if req.EntityId != "123abc" ||
		req.EntityType != "projects/my-project/locations/northamerica-northeast1/featurestores/my_featurestore/entityTypes/test_entity" ||
		len(req.FeatureSelector.IdMatcher.Ids) != 1 {
		t.Errorf("req was not built correctly: %v", req)
	}
}
