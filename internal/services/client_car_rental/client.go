package client_car_rental

import (
	"fmt"

	"github.com/cancom/terraform-provider-cancom/internal/sdk"
)

type ClientCarRental struct {
	sdk.CancomClient
}

func NewClientCarRental(token string) *ClientCarRental {
	return &ClientCarRental{
		CancomClient: sdk.CancomClient{
			Token:   token,
			HostURL: fmt.Sprintf("https://%s.%s", "car-rental", sdk.MainEndpoint),
		},
	}
}

func NewClientCarRentalWithEndpoint(token, endpoint string) *ClientCarRental {
	return &ClientCarRental{
		CancomClient: sdk.CancomClient{
			Token:   token,
			HostURL: endpoint,
		},
	}
}
