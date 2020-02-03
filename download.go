package tools

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
)

func DownloadFile(url, filename string) {
	r, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer func() { _ = r.Body.Close() }()

	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer func() { _ = f.Close() }()

	n, err := io.Copy(f, r.Body)
	fmt.Println(n, err)
}

type Reader struct {
	io.Reader
	Total   int64
	Current int64
}

func (r *Reader) Read(p []byte) (n int, err error) {
	n, err = r.Reader.Read(p)

	r.Current += int64(n)
	fmt.Printf("\r进度 %.2f%%", float64(r.Current*10000/r.Total)/100)

	return
}

func DownloadFileProgress(url, filename string) string {
	client := &http.Client{}
	var req *http.Request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	defer func() { _ = req.Body.Close() }()

	cookie1 := &http.Cookie{Name: "_identity-frontend", Value: "516d655404e7ce0bb9fdabca241faefe364d4ff92a58b35f06f2744da21f00b9a%3A2%3A%7Bi%3A0%3Bs%3A18%3A%22_identity-frontend%22%3Bi%3A1%3Bs%3A50%3A%22%5B1364574%2C%22ITaaiC2jP_to0DYbViG_alnyHpJthLar%22%2C86400%5D%22%3B%7D", HttpOnly: true}
	req.AddCookie(cookie1)
	resp, err := client.Do(req)
	uri := strings.Split(resp.Request.URL.RawQuery, "?")
	ext := path.Ext(uri[0])
	f, err := os.Create(filename + ext)
	if err != nil {
		panic(err)
	}
	defer func() { _ = f.Close() }()

	reader := &Reader{
		Reader: resp.Body,
		Total:  resp.ContentLength,
	}

	_, _ = io.Copy(f, reader)
	return ext
}
