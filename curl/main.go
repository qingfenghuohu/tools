package curl

import (
	"crypto/tls"
	"fmt"
	"github.com/qingfenghuohu/tools"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type New struct {
	content    string
	proxyData  string
	headerData map[string]interface{}
	cookieData map[string]interface{}
}

func (do *New) Header(headerData map[string]interface{}) *New {
	do.headerData = headerData
	return do
}

func (do *New) Cookie(cookieData map[string]interface{}) *New {
	do.cookieData = cookieData
	return do
}

func (do *New) Proxy(content string) *New {
	do.proxyData = content
	return do
}

func (do *New) Data(contentData map[string]interface{}) *New {
	tmp := []string{}
	for i, v := range contentData {
		tmp = append(tmp, i+"="+tools.Obj2Str(v))
	}
	do.content = strings.Join(tmp, "&")
	return do
}

func (do *New) Post(Url string) string {
	var result string
	result = do.run("POST", Url)
	do.Clear()
	return result
}

func (do *New) run(Method string, Url string) string {
	client := &http.Client{}
	Method = strings.ToUpper(Method)
	req, err := http.NewRequest(Method, Url, strings.NewReader(do.content))
	if err != nil {
		panic(err)
	}
	utmp, err := url.Parse(Url)
	if err != nil {
		panic(err)
	}
	Tr := &http.Transport{}
	if utmp.Scheme == "https" {
		Tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}
	if Method == "POST" {
		//req.Header.Set("Content-Type", "multipart/form-data")
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	} else {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	do.setHeader(req)
	do.setCookie(req)
	do.setProxy(Tr)
	client.Transport = Tr
	resp, err := client.Do(req)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	return string(body)
}

func (do *New) setHeader(req *http.Request) {
	for i, v := range do.headerData {
		req.Header.Set(i, tools.Obj2Str(v))
	}
}

func (do *New) setCookie(req *http.Request) {
	tmp := []string{}
	for i, v := range do.cookieData {
		tmp = append(tmp, i+"="+tools.Obj2Str(v))
	}
	cookie := strings.Join(tmp, "&")
	req.Header.Set("Cookie", cookie)
}

func (do *New) setProxy(tr *http.Transport) {
	if do.proxyData != "" {
		tr.Proxy = func(_ *http.Request) (*url.URL, error) {
			return url.Parse(do.proxyData)
		}
	}
}

func (do *New) Get(Url string) string {
	var result string
	result = do.run("GET", Url)
	do.Clear()
	return result
}

func (do *New) Clear() {
	do = &New{}
}
