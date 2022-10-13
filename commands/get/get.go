package get

import (
	"fmt"

	"github.com/spf13/cobra"
	client "github.com/streamnative/function-mesh/api/generated/clientset/versioned"
	"github.com/tpiperatgod/fmw/pkg/sw"
	"github.com/tpiperatgod/fmw/pkg/util"
)

type Get struct {
	functionMeshClient client.Interface
}

func NewCmdGet() *cobra.Command {
	c := &Get{}

	create := &cobra.Command{
		Use:   "get",
		Short: "Get FunctionMesh",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := c.preRunGet(); err != nil {
				err = fmt.Errorf("[Get] %s", err)
				return err
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := c.runGet(args); err != nil {
				err = fmt.Errorf("[Get] %s", err)
				return err
			}
			return nil
		},
	}
	return create
}

func (cmd *Get) preRunGet() error {
	config, _, err := util.NewKubeConfigClient()
	if err != nil {
		return err
	}

	cmd.functionMeshClient = client.NewForConfigOrDie(config)
	return nil
}

func (cmd *Get) runGet(args []string) error {
	var resourceName *string
	if len(args) > 0 {
		resourceName = &args[0]
	} else {
		resourceName = nil
	}
	if err := sw.GetFunctionMesh(cmd.functionMeshClient, resourceName); err != nil {
		return err
	}
	return nil
}
