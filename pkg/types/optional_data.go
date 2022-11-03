package types

import (
	"fmt"
	"strings"
)

var ErrInvalidSchemaSyntax = fmt.Errorf("invalid schema syntax")

const MaxSchemaLengthBytes = 2048

// OptionalFuncArgs describe optional func arguments could be parsed as rate(1)
type OptionalFuncArgs struct {
	params []string
}

// OptionalData represent tree based structure for handling optional field which can be passed as query params
// e.g. account_info{owner{name,branding}}
type OptionalData struct {
	name       string
	funcArgs   *OptionalFuncArgs
	nestedData map[string]*OptionalData
}

func (o *OptionalData) GetNested() map[string]*OptionalData {
	return o.nestedData
}

func (o *OptionalData) GetName() string {
	return o.name
}

func (o *OptionalFuncArgs) GetFuncParams() []string {
	return o.params
}

func (o *OptionalData) GetNestedByName(name string) *OptionalData {
	nested, ok := o.nestedData[name]
	if !ok {
		return nil
	}
	return nested
}

func (o *OptionalData) ContainsFunc(name string) ([]string, bool) {
	if o.nestedData == nil {
		return nil, false
	}

	data, ok := o.nestedData[name]
	if !ok {
		return nil, false
	}

	return data.funcArgs.params, true
}

func (o *OptionalData) Contains(name string) (*OptionalData, bool) {
	if o.nestedData == nil {
		return nil, false
	}

	data, ok := o.nestedData[name]
	return data, ok
}

func ParseOptionalData(fieldsSchema string) (*OptionalData, error) {
	buffer := make([]byte, MaxSchemaLengthBytes)

	n, err := strings.NewReader(fieldsSchema).Read(buffer)
	if err != nil {
		return nil, err
	}

	if strings.Count(fieldsSchema, "{") != strings.Count(fieldsSchema, "}") {
		return nil, ErrInvalidSchemaSyntax
	}

	if strings.Count(fieldsSchema, "(") != strings.Count(fieldsSchema, ")") {
		return nil, ErrInvalidSchemaSyntax
	}

	var openBracketsCount, openFunctionCount, startObject int
	var rootData, _ = newOptionalDataNode("root")

	for i := 0; i < n; i++ {
		if buffer[i] == '{' {
			openBracketsCount++
		}

		if buffer[i] == '(' {
			openFunctionCount++
		}

		if buffer[i] == ')' {
			openFunctionCount--

			// found and of function
			if openFunctionCount == 0 {
				optionalDataSchema := strings.TrimSpace(string(buffer[startObject : i+1]))
				funcArgs, err := parseOptionalFuncArgs(optionalDataSchema)
				if err != nil {
					return nil, err
				}
				rootData.nestedData[funcArgs.GetName()] = funcArgs
				startObject = i + 1
			}
		}

		if buffer[i] == '}' {
			openBracketsCount--

			// found closing fields object
			if openBracketsCount == 0 {
				optionalDataSchema := strings.TrimSpace(string(buffer[startObject : i+1]))
				optionalData, err := parseOptionalData(optionalDataSchema)
				if err != nil {
					return nil, err
				}
				rootData.nestedData[optionalData.GetName()] = optionalData
				startObject = i + 1
			}
		}
	}
	return rootData, nil
}

func parseOptionalData(schema string) (*OptionalData, error) {
	schema = strings.Trim(schema, ",")

	buffer := make([]byte, len(schema))
	n, err := strings.NewReader(schema).Read(buffer)
	if err != nil {
		return nil, err
	}

	var initNode *OptionalData
	var processedObject bool
	var newNodeStartIndex = 0
	var activeParentNodeStack = NewEmptyStack[OptionalData]()

	for i := 0; i < n; i++ {
		if buffer[i] == '}' {
			if newNodeStartIndex == i {
				activeParentNodeStack.Pop()
				newNodeStartIndex = i + 1
				continue
			}

			nodeName := buffer[newNodeStartIndex:i]

			parent := activeParentNodeStack.Top()
			node, err := newOptionalDataNode(string(nodeName))
			if err != nil {
				return nil, ErrInvalidSchemaSyntax
			}

			parent.nestedData[node.GetName()] = node

			newNodeStartIndex = i + 1
			activeParentNodeStack.Pop()
			processedObject = true
		}

		// next item of current node are coming
		if buffer[i] == ',' {
			if processedObject {
				newNodeStartIndex = i + 1
				continue
			}

			nodeName := buffer[newNodeStartIndex:i]
			parent := activeParentNodeStack.Top()

			node, err := newOptionalDataNode(string(nodeName))
			if err != nil {
				return nil, ErrInvalidSchemaSyntax
			}

			parent.nestedData[node.GetName()] = node

			newNodeStartIndex = i + 1
		}

		if buffer[i] != '{' {
			continue
		}

		parentNode := activeParentNodeStack.Top()

		nodeName := buffer[newNodeStartIndex:i]
		node, err := newOptionalDataNode(string(nodeName))
		if err != nil {
			return nil, ErrInvalidSchemaSyntax
		}

		if parentNode == nil {
			initNode = node
		} else {
			parentNode.nestedData[node.GetName()] = node
		}

		activeParentNodeStack.Push(node)
		newNodeStartIndex = i + 1
		processedObject = false
	}

	return initNode, nil
}

func parseOptionalFuncArgs(schema string) (*OptionalData, error) {
	schema = strings.Trim(schema, ",")

	buffer := make([]byte, len(schema))
	n, err := strings.NewReader(schema).Read(buffer)
	if err != nil {
		return nil, err
	}

	var newNodeStartIndex = 0
	var funcName, funcParam string
	for i := 0; i < n; i++ {

		if buffer[i] == '(' {
			funcName = string(buffer[newNodeStartIndex:i])
			newNodeStartIndex = i + 1
		}

		if buffer[i] == ')' {
			funcParam = string(buffer[newNodeStartIndex:i])
		}
	}
	return newFunctionArgs(funcName, funcParam)
}

func newFunctionArgs(name, param string) (*OptionalData, error) {
	trimmedName := strings.TrimSpace(name)
	if trimmedName == "" {
		return nil, fmt.Errorf("empty function name provided")
	}

	trimmedParams := strings.TrimSpace(param)

	params := make([]string, 0)
	for _, param := range strings.Split(trimmedParams, ",") {
		params = append(params, strings.TrimSpace(param))
	}

	return &OptionalData{
		name: strings.TrimSpace(trimmedName),
		funcArgs: &OptionalFuncArgs{
			params,
		},
	}, nil
}

func newOptionalDataNode(name string) (*OptionalData, error) {
	trimmed := strings.TrimSpace(name)
	if trimmed == "" {
		return nil, fmt.Errorf("empty name provided")
	}

	return &OptionalData{
		name:       strings.TrimSpace(name),
		nestedData: make(map[string]*OptionalData),
	}, nil
}
