package cmd

import "github.com/spf13/cobra"

func GetRootCmd(args []string) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:          "kube-dockle-exporter",
		Short:        "KubeDockleExporter is Prometheus Exporter that collects CIS benchmark executed by goodwithtech/dockle in the cluster.",
		SilenceUsage: true,
	}

	rootCmd.SetArgs(args)
	rootCmd.AddCommand(serverCmd())

	return rootCmd
}
