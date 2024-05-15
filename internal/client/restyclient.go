package client

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/KirillKhitev/metrics/internal/flags"
	"github.com/KirillKhitev/metrics/internal/metrics"
	"github.com/KirillKhitev/metrics/internal/mycrypto"
	"github.com/KirillKhitev/metrics/internal/signature"
	"github.com/KirillKhitev/metrics/internal/subnet"
)

type RestyClient struct {
	client *resty.Client
}

func NewRestyClient() (*RestyClient, error) {
	restyClient := &RestyClient{
		client: resty.New(),
	}

	return restyClient, nil
}

func (c *RestyClient) Send(ctx context.Context, metricsData []metrics.Metrics) error {
	data, err := c.prepareDataForSend(metricsData)

	if err != nil {
		return err
	}

	url := fmt.Sprintf("http://%s/updates/", flags.ArgsClient.AddrRun)

	ip, err := subnet.GetIP()
	if err != nil {
		return err
	}

	ctxt, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	request := c.client.NewRequest().
		SetContext(ctxt).
		SetBody(data).
		SetHeader("Content-Encoding", "gzip").
		SetHeader("X-Real-IP", ip)

	if flags.ArgsClient.Key != "" {
		hashSum := signature.GetHash(data, flags.ArgsClient.Key)
		request.SetHeader("HashSHA256", hashSum)
	}

	_, err = request.Post(url)

	return err
}

func (c *RestyClient) prepareDataForSend(data []metrics.Metrics) ([]byte, error) {
	dataJSON, err := json.Marshal(data)
	if err != nil {
		return dataJSON, fmt.Errorf("error by json-encode metrics: %w", err)
	}

	dataEncrypted, err := mycrypto.Encrypt(dataJSON, flags.ArgsClient.CryptoKey)
	if err != nil {
		return dataEncrypted, fmt.Errorf("error by encrypting data: %s, err: %w", dataJSON, err)
	}

	dataCompress, err := c.Compress(dataEncrypted)
	if err != nil {
		return dataCompress, fmt.Errorf("error by compress data: %s, err: %w", dataJSON, err)
	}

	return dataCompress, nil
}

// Compress сжимает данные перед отправкой на сервер.
func (c *RestyClient) Compress(data []byte) ([]byte, error) {
	var b bytes.Buffer

	w := gzip.NewWriter(&b)

	_, err := w.Write(data)
	if err != nil {
		return nil, fmt.Errorf("failed write data to compress temporary buffer: %v", err)
	}

	err = w.Close()
	if err != nil {
		return nil, fmt.Errorf("failed compress data: %v", err)
	}

	return b.Bytes(), nil
}

func (c *RestyClient) Close() error {
	return nil
}
