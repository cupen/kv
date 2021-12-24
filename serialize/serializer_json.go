package serialize

import "encoding/json"

var (
	Json = &jsonSerializer{}
)

type jsonSerializer struct{}

func (*jsonSerializer) Marshal(obj interface{}) ([]byte, error) {
	return json.Marshal(&obj)
}

func (*jsonSerializer) Unmarshal(data []byte, obj interface{}) error {
	return json.Unmarshal(data, &obj)
}
