package iface

//jwt interface
type IJwt interface {
	Encode(input map[string]interface{}) (string, error)
	Decode(input string) (map[string]interface{}, error)
}
