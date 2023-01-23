package main

import (
	"encoding/json"
	"fmt"
	"log"

	cancom "github.com/cancom/terraform-provider-cancom/cancom"
	"github.com/cancom/terraform-provider-cancom/internal/sdk"
	"github.com/cancom/terraform-provider-cancom/internal/services/client_car_rental"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func PrettyPrint(v interface{}) (err error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err == nil {
		fmt.Println(string(b))
	}
	return
}

func main() {
	client := client_car_rental.NewClientCarRental("tesf")
	cc := ((interface{})(client)).(*sdk.CancomClient)
	log.Println((cc))
}

func main2() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() *schema.Provider {
			return cancom.Provider()
		},
	})
}
