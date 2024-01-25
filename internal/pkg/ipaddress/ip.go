package ipaddress

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	client *http.Client
}

func NewClient(client *http.Client) *Client {
	return &Client{client}
}

func (c *Client) GetIpInfo(ctx context.Context, ip string) (string, error) {
	if address, err := c.findIpAddressByCSDN(ctx, ip); err == nil {
		return address, nil
	}

	return c.findIpAddressByBaiDu(ctx, ip)
}

func (c *Client) findIpAddressByCSDN(_ context.Context, ip string) (string, error) {
	var result struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
			Address string `json:"address"`
			Ip      string `json:"ip"`
		} `json:"data"`
	}

	resp, err := c.client.Get(fmt.Sprintf("https://searchplugin.csdn.net/api/v1/ip/get?ip=%s", ip))
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	if result.Code == 200 {
		return result.Data.Address, nil
	}

	return "", fmt.Errorf("获取IP信息失败")
}

func (c *Client) findIpAddressByBaiDu(_ context.Context, ip string) (string, error) {
	var result struct {
		Status string `json:"status"`
		Data   []struct {
			Fetchkey string `json:"fetchkey"`
			Location string `json:"location"`
		} `json:"data"`
	}

	resp, err := c.client.Get(fmt.Sprintf("https://opendata.baidu.com/api.php?query=[%s]&co=&resource_id=6006&oe=utf8", ip))
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	if result.Status == "0" {
		return result.Data[0].Location, nil
	}

	return "", fmt.Errorf("获取IP信息失败")
}
