package face

import (
	"errors"
	"github.com/emicklei/go-restful/v3"
	"github.com/gorilla/schema"
)

/*
 * rest tool face
 */

//face info
type Tool struct {
	decoder *schema.Decoder
}

//construct
func NewTool() *Tool {
	//self init
	this := &Tool{
		decoder: schema.NewDecoder(),
	}
	this.decoder.IgnoreUnknownKeys(true)
	return this
}

//parse request form
func (f *Tool) ParseReqForm(
				formFace interface{},
				req *restful.Request,
			) error {
	//basic check
	if formFace == nil || req == nil {
		return errors.New("invalid parameters")
	}

	//parse post form
	err := req.Request.ParseForm()
	if err != nil {
		return err
	}

	//decode form data
	err = f.decoder.Decode(formFace, req.Request.PostForm)
	if err != nil {
		return err
	}
	return nil
}