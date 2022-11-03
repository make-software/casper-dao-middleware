package serialize

import (
	"encoding/json"
)

func ToRawJSON(target interface{}) map[string]interface{} {
	rawTarget, _ := json.Marshal(target)

	rawJSON := make(map[string]interface{})
	json.Unmarshal(rawTarget, &rawJSON)
	return rawJSON
}

func ToRawJSONList(target interface{}) []map[string]interface{} {
	rawTarget, _ := json.Marshal(target)

	rawJSONList := make([]map[string]interface{}, 0)
	json.Unmarshal(rawTarget, &rawJSONList)
	return rawJSONList
}
