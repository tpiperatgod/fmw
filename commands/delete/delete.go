package delete

import (
	"fmt"

	"github.com/spf13/cobra"
	client "github.com/streamnative/function-mesh/api/generated/clientset/versioned"
	"github.com/tpiperatgod/fmw/pkg/sw"
	"github.com/tpiperatgod/fmw/pkg/util"
)

type Delete struct {
	functionMeshClient client.Interface
}

func NewCmdDelete() *cobra.Command {
	c := &Delete{}

	create := &cobra.Command{
		Use:   "delete",
		Short: "Delete FunctionMesh",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := c.preRunDelete(); err != nil {
				err = fmt.Errorf("[Delete] %s", err)
				return err
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := c.runDelete(args); err != nil {
				err = fmt.Errorf("[Delete] %s", err)
				return err
			}
			return nil
		},
	}
	return create
}

func (cmd *Delete) preRunDelete() error {
	config, _, err := util.NewKubeConfigClient()
	if err != nil {
		return err
	}

	cmd.functionMeshClient = client.NewForConfigOrDie(config)
	return nil
}

func (cmd *Delete) runDelete(args []string) error {
	resources := map[string]bool{}
	resourcesAfterDeduplicated := []string{}
	for _, name := range args {
		if resources[name] {
			continue
		}
		resources[name] = true
		resourcesAfterDeduplicated = append(resourcesAfterDeduplicated, name)
	}
	if err := sw.DeleteFunctionMesh(cmd.functionMeshClient, resourcesAfterDeduplicated); err != nil {
		return err
	}
	return nil
}
