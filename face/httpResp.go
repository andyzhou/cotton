package face

/*
 * http response face
 */

//face info
type HttpResp struct {
	Data []byte
	Err error
}

//construct
func NewHttpResp() *HttpResp {
	this := &HttpResp{}
	return this
}

//get response
func (f *HttpResp) GetResp() ([]byte, error) {
	return f.Data, f.Err
}

//set response
func (f *HttpResp) SetResp(data []byte, err error) bool {
	if data == nil && err == nil {
		return false
	}
	f.Data = data
	f.Err = err
	return true
}

//set err
func (f *HttpResp) SetErr(err error) bool {
	if err == nil {
		return false
	}
	f.Err = err
	return true
}