package main

import (
	"context"
	"log"
	"os"

	"github.com/schrodit/helm-cleanup/pkg/cleanup"
	"github.com/schrodit/helm-cleanup/pkg/common"
	"github.com/spf13/cobra"
)

// Options contains the options for Map operation
type Options struct {
	DryRun    bool
	Debug     bool
	Namespace string
}

var (
	settings *EnvSettings
)

func newRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "cleanup",
		Short:        "Cleanup helm leftovers",
		Long:         `Identifies and cleans up resources that have been managed by Helm.`,
		SilenceUsage: true,

		RunE: func(cmd *cobra.Command, args []string) error {
			opts := Options{
				DryRun:    settings.DryRun,
				Debug:     settings.Debug,
				Namespace: settings.Namespace,
			}
			kubeConfig := common.KubeConfig{
				Context: settings.KubeContext,
				File:    settings.KubeConfigFile,
			}

			return Cleanup(opts, kubeConfig)
		},
	}

	flags := cmd.PersistentFlags()
	flags.ParseErrorsWhitelist.UnknownFlags = true

	settings = new(EnvSettings)

	// When run with the Helm plugin framework, Helm plugins are not passed the
	// plugin flags that correspond to Helm global flags e.g. helm mapkubeapis v3map --kube-context ...
	// The flag values are set to corresponding environment variables instead.
	// The flags are passed as expected when run directly using the binary.
	// The below allows to use Helm's --kube-context global flag.
	if ctx := os.Getenv("HELM_KUBECONTEXT"); ctx != "" {
		settings.KubeContext = ctx
	}

	if debug := os.Getenv("HELM_DEBUG"); debug != "" {
		settings.Debug = true
	}

	// Note that the plugin's --kubeconfig flag is set by the Helm plugin framework to
	// the KUBECONFIG environment variable instead of being passed into the plugin.

	settings.AddFlags(flags)

	return cmd
}

// Checks all resources managed by helm and cleans up resources that do not
// have a corresponding Helm release.
func Cleanup(opts Options, kubeConfig common.KubeConfig) error {
	if opts.DryRun {
		log.Println("NOTE: This is in dry-run mode, the following actions will not be executed.")
		log.Println("Run without --dry-run to take the actions described below:")
		log.Println()
	}

	options := common.Options{
		DryRun:     opts.DryRun,
		Debug:      opts.Debug,
		KubeConfig: kubeConfig,
		Namespace:  opts.Namespace,
	}

	leaked, err := cleanup.ListLeakedResources(context.TODO(), options)
	if err != nil {
		return err
	}
	log.Printf("Found %d leaked resources\n", len(leaked))
	if opts.DryRun || opts.Debug {
		log.Println("LEAKED RESOURCES:")
		for _, u := range leaked {
			log.Printf("(%s %s) %s %s\n", u.GetAPIVersion(), u.GetKind(), u.GetName(), u.GetNamespace())
		}
	}
	if opts.DryRun {
		return nil
	}

	return nil
}
