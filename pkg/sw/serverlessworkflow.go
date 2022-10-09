package sw

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/serverlessworkflow/sdk-go/v2/model"
	"github.com/serverlessworkflow/sdk-go/v2/parser"
	"github.com/streamnative/function-mesh/api/compute/v1alpha1"
	client "github.com/streamnative/function-mesh/api/generated/clientset/versioned"
	"github.com/tpiperatgod/fmw/pkg/util"
)

const (
	TypeFunction = "function"
	TypeSource   = "source"
	TypeSink     = "sink"
)

var (
	AvailableTypes = map[string]bool{
		TypeFunction: true,
		TypeSource:   true,
		TypeSink:     true,
	}
)

func ParseWorkflow(filePath string) (*model.Workflow, error) {
	workflow, err := parser.FromFile(filePath)
	if err != nil {
		return nil, err
	}
	if err := validate(workflow); err != nil {
		return nil, err
	}
	return workflow, nil
}

func validate(workflow *model.Workflow) error {
	// validate functions
	functions := workflow.Functions
	for _, function := range functions {
		if function.Type != "custom" {
			return fmt.Errorf("function [%s] parse error: unknown type: %s", function.Name, function.Type)
		}
		if function.Metadata == nil {
			return fmt.Errorf("function [%s] parse error: metadata required", function.Name)
		}
		if t, ok := function.Metadata["type"]; !ok {
			return fmt.Errorf("function [%s] parse error: metadata.type required", function.Name)
		} else {
			if !AvailableTypes[t.(string)] {
				return fmt.Errorf("function [%s] parse error: invalid metadata.type: %s", function.Name, t.(string))
			}
		}
		if function.Operation == "" && function.Metadata["spec"] == nil {
			return fmt.Errorf("function [%s] parse error: specification required", function.Name)
		}
		if function.Operation != "" {
			filePath := strings.TrimPrefix(function.Operation, "file://")
			if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
				return fmt.Errorf("function [%s] parse error: operation file not exist: %s, err: %s", function.Name, function.Operation, err.Error())
			}
		}
	}
	return nil
}

func FetchFunctions(workflow *model.Workflow) map[string]*model.Function {
	funcMap := map[string]*model.Function{}
	functions := workflow.Functions
	for _, function := range functions {
		f := function
		funcMap[function.Name] = &f
	}
	return funcMap
}

func FetchOrder(workflow *model.Workflow) []string {
	order := []string{}
	start := ""
	if workflow.Start != nil {
		start = workflow.Start.StateName
		order = append(order, start)
	}
	states := workflow.States
	for _, state := range states {
		if start != "" && start == state.GetName() {
			continue
		}
		order = append(order, state.GetName())
	}
	return order
}

func CreateFunctionMesh(client client.Interface, workflow *model.Workflow) error {
	functions := FetchFunctions(workflow)
	order := FetchOrder(workflow)
	fmt.Println(functions)
	functionMesh := &v1alpha1.FunctionMesh{}
	functionMesh.Name = workflow.Name
	functionMesh.Namespace = util.Namespace

	for _, name := range order {
		fmt.Println(name)
		resource := functions[name]
		fmt.Println(resource)
		if resource.Operation != "" {
			if resourceBytes, err := ioutil.ReadFile(strings.TrimPrefix(resource.Operation, "file://")); err != nil {
				fmt.Println("read operation error:", err.Error())
				return err
			} else {
				if err = createResourceWithOperation(resourceBytes, resource, functionMesh); err != nil {
					return err
				}
			}
			continue
		}
		if spec, exist := resource.Metadata["spec"]; exist {
			if specBytes, err := json.Marshal(spec); err != nil {
				fmt.Println("parse function spec error:", err.Error())
				return err
			} else {
				if err = createResourceWithSpec(specBytes, resource, functionMesh); err != nil {
					return err
				}
			}
		}
	}

	// TODO: create FunctionMesh
	return nil
}

func createResourceWithOperation(resourceBytes []byte, resource *model.Function, functionMesh *v1alpha1.FunctionMesh) error {
	resourceType := resource.Metadata["type"]
	switch resourceType {
	case TypeSink:
		sink := &v1alpha1.Sink{}
		if err := yaml.Unmarshal(resourceBytes, sink); err != nil {
			fmt.Println("parse sink error:", err.Error())
			return err
		} else {
			functionMesh.Spec.Sinks = append(functionMesh.Spec.Sinks, sink.Spec)
			fmt.Println("add sink:", resource.Name)
			return nil
		}
	case TypeSource:
		source := &v1alpha1.Source{}
		if err := yaml.Unmarshal(resourceBytes, source); err != nil {
			fmt.Println("parse source error:", err.Error())
			return err
		} else {
			functionMesh.Spec.Sources = append(functionMesh.Spec.Sources, source.Spec)
			fmt.Println("add source:", resource.Name)
			return nil
		}
	case TypeFunction:
		function := &v1alpha1.Function{}
		if err := yaml.Unmarshal(resourceBytes, function); err != nil {
			fmt.Println("parse function error:", err.Error())
			return err
		} else {
			functionMesh.Spec.Functions = append(functionMesh.Spec.Functions, function.Spec)
			fmt.Println("add function:", resource.Name)
			return nil
		}
	default:
		return errors.New(fmt.Sprintf("unknown resource type: %s", resourceType))
	}
}

func createResourceWithSpec(specBytes []byte, resource *model.Function, functionMesh *v1alpha1.FunctionMesh) error {
	resourceType := resource.Metadata["type"]
	switch resourceType {
	case TypeSink:
		sinkSpec := &v1alpha1.SinkSpec{}
		if err := yaml.Unmarshal(specBytes, sinkSpec); err != nil {
			fmt.Println("parse sink spec error:", err.Error())
			return err
		} else {
			functionMesh.Spec.Sinks = append(functionMesh.Spec.Sinks, *sinkSpec)
			fmt.Println("add sink:", resource.Name)
			return nil
		}
	case TypeSource:
		sourceSpec := &v1alpha1.SourceSpec{}
		if err := yaml.Unmarshal(specBytes, sourceSpec); err != nil {
			fmt.Println("parse source spec error:", err.Error())
			return err
		} else {
			functionMesh.Spec.Sources = append(functionMesh.Spec.Sources, *sourceSpec)
			fmt.Println("add source:", resource.Name)
			return nil
		}
	case TypeFunction:
		functionSpec := &v1alpha1.FunctionSpec{}
		if err := json.Unmarshal(specBytes, functionSpec); err != nil {
			fmt.Println("parse function spec error:", err.Error())
			return err
		} else {
			functionMesh.Spec.Functions = append(functionMesh.Spec.Functions, *functionSpec)
			fmt.Println("add source:", resource.Name)
			return nil
		}
	default:
		return errors.New(fmt.Sprintf("unknown resource type: %s", resourceType))
	}
}
