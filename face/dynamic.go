package face

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
)

/*
 * dynamic face
 */

//inter macro define
const (
	MaxInParams = 10
)

//face info
type Dynamic struct {
	faceMap sync.Map
}

//construct
func NewDynamic() *Dynamic {
	//self init
	this := &Dynamic{
		faceMap: sync.Map{},
	}
	return this
}

//call method on all faces
func (f *Dynamic) Cast(method string, params ...interface{}) bool {
	if method == "" {
		return false
	}
	//init parameters
	inParam := make([]reflect.Value, MaxInParams)
	paramNum := len(params)
	for i, para := range params {
		if i >= MaxInParams {
			break
		}
		inParam[i] = reflect.ValueOf(para)
	}
	//call method on each face
	subFunc := func(key interface{}, face interface{}) bool {
		faceVal, ok := face.(reflect.Value)
		if ok {
			faceVal.MethodByName(method).Call(inParam[0:paramNum])
			return true
		}else{
			return false
		}
	}
	f.faceMap.Range(subFunc)

	//reset in param slice
	inParam = []reflect.Value{}
	return true
}

//dynamic call method with parameters support
func (f *Dynamic) Call(name, method string, params ...interface{}) ([]reflect.Value, error) {
	var (
		tips string
	)

	//basic check
	if name == "" || method == "" {
		return nil, errors.New("invalid parameter")
	}

	//check instance
	face, isOk := f.faceMap.Load(name)
	if !isOk {
		tips = fmt.Sprintf("No face instance for name %s", name)
		return nil, errors.New(tips)
	}

	subFace, ok := face.(reflect.Value)
	if !ok {
		tips = fmt.Sprintf("Invalid face instance for name %s", name)
		return nil, errors.New(tips)
	}

	//init parameters
	inParam := make([]reflect.Value, 0)
	totalParas := 0
	//f.params = len(params)
	for _, para := range params {
		if totalParas >= MaxInParams {
			break
		}
		inParam = append(inParam, reflect.ValueOf(para))
		totalParas++
	}

	//dynamic call method with parameter
	callResult := subFace.MethodByName(method).Call(inParam[0:totalParas])

	return callResult, nil
}

//get face instance
func (f *Dynamic) GetFace(name string) interface{} {
	if name == "" {
		return nil
	}
	face, ok := f.faceMap.Load(name)
	if !ok {
		return nil
	}
	return face
}

//bind face with name
func (f *Dynamic) BindFace(name string, face interface{}) bool {
	if name == "" || face == nil {
		return false
	}
	//check is exists or not
	v := f.GetFace(name)
	if v != nil {
		return true
	}
	f.faceMap.Store(name, reflect.ValueOf(face))
	return true
}