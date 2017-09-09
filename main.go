package main

import (
    "encoding/json"
	"net/http"
    "log"     
    "time")

//
// TimeHandler
// 
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

