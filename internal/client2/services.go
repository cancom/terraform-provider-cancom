package client2

import (
	"github.com/cancom/terraform-provider-cancom/internal/sdk"
	"github.com/cancom/terraform-provider-cancom/internal/services/client_car_rental"
)

type FactoryFunc func(token string, endpoint *string) *sdk.CancomClient
type factoryFuncWithout func(token string) *sdk.CancomClient

type service struct {
	Name    string
	Factory FactoryFunc
}

func createAndCast(withEndpoint interface{}, without interface{}) FactoryFunc {
	return func(token string, endpoint *string) *sdk.CancomClient {
		if endpoint == nil {
			return without.(factoryFuncWithout)(token)
		}

		return withEndpoint.(FactoryFunc)(token, endpoint)
	}
}

var SupportedServices = map[string]*service{
	"car-rental": {Name: "car-rental", Factory: createAndCast(client_car_rental.NewClientCarRentalWithEndpoint, client_car_rental.NewClientCarRental)},
}
