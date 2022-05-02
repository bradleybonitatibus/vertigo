package vertigo

import (
	"context"
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

	_, err := NewClient(ctx, cfg)
	if err != nil {
		t.Error(err)
	}
}
