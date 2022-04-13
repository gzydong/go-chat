package im

type ClientCallbackInterface interface {
	Open(client ClientInterface)
	Message(message *ReceiveContent)
	Close(client ClientInterface, code int, text string)
}

type OpenCallback func(client ClientInterface)

type MessageCallback func(message *ReceiveContent)

type CloseCallback func(client ClientInterface, code int, text string)

type ClientCallbackOption func(callBack *ClientCallback)

type ClientCallback struct {
	openCallBack    OpenCallback
	messageCallBack MessageCallback
	closeCallBack   CloseCallback
}

func NewClientCallback(opts ...ClientCallbackOption) *ClientCallback {

	o := &ClientCallback{}

	for _, opt := range opts {
		opt(o)
	}

	return o
}

func (c *ClientCallback) Open(client ClientInterface) {
	if c.openCallBack != nil {
		c.openCallBack(client)
	}
}

func (c *ClientCallback) Message(message *ReceiveContent) {
	if c.messageCallBack != nil {
		c.messageCallBack(message)
	}
}

func (c *ClientCallback) Close(client ClientInterface, code int, text string) {
	if c.closeCallBack != nil {
		c.closeCallBack(client, code, text)
	}
}

func WithOpenCallback(call OpenCallback) ClientCallbackOption {
	return func(callBack *ClientCallback) {
		callBack.openCallBack = call
	}
}

func WithMessageCallback(call MessageCallback) ClientCallbackOption {
	return func(callBack *ClientCallback) {
		callBack.messageCallBack = call
	}
}

func WithCloseCallback(call CloseCallback) ClientCallbackOption {
	return func(callBack *ClientCallback) {
		callBack.closeCallBack = call
	}
}
