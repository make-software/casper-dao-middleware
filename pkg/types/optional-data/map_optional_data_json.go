package optional_data

import (
	"errors"
	"fmt"

	"casper-dao-middleware/pkg/types"
)

func MapJSON(optional *types.OptionalData, target, data map[string]interface{}) error {
	// if no data found for mapping just set is as nil
	if data == nil {
		target[optional.GetName()] = nil
		return nil
	}

	dataMap := map[string]interface{}{
		optional.GetName(): data,
	}

	return mapWithOptionalJson(optional, target, dataMap)
}

func mapWithOptionalJson(optional *types.OptionalData, target, data map[string]interface{}) error {
	// if no data found for mapping just set it to nil
	if data == nil || len(data) == 0 {
		target[optional.GetName()] = nil
		return nil
	}

	if len(optional.GetNested()) == 0 {
		var mapped interface{}
		if value, ok := data[optional.GetName()]; ok {
			mapped = value
		}
		target[optional.GetName()] = mapped
		return nil
	}

	for _, child := range optional.GetNested() {
		var ok bool
		var childNode interface{}

		childNode, ok = target[optional.GetName()]
		if !ok {
			childNode = make(map[string]interface{})
			target[optional.GetName()] = childNode
		}

		childData, ok := data[optional.GetName()]
		if !ok {
			return fmt.Errorf("data and optional tree is no identical no field %s", child.GetName())
		}

		childDataMap, ok := childData.(map[string]interface{})
		if !ok {
			return errors.New("")
		}

		if err := mapWithOptionalJson(child, childNode.(map[string]interface{}), childDataMap); err != nil {
			return err
		}
	}

	return nil
}
