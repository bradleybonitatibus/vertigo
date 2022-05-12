package vertigo

import (
	"context"
	aiplatformpb "google.golang.org/genproto/googleapis/cloud/aiplatform/v1beta1"
	"os"
	"testing"
)

func TestNewClient(t *testing.T) {
	ctx := context.Background()
	cfg := &Config{
		ProjectID:        os.Getenv("GOOGLE_PROJECT_ID"),
		Region:           nane,
		FeatureStoreName: os.Getenv("VERTIGO_FEATURESTORE_NAME"),
	}

	client, err := NewClient(ctx, cfg)
	defer client.Close()
	if err != nil {
		t.Error(err)
	}

}

func TestEntity_ScanStruct(t *testing.T) {
	type testEntity struct {
		MyFeature      int64   `vertex:"my_feature"`
		AnotherFeature float64 `vertex:"another_feature"`
	}
	entity := Entity{
		header: &aiplatformpb.ReadFeatureValuesResponse_Header{
			EntityType: "my_entity",
			FeatureDescriptors: []*aiplatformpb.ReadFeatureValuesResponse_FeatureDescriptor{
				&aiplatformpb.ReadFeatureValuesResponse_FeatureDescriptor{
					Id: "my_feature",
				},
				&aiplatformpb.ReadFeatureValuesResponse_FeatureDescriptor{
					Id: "another_feature",
				},
			},
		},
		data: []*aiplatformpb.ReadFeatureValuesResponse_EntityView_Data{
			&aiplatformpb.ReadFeatureValuesResponse_EntityView_Data{
				Data: &aiplatformpb.ReadFeatureValuesResponse_EntityView_Data_Value{
					Value: &aiplatformpb.FeatureValue{
						Value: &aiplatformpb.FeatureValue_Int64Value{
							Int64Value: 25,
						},
					},
				},
			},
			&aiplatformpb.ReadFeatureValuesResponse_EntityView_Data{
				Data: &aiplatformpb.ReadFeatureValuesResponse_EntityView_Data_Value{
					Value: &aiplatformpb.FeatureValue{
						Value: &aiplatformpb.FeatureValue_DoubleValue{
							DoubleValue: 50.0,
						},
					},
				},
			},
		},
		ID: "123",
	}

	c := testEntity{}
	err := entity.ScanStruct(&c)
	if err != nil {
		t.Error(err)
	}
	if c.MyFeature != 25 || c.AnotherFeature != 50 {
		t.Errorf("Expected 25 and 50, got %v, and %v", c.MyFeature, c.AnotherFeature)
	}
}
