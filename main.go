package main

import (
	"context"
	"flag"

	"github.com/pingidentity/terraform-provider-pingdirectory/internal/provider"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

func main() {
	var debug bool
	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	providerserver.Serve(context.Background(), provider.New, providerserver.ServeOpts{
		Address: "registry.terraform.io/pingidentity/pingdirectory",
		Debug:   debug,
	})
}
