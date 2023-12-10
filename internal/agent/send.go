package agent

import (
	"fmt"
	"github.com/KirillKhitev/metrics/internal/config"
	"github.com/go-resty/resty/v2"
)

func SendUpdate(client *resty.Client, t, name, value string) (*resty.Response, error) {
	url := fmt.Sprintf("http://%s%s/update/%s/%s/%s", config.ServerHost, config.ServerPort, t, name, value)
	resp, err := client.R().Post(url)

	return resp, err
}
