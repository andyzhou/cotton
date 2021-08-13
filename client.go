package cotton

import (
	"github.com/andyzhou/cotton/face"
	"github.com/andyzhou/cotton/iface"
)

/*
 * http client
 */

//client face
type Client struct {
	queue iface.IHttpQueue
}

//construct
func NewClient() *Client {
	this := &Client{
		queue: face.NewHttpQueue(),
	}
	return this
}

//quit
func (c *Client) Quit() {
	c.queue.Quit()
}

//send request
func (c *Client) SendReq(req iface.IHttpReq) iface.IHttpResp {
	return c.queue.SendReq(req)
}