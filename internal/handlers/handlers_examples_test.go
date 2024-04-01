package handlers

import (
	"fmt"
	"io"
	"net/http/httptest"
	"strings"

	"github.com/KirillKhitev/metrics/internal/storage"
)

func ExampleGetHandler_ServeHTTP() {
	var appStorage = storage.MemStorage{
		Counter: map[string]int64{
			"PollCount": 100,
			"Countet1":  20,
		},
		Gauge: map[string]float64{
			"Alloc": 324.25,
			"Mem":   23450.00,
			"CPU":   600.50,
		},
	}

	ch := &GetHandler{
		Storage: &appStorage,
	}

	request := httptest.NewRequest("GET", "/value/counter/PollCount", nil)
	w := httptest.NewRecorder()

	ch.ServeHTTP(w, request)

	res := w.Result()
	defer res.Body.Close()

	resBody, _ := io.ReadAll(res.Body)

	fmt.Println(res.StatusCode)
	fmt.Println(string(resBody))

	// Output:
	// 200
	// 100
}

func ExampleGetJSONHandler_ServeHTTP() {
	var appStorage = storage.MemStorage{
		Counter: map[string]int64{
			"PollCount": 100,
			"Countet1":  20,
		},
		Gauge: map[string]float64{
			"Alloc": 324.25,
			"Mem":   23450.00,
			"CPU":   600.50,
		},
	}

	ch := &GetJSONHandler{
		Storage: &appStorage,
	}

	s := `{"id":"PollCount","type":"counter"}`

	buf := strings.NewReader(s)

	request := httptest.NewRequest("POST", "/value", buf)
	w := httptest.NewRecorder()

	ch.ServeHTTP(w, request)

	res := w.Result()
	defer res.Body.Close()

	resBody, _ := io.ReadAll(res.Body)

	fmt.Println(res.StatusCode)
	fmt.Println(string(resBody))

	// Output:
	// 200
	// {
	//     "id": "PollCount",
	//     "type": "counter",
	//     "delta": 100
	// }
}

func ExampleListHandler_ServeHTTP() {
	var appStorage = storage.MemStorage{
		Counter: map[string]int64{
			"PollCount": 100,
			"Countet1":  20,
		},
		Gauge: map[string]float64{
			"Alloc": 324.25,
			"Mem":   23450.00,
			"CPU":   600.50,
		},
	}

	ch := &ListHandler{
		Storage: &appStorage,
	}

	request := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	ch.ServeHTTP(w, request)

	res := w.Result()
	defer res.Body.Close()

	resBody, _ := io.ReadAll(res.Body)

	fmt.Println(res.StatusCode)
	fmt.Println(len(resBody))

	// Output:
	// 200
	// 142
}

func ExampleUpdateHandler_ServeHTTP() {
	var appStorage = storage.MemStorage{
		Counter: map[string]int64{
			"PollCount": 100,
			"Countet1":  20,
		},
		Gauge: map[string]float64{
			"Alloc": 324.25,
			"Mem":   23450.00,
			"CPU":   600.50,
		},
	}

	ch := &UpdateHandler{
		Storage: &appStorage,
	}

	request := httptest.NewRequest("POST", "/update/counter/PollCount/200", nil)
	w := httptest.NewRecorder()

	ch.ServeHTTP(w, request)

	res := w.Result()
	defer res.Body.Close()

	fmt.Println(res.StatusCode)

	// Output:
	// 200
}

func ExampleUpdateJSONHandler_ServeHTTP() {
	var appStorage = storage.MemStorage{
		Counter: map[string]int64{
			"PollCount": 100,
			"Countet1":  20,
		},
		Gauge: map[string]float64{
			"Alloc": 324.25,
			"Mem":   23450.00,
			"CPU":   600.50,
		},
	}

	ch := &UpdateJSONHandler{
		Storage: &appStorage,
	}

	s := `{"id":"PollCount","type":"counter","delta":150}`

	buf := strings.NewReader(s)

	request := httptest.NewRequest("POST", "/update", buf)
	w := httptest.NewRecorder()

	ch.ServeHTTP(w, request)

	res := w.Result()
	defer res.Body.Close()

	resBody, _ := io.ReadAll(res.Body)

	fmt.Println(res.StatusCode)
	fmt.Println(string(resBody))

	// Output:
	// 200
	// {
	//     "id": "PollCount",
	//     "type": "counter",
	//     "delta": 250
	// }
}

func ExampleUpdatesHandler_ServeHTTP() {
	var appStorage = storage.MemStorage{
		Counter: map[string]int64{
			"PollCount": 100,
			"Countet1":  20,
		},
		Gauge: map[string]float64{
			"Alloc": 324.25,
			"Mem":   23450.00,
			"CPU":   600.50,
		},
	}

	ch := &UpdatesHandler{
		Storage: &appStorage,
	}

	s := `[
    {
        "id": "Alloc",
        "type": "gauge",
        "value": 11.55
    },
    {
        "id": "PollCount",
        "type": "counter",
        "delta": 150
    }
]`

	buf := strings.NewReader(s)

	request := httptest.NewRequest("POST", "/updates", buf)
	w := httptest.NewRecorder()

	ch.ServeHTTP(w, request)

	res := w.Result()
	defer res.Body.Close()

	resBody, _ := io.ReadAll(res.Body)

	fmt.Println(res.StatusCode)
	fmt.Println(string(resBody))

	// Output:
	// 200
	// [
	//     {
	//         "id": "PollCount",
	//         "type": "counter",
	//         "delta": 250
	//     },
	//     {
	//         "id": "Alloc",
	//         "type": "gauges",
	//         "value": 11.55
	//     }
	// ]
}
