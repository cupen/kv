package serialize

type Serializer interface {
	Marshal(obj interface{}) ([]byte, error)
	Unmarshal(data []byte, obj interface{}) error
}
