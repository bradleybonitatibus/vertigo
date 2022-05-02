package vertigo

import (
	aiplatform "cloud.google.com/go/aiplatform/apiv1beta1"
	"context"
	"errors"
	"fmt"
	"google.golang.org/api/option"
	aiplatformpb "google.golang.org/genproto/googleapis/cloud/aiplatform/v1beta1"
	"reflect"
)

// Client is the Vertigo client, which uses the aiplatformv1beta1 gRPC API to communicate
// with the FeaturestoreOnlineServingClient.
type Client struct {
	cfg *Config
	v   *aiplatform.FeaturestoreOnlineServingClient
}

// NewClient creates a Client using the provided Config.
func NewClient(ctx context.Context, cfg *Config) (*Client, error) {
	fmt.Println(cfg.APIEndpoint())
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
	ID     string
}

// ScanStruct will parse the ReadFeatureValues response from the online serving client
// and load the features into dst.
// DST must be a pointer to a struct and have valid `vertex` tags that map to the
// feature IDs of the entity being parsed.
func (e *Entity) ScanStruct(dst interface{}) error {
	if err := isStructPointer(dst); err != nil {
		return err
	}
	mapping := loadMap(dst)

	if len(e.header.FeatureDescriptors) != len(e.data) {
		return errors.New("feature descriptors do not match entity view data entries")
	}

	v := reflect.ValueOf(dst)
	for i, fd := range e.header.FeatureDescriptors {
		fv := e.data[i].GetValue()
		lookup, ok := mapping[fd.Id]
		if !ok {
			continue
		}
		var sf reflect.Value
		if isValuePointer(v) {
			sf = reflect.Indirect(v).Field(lookup.fieldIdx)
		} else {
			sf = v.Field(lookup.fieldIdx)
		}

		if fv.GetValue() == nil {
			continue
		}

		switch fv.Value.(type) {
		case *aiplatformpb.FeatureValue_BoolValue:
			bv := fv.GetBoolValue()
			if isValuePointer(v) {
				sf.Set(reflect.ValueOf(&bv))
			} else {
				sf.SetBool(bv)
			}

		case *aiplatformpb.FeatureValue_BoolArrayValue:
			values := fv.GetBoolArrayValue().Values
			sf.Set(reflect.MakeSlice(reflect.TypeOf(values), len(values), len(values)))

		case *aiplatformpb.FeatureValue_Int64Value:
			iv := fv.GetInt64Value()
			if isValuePointer(v) {
				sf.Set(reflect.ValueOf(&iv))
			} else {
				sf.SetInt(iv)
			}

		case *aiplatformpb.FeatureValue_Int64ArrayValue:
			setSlice(v, fv.GetInt64ArrayValue().Values)

		case *aiplatformpb.FeatureValue_DoubleValue:
			dv := fv.GetDoubleValue()
			if isValuePointer(v) {
				sf.Set(reflect.ValueOf(&dv))
			} else {
				sf.SetFloat(dv)
			}

		case *aiplatformpb.FeatureValue_DoubleArrayValue:
			setSlice(v, fv.GetDoubleArrayValue().Values)

		case *aiplatformpb.FeatureValue_StringValue:
			stringV := fv.GetStringValue()
			if isValuePointer(v) {
				sf.Set(reflect.ValueOf(&stringV))
			} else {
				sf.SetString(stringV)
			}

		case *aiplatformpb.FeatureValue_StringArrayValue:
			setSlice(v, fv.GetStringArrayValue().Values)

		case *aiplatformpb.FeatureValue_BytesValue:
			sf.SetBytes(fv.GetBytesValue())
		}
	}
	return nil
}

// GetEntity calls the Vertex AI Online Serving API and retrieves the response in the
// form of an Entity and error if one occurs.
func (c *Client) GetEntity(ctx context.Context, query *Query) (*Entity, error) {
	res, err := c.v.ReadFeatureValues(ctx, query.BuildRequest(c.cfg))
	if err != nil {
		return nil, err
	}
	return &Entity{
		ID:     res.EntityView.EntityId,
		header: res.Header,
		data:   res.EntityView.Data,
	}, nil
}

// Close closes the underlying vertex AI gRPC client.
func (c *Client) Close() error {
	return c.v.Close()
}
