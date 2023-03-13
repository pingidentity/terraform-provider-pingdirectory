package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/provider"
)

// Run "go generate" to format example terraform files, generate the docs for the registry/website, and
// generate resource source files

// Format examples
//go:generate terraform fmt -recursive ./examples/

// Generate resource source files
//go:generate ./scripts/generate-resource-source-files.sh

// Run the docs generation tool
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

func main() {
	var debug bool
	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	err := providerserver.Serve(context.Background(), provider.New, providerserver.ServeOpts{
		Address: "registry.terraform.io/pingidentity/pingdirectory",
		Debug:   debug,
	})

	if err != nil {
		fmt.Println(err)
	}
}
