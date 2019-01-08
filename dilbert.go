package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
)

const (
	dilbertHome     = "http://dilbert.com"
	imgClass        = "img-comic"
	fileNamePattern = "dilbert-%d.gif"
)

func main() {
	body, err := httpGet(dilbertHome)
	if err != nil {
		fmt.Fprintf(os.Stderr, "httpGet %s: %v\n", dilbertHome, err)
		os.Exit(1)
	}
	defer body.Close()
	doc, err := html.Parse(body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "parse response: %v\n", err)
		os.Exit(1)
	}
	imageURLs := extractImages(doc, imgClass)
	for i, imgURL := range imageURLs {
		fileName := fmt.Sprintf(fileNamePattern, i)
		f, err := os.Create(fileName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "open file %s: %v", fileName, err)
		}
		defer f.Close()
		imgBody, err := httpGet(imgURL)
		if err != nil {
			fmt.Fprintf(os.Stderr, "open file %s: %v", imgURL, err)
		}
		io.Copy(f, imgBody)
	}
}

func httpGet(url string) (io.ReadCloser, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("GET %s: %v", url, err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET %s: %s", url, resp.Status)
	}
	return resp.Body, nil
}

func extractImages(n *html.Node, class string) []string {
	var imageURLs []string
	if n.Type == html.ElementNode && n.Data == "img" {
		var imgURL string
		var add bool
		for _, attr := range n.Attr {
			if attr.Key == "class" {
				for _, v := range strings.Fields(attr.Val) {
					if v == class {
						add = true
					}
				}
			} else if attr.Key == "src" {
				imgURL = attr.Val
			}
		}
		if add && len(imgURL) > 0 {
			imageURLs = append(imageURLs, imgURL)
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		for _, imgURL := range extractImages(c, class) {
			if !(strings.HasPrefix(imgURL, "http:") || strings.HasPrefix(imgURL, "https:")) {
				imgURL = "http:" + imgURL
			}
			imageURLs = append(imageURLs, imgURL)
		}
	}
	return imageURLs
}
