package main

import (
	"my-web-template/internal/core/bootstrap"
)

func main() {
	if err := bootstrap.Run(); err != nil {
		panic(err)
	}
}
