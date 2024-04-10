package main

import (
	"context"
	"flag"
	"log"

	"github.com/SchwarzIT/terraform-provider-stackit/stackit"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

var (
	// goreleaser configuration will override this value
	version string = "dev"
)

func main() {
	var debug bool
	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := providerserver.ServeOpts{
		Address: "github.com/schwarzit/stackit",
		Debug:   debug,
	}

	err := providerserver.Serve(context.Background(), stackit.New(version), opts)
	if err != nil {
		log.Fatal(err.Error())
	}
}
