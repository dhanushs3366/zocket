package main

import (
	"github.com/dhanushs3366/zocket/handler"
	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	h, err := handler.Init()

	if err != nil {
		panic(err)
	}
	err = h.Run(8080)

	panic(err)
}
