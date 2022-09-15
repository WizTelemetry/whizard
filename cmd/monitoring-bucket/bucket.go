package main

import (
	"context"
	"flag"
	"github.com/kubesphere/whizard/pkg/bucket"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/klog/v2"
	"os"
	"time"
)

var (
	tenantLabelName   string
	defaultTenantId   string
	storageConfig     string
	storageConfigFile string
	interval          time.Duration
	cleanupTimeout    time.Duration
)

func AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&tenantLabelName, "tenant.label-name", "tenant_id", "Label name through which the tenant will be announcede")
	fs.StringVar(&defaultTenantId, "tenant.default-id", "default-tenant", "Default tenant ID to use when none is provided via a header.")
	fs.StringVar(&storageConfig, "objstore.config", "", "The storage config used to access the object storage")
	fs.StringVar(&storageConfigFile, "objstore.config-file", "", "The storage config file used to access the object storage")
	fs.DurationVar(&interval, "interval", time.Minute, "The garbage collection interval")
	fs.DurationVar(&cleanupTimeout, "cleanup-timeout", time.Hour, "The maximum time a cleanup can execute")
}

func NewCommand() *cobra.Command {

	cmd := &cobra.Command{
		Use:          "bucket garbage collection",
		Short:        `Whizard bucket garbage collection`,
		Run:          run,
		SilenceUsage: true,
	}

	AddFlags(cmd.Flags())
	cmd.Flags().AddGoFlagSet(flag.CommandLine)

	return cmd
}

func run(_ *cobra.Command, _ []string) {
	if storageConfig == "" && storageConfigFile == "" {
		klog.Errorf("storage config or storage config file must be specified")
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	b := bucket.NewBackend(ctx, tenantLabelName, defaultTenantId, storageConfig, storageConfigFile, interval, cleanupTimeout)
	if err := b.Run(); err != nil {
		klog.Error(err)
		os.Exit(1)
	}
}

func main() {
	command := NewCommand()

	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
