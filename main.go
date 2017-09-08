package main

import (
	"net/http"
    "log"     
    "time"
)

type TimeHandler struct {
  format string
}

func (th *TimeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  tm := time.Now().Format(th.format)
  w.Write([]byte("The time is: " + tm))
}

func main() {

	mux := http.NewServeMux()
    th :=  &TimeHandler{ format: time.RFC1123 }

    mux.Handle("/time", th)
    log.Println("Listening on port 3000...")
    http.ListenAndServe(":3000", mux)
}

// Sources:
// [1]Why do I not like any Golang URL Routers? - https://husobee.github.io/golang/url-router/2015/06/15/why-do-all-golang-url-routers-suck.html
// [2]A Recap of request handling in go - http://www.alexedwards.net/blog/a-recap-of-request-handling
// [3]Time format constants godoc - https://godoc.org/time#pkg-constants