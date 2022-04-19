package service

import (
	"context"
	"encoding/json"
	"errors"
	"net/url"
	"strings"

	"go-chat/config"
	"go-chat/internal/pkg/client"
	"go-chat/internal/pkg/sliceutil"
)

type IpAddressService struct {
	*BaseService
	conf       *config.Config
	httpClient *client.HttpClient
}

func NewIpAddressService(baseService *BaseService, conf *config.Config, httpClient *client.HttpClient) *IpAddressService {
	return &IpAddressService{BaseService: baseService, conf: conf, httpClient: httpClient}
}

func (s *IpAddressService) getCache(ip string) (string, error) {
	return s.rds.HGet(context.Background(), "rds:hash:ip-address", ip).Result()
}

func (s *IpAddressService) setCache(ip string, value string) error {
	return s.rds.HSet(context.Background(), "rds:hash:ip-address", ip, value).Err()
}

type IpAddressResponse struct {
	Code   string `json:"resultcode"`
	Reason string `json:"reason"`
	Result struct {
		Country  string `json:"Country"`
		Province string `json:"Province"`
		City     string `json:"City"`
		Isp      string `json:"Isp"`
	} `json:"result"`
	ErrorCode int `json:"error_code"`
}

func (s *IpAddressService) FindAddress(ip string) (string, error) {
	if val, err := s.getCache(ip); err == nil {
		return val, nil
	}

	params := &url.Values{}
	params.Add("key", s.conf.App.JuheKey)
	params.Add("ip", ip)

	resp, err := s.httpClient.Get("http://apis.juhe.cn/ip/ipNew", params)
	if err != nil {
		return "", err
	}

	data := &IpAddressResponse{}
	if err := json.Unmarshal(resp, data); err != nil {
		return "", err
	}

	if data.Code != "200" {
		return "", errors.New(data.Reason)
	}

	arr := []string{data.Result.Country, data.Result.Province, data.Result.City, data.Result.Isp}
	val := strings.Join(sliceutil.UniqueString(arr), " ")
	val = strings.TrimSpace(val)

	_ = s.setCache(ip, val)

	return val, nil
}
