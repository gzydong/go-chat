package im

type ICallback interface {
	Open(client IClient)
	Message(client IClient, message []byte)
	Close(client IClient, code int, text string)
}

type (
	OpenCallback         func(client IClient)
	MessageCallback      func(client IClient, message []byte)
	CloseCallback        func(client IClient, code int, text string)
	ClientCallbackOption func(callBack *ClientCallback)
)

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

func (c *ClientCallback) Open(client IClient) {
	if c.openCallBack != nil {
		c.openCallBack(client)
	}
}

func (c *ClientCallback) Message(client IClient, message []byte) {
	if c.messageCallBack != nil {
		c.messageCallBack(client, message)
	}
}

func (c *ClientCallback) Close(client IClient, code int, text string) {
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
