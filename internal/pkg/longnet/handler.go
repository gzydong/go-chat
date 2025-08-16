package longnet

import (
	"log"
)

var _ IHandler = (*Handler)(nil)

type HandlerOption func(r *Handler)

type Handler struct {
	open    func(smg ISessionManager, c ISession)
	message func(smg ISessionManager, c ISession, message []byte)
	close   func(cid int64, uid int64)
}

func NewHandler(opts ...HandlerOption) IHandler {
	r := &Handler{}

	for _, opt := range opts {
		opt(r)
	}

	return r
}

func (r *Handler) OnOpen(smg ISessionManager, c ISession) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("panic:", err)
		}
	}()

	if r.open != nil {
		r.open(smg, c)
	}
}

func (r *Handler) OnMessage(smg ISessionManager, c ISession, message []byte) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("panic:", err)
		}
	}()

	if r.message != nil {
		r.message(smg, c, message)
	}
}

func (r *Handler) OnClose(cid int64, uid int64) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("panic:", err)
		}
	}()

	if r.close != nil {
		r.close(cid, uid)
	}
}

func WithOpenHandler(fn func(smg ISessionManager, c ISession)) HandlerOption {
	return func(r *Handler) {
		r.open = fn
	}
}

func WithMessageHandler(fn func(smg ISessionManager, c ISession, data []byte)) HandlerOption {
	return func(r *Handler) {
		r.message = fn
	}
}

func WithCloseHandler(fn func(cid int64, uid int64)) HandlerOption {
	return func(r *Handler) {
		r.close = fn
	}
}
