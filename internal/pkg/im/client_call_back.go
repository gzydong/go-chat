package im

type ClientCallBackInterface interface {
	Open(client ClientInterface)
	Message(message *ReceiveContent)
	Close(client ClientInterface, code int, text string)
}

type OpenCallBack func(client ClientInterface)

type MessageCallBack func(message *ReceiveContent)

type CloseCallBack func(client ClientInterface, code int, text string)

type ClientCallBackOption func(callBack *ClientCallBack)

type ClientCallBack struct {
	openCallBack    OpenCallBack
	messageCallBack MessageCallBack
	closeCallBack   CloseCallBack
}

func NewClientCallBack(opts ...ClientCallBackOption) *ClientCallBack {

	o := &ClientCallBack{}

	for _, opt := range opts {
		opt(o)
	}

	return o
}

func (c *ClientCallBack) Open(client ClientInterface) {
	if c.openCallBack != nil {
		c.openCallBack(client)
	}
}

func (c *ClientCallBack) Message(message *ReceiveContent) {
	if c.messageCallBack != nil {
		c.messageCallBack(message)
	}
}

func (c *ClientCallBack) Close(client ClientInterface, code int, text string) {
	if c.closeCallBack != nil {
		c.closeCallBack(client, code, text)
	}
}

func WithClientCallBackOpen(call OpenCallBack) ClientCallBackOption {
	return func(callBack *ClientCallBack) {
		callBack.openCallBack = call
	}
}

func WithClientCallBackMessage(call MessageCallBack) ClientCallBackOption {
	return func(callBack *ClientCallBack) {
		callBack.messageCallBack = call
	}
}

func WithClientCallBackClose(call CloseCallBack) ClientCallBackOption {
	return func(callBack *ClientCallBack) {
		callBack.closeCallBack = call
	}
}
