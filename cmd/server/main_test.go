package main

import (
	"github.com/KirillKhitev/metrics/internal/server"
	"github.com/KirillKhitev/metrics/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func testRequest(t *testing.T, ts *httptest.Server, method, path string) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, nil)
	require.NoError(t, err)

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp, string(respBody)
}

func TestRouter(t *testing.T) {
	appStorage := storage.MemStorage{}
	if err := appStorage.Init(); err != nil {
		t.Errorf("dont init appStorage")
		return
	}

	ts := httptest.NewServer(server.GetRouter(&appStorage))
	defer ts.Close()

	var testTable = []struct {
		url    string
		status int
		method string
	}{
		{"/", http.StatusMethodNotAllowed, http.MethodPost},
		{"/", http.StatusOK, http.MethodGet},
		{"/other_command/", http.StatusNotFound, http.MethodPost},
		{"/update/other_type", http.StatusNotFound, http.MethodPost},
		{"/update/other_type/nameMetric", http.StatusBadRequest, http.MethodPost},
		{"/update/other_type/nameMetric/10", http.StatusBadRequest, http.MethodPost},
		{"/update/counter/", http.StatusNotFound, http.MethodPost},
		{"/update/counter/nameMetric/", http.StatusBadRequest, http.MethodPost},
		{"/update/counter/nameMetric/string", http.StatusBadRequest, http.MethodPost},
		{"/update/counter/nameMetric/10", http.StatusOK, http.MethodPost},
		{"/update/gauge/", http.StatusNotFound, http.MethodPost},
		{"/update/gauge/nameMetric/", http.StatusBadRequest, http.MethodPost},
		{"/update/gauge/nameMetric/string", http.StatusBadRequest, http.MethodPost},
		{"/update/gauge/nameMetric/10", http.StatusOK, http.MethodPost},
		{"/value/", http.StatusInternalServerError, http.MethodPost},
		{"/value/", http.StatusBadRequest, http.MethodGet},
		{"/value/wrong_type/", http.StatusNotFound, http.MethodGet},
		{"/value/wrong_type/Alloc", http.StatusBadRequest, http.MethodGet},
		{"/value/counter/", http.StatusNotFound, http.MethodGet},
		{"/value/counter/PollCount", http.StatusNotFound, http.MethodGet},
		{"/value/gauge/Alloc", http.StatusNotFound, http.MethodGet},
	}
	for _, tt := range testTable {
		t.Run(tt.url, func(t *testing.T) {
			resp, _ := testRequest(t, ts, tt.method, tt.url)
			defer resp.Body.Close()
			assert.Equal(t, tt.status, resp.StatusCode)
		})
	}
}
