package vertigo

import (
	"errors"
	"testing"
)

const nane = "northamerica-northeast1"

func TestNewConfigBuilder(t *testing.T) {
	const myProj = "my-project"
	const myFS = "my_featurestore"
	type test struct {
		region    string
		projectID string
		fsName    string
		err       error
	}

	tests := []test{
		{
			region:    "",
			projectID: myProj,
			fsName:    myFS,
			err:       nil,
		},
		{
			region:    nane,
			projectID: "",
			fsName:    myFS,
			err:       ErrInvalidProjectID,
		},
		{
			region:    nane,
			projectID: myProj,
			fsName:    "",
			err:       ErrInvalidFeatureStoreName,
		},
	}

	for _, tc := range tests {
		cfg, err := NewConfigBuilder().WithRegion(tc.region).
			WithProjectID(tc.projectID).
			WithFeatureStoreName(tc.fsName).
			Apply()

		if errors.Is(err, tc.err) {
			continue
		}

		if (cfg.Region != tc.region && cfg.Region != DefaultRegion) || cfg.ProjectID != tc.projectID || cfg.FeatureStoreName != tc.fsName {
			t.Errorf("builder failed to set Config fields: %v", tc)
		}

	}
}
