package iface

//json interface
type IJson interface {
	EncodeSelf() ([]byte, error)
	Decode(data []byte, i interface{}) error
	Encode(i interface{}) ([]byte, error)
	DecodeSimple(data []byte, kv map[string]interface{}) error
	EncodeSimple(data map[string]interface{}) ([]byte, error)
}
