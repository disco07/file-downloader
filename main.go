package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/schollz/progressbar/v3"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

func downloader(url string) error {
	if strings.TrimSpace(url) == "" {
		return errors.New("invalid url")
	}
	client := http.Client{}

	res, err := http.Head(url)
	if err != nil {
		fmt.Printf("error: %s", err.Error())
		return errors.New("unsupported protocol scheme")
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
	nbPart := 3
	offset := cntLen / nbPart

	wg := sync.WaitGroup{}

	for i := 0; i < nbPart; i++ {
		wg.Add(1)
		name := fmt.Sprintf("part%d", i)
		start := i * offset
		end := (i + 1) * offset

		i := i

		go func() {
			defer wg.Done()
			part, err := os.Create(name)
			if err != nil {
				return
			}
			defer part.Close()

			req, err := http.NewRequest(http.MethodGet, url, nil)
			if err != nil {
				return
			}

			req.Header.Add("Range", fmt.Sprintf("bytes=%d-%d", start, end))
			res, err := client.Do(req)
			if err != nil {
				return
			}
			defer res.Body.Close()

			f, err := os.OpenFile(name, os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				return
			}
			defer f.Close()

			bar := progressbar.DefaultBytes(
				res.ContentLength,
				fmt.Sprintf("downloading-worker %d", i+1),
			)
			io.Copy(io.MultiWriter(f, bar), res.Body)

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
	if err != nil {
		return err
	}
	defer out.Close()
	for i := 0; i < nbPart; i++ {
		name := fmt.Sprintf("part%d", i)
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
	flag.Parse()
	start := time.Now()
	err := downloader(url)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println(time.Since(start))
}
