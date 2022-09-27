package common

import (
	client "github.com/SchwarzIT/community-stackit-go-client"
	"github.com/hashicorp/terraform-plugin-framework/provider"
)

type Provider interface {
	provider.Provider
	// Returns true when the provider has been configured
	IsConfigured() bool

	// Client - returns the STACKIT client
	Client() *client.Client

	// ServiceAccountID - returns the service account id
	ServiceAccountID() string
}
