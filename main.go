package main

import (
	"log"
	"net/http"
	"os"

	"github.com/elazarl/goproxy"
)

func isPreFlightRequest() goproxy.ReqConditionFunc {
	return func(req *http.Request, ctx *goproxy.ProxyCtx) bool {
		if req.Method == "OPTIONS" {
			if len(req.Header.Get("Origin")) > 0 {
				return true
			}
		}
		return false
	}
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage:", os.Args[0], `[interface]<:port>`)
	}

	proxy := goproxy.NewProxyHttpServer()
	// proxy.Verbose = true /* for debugging */
	proxy.OnRequest(isPreFlightRequest()).DoFunc(func(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		preFlightResp := goproxy.NewResponse(r, r.Header.Get("Content-Type"), 200, "")
		preFlightResp.Header.Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		preFlightResp.Header.Set("Access-Control-Allow-Headers", "content-type")
		preFlightResp.Header.Set("Access-Control-Allow-Methods", "GET, POST")
		return r, preFlightResp
	})

	log.Fatal(http.ListenAndServe(os.Args[1], proxy))
}
