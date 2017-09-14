package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/gorilla/mux"
)

// GithubRepo ...
type GithubRepo struct {
	Name string `json:"name"`

	Owner struct {
		Name string `json:"login"`
		ID   uint   `json:"id"`
	} `json:"owner"`

	Contributors []struct {
		Name          string `json:"login"`
		Contributions uint   `json:"contributions"`
	}

	Languages map[string]uint
}

// ServiceResponse ...
type ServiceResponse struct {
	Repository string   `json:"repository"` // e.g. Ordbase
	Owner      string   `json:"owner"`      // e.g. FylkesmannenIKT
	Committer  string   `json:"committer"`  // e.g. Arxcis
	Commits    uint     `json:"commits"`    // e.g. 115
	Languages  []string `json:"languages"`  // e.g. [shell, java, scala, ...]
}

func parseJSON(target interface{}, url string) {

	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()
	json.NewDecoder(resp.Body).Decode(target)
}

func getProjectinfo(user string, repo string) string {

	githubRepo := GithubRepo{}
	{
		wg := sync.WaitGroup{}
		wg.Add(3)

		repoURL := "https://api.github.com/repos/" + user + "/" + repo
		go func() {
			defer wg.Done()

			parseJSON(githubRepo, repoURL)
		}()

		languagesURL := repoURL + "/languages"
		go func() {
			defer wg.Done()
			parseJSON(githubRepo.Languages, languagesURL)
		}()

		contributorsURL := repoURL + "/contributors"
		go func() {
			defer wg.Done()
			parseJSON(githubRepo.Contributors, contributorsURL)
		}()
		wg.Wait()
	}
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {

	defaultBase := "projectinfo/v1/github.com/"
	defaultUser := "FylkesmannenIKT"
	defaultRepo := "OrdBase"

	fmt.Fprintln(w, "Url:  "+r.URL.Path)
	fmt.Fprintln(w, "Base: "+defaultBase)
	fmt.Fprintln(w, "User: "+defaultUser)
	fmt.Fprintln(w, "Repo: "+defaultRepo)

	fmt.Fprintln(w, "----------------- ")
	fmt.Fprintln(w, "Projectinfo: "+getProjectinfo(defaultUser, defaultRepo))
}

func githubProjectinfo(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Github - projectinfo")

	vars := mux.Vars(r)
	fmt.Fprintln(w, "Url:  "+r.URL.Path)
	fmt.Fprintln(w, "Base: "+vars["base"])
	fmt.Fprintln(w, "User: "+vars["user"])
	fmt.Fprintln(w, "Repo: "+vars["repo"])

	fmt.Fprintln(w, "----------------- ")
	fmt.Fprintln(w, "Projectinfo: "+getProjectinfo(vars["user"], vars["repo"]))
}

//
// MAIN
//
func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", defaultHandler)
	router.HandleFunc("/projectinfo/v1/{base}/{user}/{repo}", githubProjectinfo)
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), router))
}
