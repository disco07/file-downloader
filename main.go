package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
)

func downloader(url string) error {
	client := http.Client{}

	res, err := http.Head(url)
	if err != nil {
		return err
	}
	urlSplit := strings.Split(url, "/")
	filename := urlSplit[len(urlSplit)-1]
	if res.Header.Get("Accept-Ranges") != "bytes" {
		return errors.New("impossible de télécharger ce fichier")
	}

	cntLen, err := strconv.Atoi(res.Header.Get("Content-Length"))
	if err != nil {
		return err
	}
	nbPart := 6
	offset := cntLen / nbPart

	wg := sync.WaitGroup{}

	for i := 0; i < nbPart; i++ {
		wg.Add(1)
		name := fmt.Sprintf("part%d.part", i)
		start := i * offset
		end := (i + 1) * offset

		go func() {
			defer wg.Done()
			part, err := os.Create(name)
			if err != nil {
				return
			}
			defer part.Close()

			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				return
			}

			req.Header.Add("Range", fmt.Sprintf("bytes=%d-%d", start, end))
			res, err := client.Do(req)
			if err != nil {
				return
			}
			defer res.Body.Close()

			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				return
			}
			_, err = part.Write(body)
			if err != nil {
				return
			}
		}()
	}

	wg.Wait()

	out, err := os.Create(filename)
	for i := 0; i < nbPart; i++ {
		name := fmt.Sprintf("part%d.part", i)
		file, err := ioutil.ReadFile(name)
		if err != nil {
			return err
		}
		out.WriteAt(file, int64(i*offset))

		if err := os.Remove(name); err != nil {
			return err
		}
	}

	return nil
}

func main() {
	var url string
	flag.StringVar(&url, "u", "https://agritrop.cirad.fr/584726/1/Rapport.pdf", "url of the file to download")
	err := downloader(url)
	if err != nil {
		log.Fatal(err)
		return
	}
}
