package main

import (
	"github.com/KirillKhitev/metrics/internal/server"
)

func main() {
	if err := server.Run(); err != nil {
		panic(err)
	}
}
