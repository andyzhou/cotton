package face

import (
	"errors"
	"github.com/andyzhou/cotton/iface"
	"github.com/emicklei/go-restful/v3"
	"log"
)

/*
 * api face
 */
const (
	EntryMethodDefault = "Entry"
)

type Api struct {
	dynamicFace iface.IDynamic
	entryMethod string
}

//construct
func NewApi() *Api {
	//self init
	this := &Api{
		dynamicFace: NewDynamic(),
		entryMethod: EntryMethodDefault,
	}
	return this
}

//dynamic call api face
func (f *Api) Call(
			name string,
			req *restful.Request,
			resp *restful.Response,
			extPathParams ...interface{},
		) error  {
	//check
	isExists := f.ApiIsExists(name)
	if !isExists {
		return errors.New("can't get face info")
	}

	//format final param slice
	finalParams := make([]interface{}, 0)
	finalParams = append(finalParams, req, resp)
	finalParams = append(finalParams, extPathParams...)

	//dynamic call api face
	_, err := f.dynamicFace.Call(
		name,
		f.entryMethod,
		finalParams...,
	)
	if err != nil {
		log.Printf("Api:Call, call %s failed, err:%v\n", name, err)
		return err
	}

	////analyze response data
	//if respSlice == nil || len(respSlice) < 2 {
	//	return define.ErrCodeInvalidRespData, nil
	//}
	//
	////get key data
	//errCode := respSlice[0].Int()
	//data := respSlice[1].Interface()
	//return int(errCode), data
	return nil
}

//check sub api face
func (f *Api) ApiIsExists(name string) bool {
	v := f.dynamicFace.GetFace(name)
	if v == nil {
		return false
	}
	return true
}

//bind api face
func (f *Api) BindApi(name string, face interface{}) error {
	if name == "" || face == nil {
		return errors.New("invalid parameter")
	}
	isExists := f.ApiIsExists(name)
	if isExists {
		return nil
	}
	bRet := f.dynamicFace.BindFace(name, false)
	if !bRet {
		return errors.New("bind face failed")
	}
	return nil
}

//set api face entry method
func (f *Api) SetEntryMethod(name string) {
	if name == "" {
		return
	}
	f.entryMethod = name
}