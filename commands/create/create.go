package create

import (
	"fmt"

	"github.com/spf13/cobra"
	client "github.com/streamnative/function-mesh/api/generated/clientset/versioned"
	"github.com/tpiperatgod/fmw/pkg/sw"
	"github.com/tpiperatgod/fmw/pkg/util"
)

type Create struct {
	functionMeshClient client.Interface
}

func NewCmdCreate() *cobra.Command {
	c := &Create{}

	create := &cobra.Command{
		Use:   "create",
		Short: "Create FunctionMesh with Serverlessworkflow configuration",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := c.preRunCreate(); err != nil {
				err = fmt.Errorf("[Create] %s", err)
				return err
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := c.runCreate(); err != nil {
				err = fmt.Errorf("[Create] %s", err)
				return err
			}
			return nil
		},
	}
	return create
}

func (cmd *Create) preRunCreate() error {
	config, _, err := util.NewKubeConfigClient()
	if err != nil {
		return err
	}

	cmd.functionMeshClient = client.NewForConfigOrDie(config)
	return nil
}

func (cmd *Create) runCreate() error {
	wf, err := sw.ParseWorkflow(util.FilePath)
	if err != nil {
		return err
	}

	if err = sw.CreateFunctionMesh(cmd.functionMeshClient, wf); err != nil {
		return err
	}
	return nil
}
