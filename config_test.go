package vertigo

import "testing"

func TestNewConfigBuilder(t *testing.T) {
	const nane = "northamerica-northeast1"
	const myProj = "my-project"
	const myFS = "my_featurestore"
	b := NewConfigBuilder()

	cfg := b.WithRegion(nane).
		WithProjectID(myProj).
		WithFeatureStoreName(myFS).
		Apply()

	if cfg.Region != nane ||
		cfg.ProjectID != myProj ||
		cfg.FeatureStoreName != myFS {
		t.Error("Builder failed to set cfg values correctly")
	}
}
