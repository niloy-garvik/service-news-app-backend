package utils

import "encoding/json"

type JSONData []byte
type MapData map[string]interface{}

func ConvertToJson(data interface{}) []byte {
	resultJson, jsonConversionError := json.Marshal(data)

	if jsonConversionError != nil {
		return JSONData{}
	} else {
		return resultJson
	}
}

func ConvertJsonToArrayOfMap(data []byte) []map[string]interface{} {
	var mapData []map[string]interface{}

	jsonConversionError := json.Unmarshal(data, &mapData)

	if jsonConversionError != nil {
		return []map[string]interface{}{}
	} else {
		return mapData
	}
}

func ConvertJsonToMap(data []byte) map[string]interface{} {
	var mapData map[string]interface{}

	jsonConversionError := json.Unmarshal(data, &mapData)

	if jsonConversionError != nil {
		return map[string]interface{}{}
	} else {
		return mapData
	}
}

func ConvertToMap(str string) MapData {

	var mapData map[string]interface{}
	jsonConversionError := json.Unmarshal([]byte(str), &mapData)

	if jsonConversionError != nil {
		return MapData{}
	} else {
		return mapData
	}

}

func StringifyObject(data interface{}) string {
	resultJson, jsonConversionError := json.Marshal(data)

	if jsonConversionError != nil {
		return string(JSONData{})
	} else {
		return string(resultJson)
	}
}
