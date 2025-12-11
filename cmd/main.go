package main

import (
	"github.com/informeai/casd/handlers"
)

func main() {
	if err := handlers.Execute(); err != nil {
		panic(err)
	}
}
