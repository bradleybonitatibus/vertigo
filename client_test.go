package vertigo

import (
	"context"
	"os"
	"testing"
)

func TestNewVertigoClient(t *testing.T) {
	ctx := context.Background()
	cfg := &Config{
		ProjectID:        os.Getenv("GOOGLE_PROJECT_ID"),
		Region:           nane,
		FeatureStoreName: os.Getenv("VERTIGO_FEATURESTORE_NAME"),
	}

	_, err := NewVertigoClient(ctx, cfg)
	if err != nil {
		t.Error(err)
	}
}
