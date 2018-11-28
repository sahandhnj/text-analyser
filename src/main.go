package main

import (
	"fmt"
	"log"
	"net/http"
)

const Address = ":3005"

func main() {
	http.HandleFunc("/analyse", analyseRequest)

	fmt.Println("Listening on " + Address)
	err := http.ListenAndServe(Address, nil)
	if err != nil {
		log.Fatal(err)
	}
}
