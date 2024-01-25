package nsqd

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
)

type NsqStatInfo struct {
	Version   string `json:"version"`
	Health    string `json:"health"`
	StartTime int    `json:"start_time"`
	Topics    []struct {
		TopicName string `json:"topic_name"`
		Channels  []struct {
			ChannelName   string `json:"channel_name"`
			Depth         int    `json:"depth"`
			BackendDepth  int    `json:"backend_depth"`
			InFlightCount int    `json:"in_flight_count"`
			DeferredCount int    `json:"deferred_count"`
			MessageCount  int    `json:"message_count"`
			RequeueCount  int    `json:"requeue_count"`
			TimeoutCount  int    `json:"timeout_count"`
			ClientCount   int    `json:"client_count"`
			Clients       []struct {
				ClientId                      string `json:"client_id"`
				Hostname                      string `json:"hostname"`
				Version                       string `json:"version"`
				RemoteAddress                 string `json:"remote_address"`
				State                         int    `json:"state"`
				ReadyCount                    int    `json:"ready_count"`
				InFlightCount                 int    `json:"in_flight_count"`
				MessageCount                  int    `json:"message_count"`
				FinishCount                   int    `json:"finish_count"`
				RequeueCount                  int    `json:"requeue_count"`
				ConnectTs                     int    `json:"connect_ts"`
				SampleRate                    int    `json:"sample_rate"`
				Deflate                       bool   `json:"deflate"`
				Snappy                        bool   `json:"snappy"`
				UserAgent                     string `json:"user_agent"`
				Tls                           bool   `json:"tls"`
				TlsCipherSuite                string `json:"tls_cipher_suite"`
				TlsVersion                    string `json:"tls_version"`
				TlsNegotiatedProtocol         string `json:"tls_negotiated_protocol"`
				TlsNegotiatedProtocolIsMutual bool   `json:"tls_negotiated_protocol_is_mutual"`
			} `json:"clients"`
			Paused               bool `json:"paused"`
			E2EProcessingLatency struct {
				Count       int         `json:"count"`
				Percentiles interface{} `json:"percentiles"`
			} `json:"e2e_processing_latency"`
		} `json:"channels"`
		Depth                int  `json:"depth"`
		BackendDepth         int  `json:"backend_depth"`
		MessageCount         int  `json:"message_count"`
		MessageBytes         int  `json:"message_bytes"`
		Paused               bool `json:"paused"`
		E2EProcessingLatency struct {
			Count       int         `json:"count"`
			Percentiles interface{} `json:"percentiles"`
		} `json:"e2e_processing_latency"`
	} `json:"topics"`
}

func doGet[T any](ctx context.Context, uri string) (*T, error) {
	req, _ := http.NewRequest("GET", uri, nil)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var resp T
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func doDelete[T any](ctx context.Context, uri string) (*T, error) {
	req, _ := http.NewRequest("DELETE", uri, nil)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var resp T
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
