package commands

import (
	"github.com/spf13/cobra"
	"github.com/tpiperatgod/fmw/commands/create"
	"github.com/tpiperatgod/fmw/pkg/util"
)

var (
	version string
)

var Root = &cobra.Command{
	Use:           "fmw command [flags]",
	Short:         "Serverlessworkflow Specification Conversion Tool for FunctionMesh",
	Version:       version,
	SilenceErrors: true,
	SilenceUsage:  true,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func Execute() error {
	Root.AddCommand(create.NewCmdCreate())

	Root.PersistentFlags().StringVar(&util.KubeConfig, "kubeconfig", "", "config of Kubernetes cluster")
	Root.PersistentFlags().StringVarP(&util.FilePath, "file-path", "f", "", "file path")
	Root.PersistentFlags().StringVarP(&util.Namespace, "namespace", "n", "default", "namespace")

	return Root.Execute()
}
