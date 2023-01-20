package main

import (
	"os"

	"github.com/mattn/go-colorable"

	"github.com/hashicorp/terraform-plugin-docs/internal/cmd"
)

func main() {
	name := "tfplugindocs"

	os.Exit(cmd.Run(
		name,
		"dev",
		os.Args[1:],
		os.Stdin,
		colorable.NewColorableStdout(),
		colorable.NewColorableStderr(),
	))
}
