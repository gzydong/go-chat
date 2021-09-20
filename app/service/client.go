package service

import "go-chat/app/cache"

type ClientService struct {
}

func NewClientService() *ClientService {
	return new(ClientService)
}

// Bind ...
func (c *ClientService) Bind(channel string, uuid string, id int) {
	cache.NewWsClient().Set(channel, uuid, id)
}

// UnBind ...
func (c *ClientService) UnBind(channel string, uuid string) {
	cache.NewWsClient().Del(channel, uuid)
}
