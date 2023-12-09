package agent

import (
	"fmt"
	"github.com/KirillKhitev/metrics/internal/config"
	"net/http"
)

func SendUpdate(t, name, value string) (*http.Response, error) {
	url := fmt.Sprintf("http://%s%s/update/%s/%s/%s", config.ServerHost, config.ServerPort, t, name, value)
	resp, err := http.Post(url, "text/plain", nil)

	return resp, err
}
