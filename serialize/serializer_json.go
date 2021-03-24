package serialize

import "encoding/json"

var (
	JsonSerializer = &jsonSerializer{}
)

type Serializer interface {
	Marshal(obj interface{}) ([]byte, error)
	Unmarshal(data []byte, obj interface{}) error
}

type jsonSerializer struct{}

func (*jsonSerializer) Marshal(obj interface{}) ([]byte, error) {
	return json.Marshal(&obj)
}

func (*jsonSerializer) Unmarshal(data []byte, obj interface{}) error {
	return json.Unmarshal(data, &obj)
}
