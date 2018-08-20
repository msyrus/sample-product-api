package main

import (
	"github.com/msyrus/simple-product-inv/version"
	"github.com/spf13/cobra"
)

// rootCmd is the root of all sub commands in the binary
// it doesn't have a Run method as it executes other sub commands
var rootCmd = &cobra.Command{
	Use:     "product",
	Short:   "product is a http server to serve products",
	Version: version.Version,
}

func init() {
	// Here all other sub commands should be registered to the rootCmd
}

func main() {
	rootCmd.Execute()
}
