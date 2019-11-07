package main

import (
    "log"
    "net/http"
    "net/http/httputil"
    "net/url"
    "github.com/gorilla/mux"
    "github.com/namsral/flag"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

type chainMapping struct {
    endpoints map[string]string
}

var chainMap = chainMapping{
    endpoints: make(map[string]string),
}

func (chainMap *chainMapping) createMap() {
    chainMap.endpoints["aca376f206b8fc25a6ed44dbdc66547c36c6c33e3a119ffbeaef943642f0e906"] = "https://eos.greymass.com"
}

func handler(w http.ResponseWriter, r *http.Request) {
    if host, ok := chainMap.endpoints[r.Header.Get("X-Chain-Id")]; ok {
        url, _ := url.Parse(host)
        proxy := httputil.NewSingleHostReverseProxy(url)
        r.URL.Host = url.Host
        r.URL.Scheme = url.Scheme
        r.Host = url.Host
        proxy.ServeHTTP(w, r)
    }
}

func main() {
    flag.Parse()
    log.SetFlags(0)

    chainMap.createMap()

    r := mux.NewRouter()
    r.PathPrefix("/").HandlerFunc(handler)
    http.Handle("/", r)

    log.Printf("Listening on %s", *addr)
    log.Fatal(http.ListenAndServe(*addr, nil))
}
