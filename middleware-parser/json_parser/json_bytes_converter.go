package json_parser

import "encoding/json"

func ConvertObjectToByteArray(body interface{}) ([]byte, error) {
	return json.Marshal(body)
}
func ConvertByteArrayToObject(bytes []byte, typed interface{}) error {
	return json.Unmarshal(bytes, typed)
}
