package create

import (
	"fmt"

	"github.com/spf13/cobra"
	client "github.com/streamnative/function-mesh/api/generated/clientset/versioned"
	"github.com/tpiperatgod/fmw/pkg/sw"
	"github.com/tpiperatgod/fmw/pkg/util"
)

var (
	Create = &cobra.Command{
		Use:   "create",
		Short: "Create FunctionMesh with Serverlessworkflow configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := RunCreate()
			if err != nil {
				err = fmt.Errorf("[Create] %s", err)
				return err
			}
			return nil
		},
	}
	FunctionMeshClient client.Interface
)

func RunCreate() error {
	wf, err := sw.ParseWorkflow(util.FilePath)
	if err != nil {
		return err
	}

	if err = sw.CreateFunctionMesh(FunctionMeshClient, wf); err != nil {
		return err
	}
	return nil
}
