package im

type ICallback interface {
	Open(client IClient)
	Message(client IClient, message []byte)
	Close(client IClient, code int, text string)
	Destroy(client IClient)
}

type (
	OpenCallback         func(client IClient)
	MessageCallback      func(client IClient, message []byte)
	CloseCallback        func(client IClient, code int, text string)
	DestroyCallback      func(client IClient)
	ClientCallbackOption func(callBack *ClientCallback)
)

type ClientCallback struct {
	open    OpenCallback
	message MessageCallback
	close   CloseCallback
	destroy DestroyCallback
}

func NewClientCallback(opts ...ClientCallbackOption) ICallback {

	o := &ClientCallback{}

	for _, opt := range opts {
		opt(o)
	}

	return o
}

func (c *ClientCallback) Open(client IClient) {
	if c.open != nil {
		c.open(client)
	}
}

func (c *ClientCallback) Message(client IClient, message []byte) {
	if c.message != nil {
		c.message(client, message)
	}
}

func (c *ClientCallback) Close(client IClient, code int, text string) {
	if c.close != nil {
		c.close(client, code, text)
	}
}

func (c *ClientCallback) Destroy(client IClient) {
	if c.destroy != nil {
		c.destroy(client)
	}
}

func WithOpenCallback(call OpenCallback) ClientCallbackOption {
	return func(callBack *ClientCallback) {
		callBack.open = call
	}
}

func WithMessageCallback(call MessageCallback) ClientCallbackOption {
	return func(callBack *ClientCallback) {
		callBack.message = call
	}
}

func WithCloseCallback(call CloseCallback) ClientCallbackOption {
	return func(callBack *ClientCallback) {
		callBack.close = call
	}
}

func WithDestroyCallback(call DestroyCallback) ClientCallbackOption {
	return func(callBack *ClientCallback) {
		callBack.destroy = call
	}
}
