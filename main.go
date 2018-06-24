package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/gocolly/colly"
)

func main() {
	ENDPOINT := os.Args[1]

	c := colly.NewCollector()
	c.WithTransport(&http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	})

	// Find and visit all links
	c.OnHTML("a.project[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		repositoryDownload := ENDPOINT + link + "/repository/archive.zip?ref=master"
		go DownloadFile(strings.Replace(link, "/", "-", -1)[1:]+".zip", repositoryDownload)
	})

	c.OnHTML("li.next a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		c.Visit(e.Request.AbsoluteURL(link))
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.Visit(ENDPOINT + "/explore/projects")
}

func DownloadFile(filepath string, url string) error {

	// Create the file
	out, err := os.Create("repos/" + filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
