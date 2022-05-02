# `vertigo`

An alternative Vertex AI (AI Platform) Online Featurestore Client.

## Goals

The main goal of this project is to provide an alternative way of interacting with the Online Featurestore
Service. Specifically, several semantics are introduced in this API:

- `Query`
- `Config`
- `Entity`

The `Query` type contains parameterized inputs for the ReadFeatureValues RPC, containing
the entity ID, features to return (projection), and the entity type.

The `Config` type contains project level configuration, around the GCP Region, GCP Project ID,
and the Featurestore you are querying from.

Lastly, the `Entity` type is a wrapper struct around the `Header` and `EntityView.Data`
fields in the `ReadFeatureValuesResponse`.
The `Entity` has a receiver function that can scan the header and data into a user-provided
struct pointer.

You can provide an `interface{}` (pointer to a struct) type to the `entity.ScanStruct` function,
and as long as your struct has a `vertex` tag with the corresponding feature name, it will
load the values into the struct.

## Example

The following is an example of leveraging the Vertex AI Featurestore for a "customer" entity.

```go
package main

import (
	"context"
	"github.com/bradleybonitatibus/vertigo"
	"log"
)

// MyCustomer is a user provided struct that contain `"vertex"` tags that map to the entity feature ID in
// the Vertex Feature Store.
type MyCustomer struct {
	Segment               string   `json:"segment" vertex:"segment"`
	MarketAudiences       []string `json:"market_audiences" vertex:"market_audiences"`
	SixMonthSpend         *float64 `json:"six_month_spend" vertex:"six_month_spend"`
	AnotherNumericFeature int64    `json:"another_numeric_feature" vertex:"another_numeric_feature"`
}

func main() {
	var region string
	var projectID string
	var featurestoreName string
	// region = "my-gcp-region"
	// projectID = "my-project-id"
	// featurestoreName = "my_featurestore_name"

	client, err := vertigo.NewClient(context.Background(), &vertigo.Config{
		Region:           region,
		ProjectID:        projectID,
		FeatureStoreName: featurestoreName,
	})
	defer client.Close()
	if err != nil {
		log.Fatalf("vertigo.NewClient: %v", err)
	}
	myCust := MyCustomer{}

	entity, err := client.GetEntity(context.Background(), &vertigo.Query{
		EntityType: "my_customer",
		EntityID:   "123abc",
		Features:   []string{"*"},
	})
	if err != nil {
		log.Fatalf("client.GetEntity: %v", err)
	}
	err = entity.ScanStruct(&myCust)
	if err != nil {
		log.Fatalf("entity.ScanStruct: %v", err)
	}
	// continue using MyCustomer as you wish.
}
```
