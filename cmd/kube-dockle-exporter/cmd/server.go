package cmd

import (
	"fmt"
	"kube-dockle-exporter/pkg/server"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func serverCmd() *cobra.Command {
	serverArgs := server.DefaultArgs()

	cmd := &cobra.Command{
		Use:          "server",
		Short:        "Starts KubeDockleExporter as a server",
		SilenceUsage: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				return fmt.Errorf("%q is an invalid argument", args[0])
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			err := server.Run(serverArgs)
			if err != nil {
				log.Fatalf("Failed to run server.Run: %s\n", err.Error())
			}
		},
	}

	cmd.PersistentFlags().StringVarP(
		&serverArgs.APIAddress,
		"api-address",
		"",
		serverArgs.APIAddress,
		"Address to use API",
	)
	cmd.PersistentFlags().Int64VarP(
		&serverArgs.APIMaxConnections,
		"api-max-connections",
		"",
		serverArgs.APIMaxConnections,
		"Max connections of API",
	)
	cmd.PersistentFlags().StringVarP(
		&serverArgs.MonitorAddress,
		"monitor-address",
		"",
		serverArgs.MonitorAddress,
		"Address to use self-monitoring information",
	)
	cmd.PersistentFlags().Int64VarP(
		&serverArgs.MonitorMaxConnections,
		"monitor-max-connections",
		"",
		serverArgs.MonitorMaxConnections,
		"Max connections of self-monitoring information",
	)
	cmd.PersistentFlags().StringVarP(
		&serverArgs.MonitoringJaegerEndpoint,
		"monitoring-jaeger-endpoint",
		"",
		serverArgs.MonitoringJaegerEndpoint,
		"Address to use for distributed tracing",
	)
	cmd.PersistentFlags().BoolVarP(
		&serverArgs.EnableProfiling,
		"enable-profiling",
		"",
		serverArgs.EnableProfiling,
		"Enable profiling",
	)
	cmd.PersistentFlags().BoolVarP(
		&serverArgs.EnableTracing,
		"enable-tracing",
		"",
		serverArgs.EnableTracing,
		"Enable distributed tracing",
	)
	cmd.PersistentFlags().Float64VarP(
		&serverArgs.TracingSampleRate,
		"tracing-sample-rate",
		"",
		serverArgs.TracingSampleRate,
		"Tracing sample rate",
	)
	cmd.PersistentFlags().BoolVarP(
		&serverArgs.KeepAlived,
		"enable-keep-alived",
		"",
		serverArgs.KeepAlived,
		"Enable HTTP KeepAlive",
	)
	cmd.PersistentFlags().BoolVarP(
		&serverArgs.ReUsePort,
		"enable-reuseport",
		"",
		serverArgs.ReUsePort,
		"Enable SO_REUSEPORT",
	)
	cmd.PersistentFlags().Int64VarP(
		&serverArgs.TCPKeepAliveInterval,
		"tcp-keep-alive-interval",
		"",
		serverArgs.TCPKeepAliveInterval,
		"Interval of TCP KeepAlive",
	)
	cmd.PersistentFlags().Int64VarP(
		&serverArgs.DockleConcurrency,
		"dockle-concurrency",
		"",
		serverArgs.DockleConcurrency,
		"Concurrency of dockle execution",
	)
	cmd.PersistentFlags().Int64VarP(
		&serverArgs.CollectorLoopInterval,
		"collector-loop-interval",
		"",
		serverArgs.CollectorLoopInterval,
		"Interval to execute collect result from dockle",
	)
	cmd.PersistentFlags().BoolVarP(
		&serverArgs.Verbose,
		"verbose",
		"",
		serverArgs.Verbose,
		"Verbose logging",
	)

	if err := viper.BindPFlags(cmd.PersistentFlags()); err != nil {
		log.Fatalf("Failed to execute command: %s\n", err)
	}

	return cmd
}
