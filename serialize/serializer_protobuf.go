package serialize

var (
	ProtobufSerializer = &protobufSerializer{}
)

type ProtobufMessage interface {
	Marshal() ([]byte, error)
	Unmarshal(data []byte) error
}

type protobufSerializer struct{}

func (*protobufSerializer) Marshal(obj ProtobufMessage) ([]byte, error) {
	return obj.Marshal()
}

func (*protobufSerializer) Unmarshal(data []byte, obj ProtobufMessage) error {
	return obj.Unmarshal(data)
}
