package serialize

var (
	Protobuf = &protobufSerializer{}
)

type protoMessage interface {
	Marshal() ([]byte, error)
	Unmarshal(data []byte) error
}

type protobufSerializer struct{}

func (*protobufSerializer) Marshal(obj protoMessage) ([]byte, error) {
	return obj.Marshal()
}

func (*protobufSerializer) Unmarshal(data []byte, obj protoMessage) error {
	return obj.Unmarshal(data)
}
