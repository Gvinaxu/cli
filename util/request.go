package util

import (
	"time"

	"github.com/valyala/fasthttp"
)

func DoTimeout(arg *fasthttp.Args, method string, requestURI string, heads map[string]interface{}) (statusCode int, body []byte, err error) {
	req := &fasthttp.Request{}
	switch method {
	case "GET":
		req.Header.SetMethod(method)
		if arg != nil {
			requestURI = requestURI + "?" + arg.String()
		}

	case "POST":
		req.Header.SetMethod(method)
		if arg != nil {
			arg.WriteTo(req.BodyWriter())
		}

	}
	if heads != nil {
		for k, v := range heads {
			req.Header.Add(k, v.(string))
		}
	}
	req.SetRequestURI(requestURI)

	resp := &fasthttp.Response{}
	err = fasthttp.DoTimeout(req, resp, time.Second*30)

	return resp.StatusCode(), resp.Body(), err
}

func DoJsonTimeout(method string, url, bodyjson string, heads map[string]interface{}) (statusCode int, body []byte, err error) {
	req := &fasthttp.Request{}
	resp := &fasthttp.Response{}

	switch method {
	case "GET":
		req.Header.SetMethod(method)
	case "POST":
		req.Header.SetMethod(method)
	}
	if heads != nil {
		for k, v := range heads {
			req.Header.Add(k, v.(string))
		}
	}
	req.Header.SetContentType("application/json")
	req.SetBodyString(bodyjson)

	req.SetRequestURI(url)

	err = fasthttp.DoTimeout(req, resp, time.Second*30)
	return resp.StatusCode(), resp.Body(), err
}
