package main

import (
	"context"
	"log"

	"github.com/SchwarzIT/terraform-provider-stackit/stackit"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

var (
	// goreleaser configuration will override this value
	version string = "dev"
)

func main() {

	err := providerserver.Serve(context.Background(), stackit.New(version), providerserver.ServeOpts{
		Address: "registry.terraform.io/schwarzit/stackit",
	})

	if err != nil {
		log.Fatal(err.Error())
	}
}
