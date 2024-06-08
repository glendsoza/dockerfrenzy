package main

import (
	"crazydocker/pkg/api"
	"log"
)

func main() {
	if err := api.Run(); err != nil {
		log.Fatal(err)
	}
}
