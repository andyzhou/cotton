package iface

//jwt interface
type IJwt interface {
	Encode(input map[string]interface{}) string
	Decode(input string) map[string]interface{}
}
