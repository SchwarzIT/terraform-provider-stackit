package main

import (
	"context"
	"log"

	"github.com/SchwarzIT/terraform-provider-stackit/stackit"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

// Generate the Terraform provider documentation using `tfplugindocs`:
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs@latest

var (
	// goreleaser configuration will override this value
	version string = "dev"
)

func main() {

	err := providerserver.Serve(context.Background(), stackit.New(version), providerserver.ServeOpts{
		Address:         "registry.terraform.io/schwarzit/stackit",
		ProtocolVersion: 6,
	})

	if err != nil {
		log.Fatal(err.Error())
	}
}
