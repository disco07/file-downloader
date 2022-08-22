package main

import (
	"fmt"
	"log"
	"net/http"
)

type app struct {
	client http.Client
}

func main() {
	//url := os.Args[1]
	//urlSplit := strings.Split(url, "/")
	//filename := urlSplit[len(urlSplit)-1]
	res, err := http.Head("https://agritrop.cirad.fr/584726/1/Rapport.pdf")
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println(res.Header.Get("Accept-Ranges"))
	if res.Header.Get("Accept-Ranges") == "bytes" {

	}
}
