package face

import (
	"encoding/json"
	"errors"
)

/*
 * json face
 */

//face info
type Json struct {
}

//construct
func NewJson() *Json {
	this := &Json{}
	return this
}

//encode self
func (f *Json) EncodeSelf() ([]byte, error) {
	//encode json
	resp, err := json.Marshal(f)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

//decode json data
func (f *Json) Decode(data []byte, i interface{}) error {
	if data == nil || len(data) <= 0 || i == nil {
		return errors.New("invalid parameter")
	}
	err := json.Unmarshal(data, i)
	return err
}

//encode json data
func (f *Json) Encode(i interface{}) ([]byte, error) {
	if i == nil {
		return nil, errors.New("invalid parameter")
	}
	resp, err := json.Marshal(i)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

//decode simple kv data
func (f *Json) DecodeSimple(data []byte, kv map[string]interface{}) error {
	if data == nil || kv == nil {
		return errors.New("invalid parameter")
	}
	err := json.Unmarshal(data, &kv)
	return err
}

//encode simple kv data
func (f *Json) EncodeSimple(data map[string]interface{}) ([]byte, error) {
	if data == nil {
		return nil, errors.New("invalid parameter")
	}
	byte, err := json.Marshal(data)
	return byte, err
}
