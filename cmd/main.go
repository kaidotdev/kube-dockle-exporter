package main

import (
	"kube-dockle-exporter/cmd/kube-dockle-exporter/cmd"
	"os"
)

func main() {
	rootCmd := cmd.GetRootCmd(os.Args[1:])

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
