package main

import (
    "net/http"
    "log"     
    "github.com/gorilla/mux"
    "fmt"
    "./githubapi"
)

//
// TimeHandler
//


func main() {
    router := mux.NewRouter().StrictSlash(true)

    router.HandleFunc("/projectinfo/v1/{base}/{user}/{repo}", GithubProjectinfo)
    log.Fatal(http.ListenAndServe(":8080", router))
}

func GithubProjectinfo(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "Github - projectinfo")

    vars := mux.Vars(r)
    fmt.Fprintln(w, "Url:  " + r.URL.Path)
    fmt.Fprintln(w, "Base: " + vars["base"])
    fmt.Fprintln(w, "User: " + vars["user"])
    fmt.Fprintln(w, "Repo: " + vars["repo"])

    fmt.Fprintln(w, "----------------- ")
    fmt.Fprintln(w, "Projectinfo: " + githubapi.GetProjectinfo(vars["user"], vars["repo"]))
}
