package main

import (
	"context"
	"log"

	"github.com/SchwarzIT/terraform-provider-stackit/stackit"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

func main() {
	err := providerserver.Serve(context.Background(), stackit.New, providerserver.ServeOpts{
		Address: "github.com/schwarzit/stackit",
	})

	if err != nil {
		log.Fatal(err.Error())
	}
}
