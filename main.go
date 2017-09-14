package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/gorilla/mux"
)

// ParseGitRepository ...
func ParseGitRepository(user string, repo string) (interface{}, error) {

	type GitRepositoryIn struct {
		Name  string `json:"name"`
		Owner struct {
			Name string `json:"login"`
		} `json:"owner"`

		Contributors []struct {
			Name          string `json:"login"`
			Contributions uint   `json:"contributions"`
		}
		Languages map[string]uint
	}

	type GitRepositoryOut struct {
		Repository string   `json:"repository"` // e.g. Ordbase
		Owner      string   `json:"owner"`      // e.g. FylkesmannenIKT
		Committer  string   `json:"committer"`  // e.g. Arxcis
		Commits    uint     `json:"commits"`    // e.g. 115
		Languages  []string `json:"languages"`  // e.g. [shell, java, scala, ...]
	}

	githubRepo := &GitRepositoryIn{}
	{
		UnmarshalJSON := func(target interface{}, url string, wg *sync.WaitGroup, errorChannel chan error) {

			resp, err := http.Get(url)
			if err != nil {
				errorChannel <- err
				wg.Done()
				return
			}
			defer resp.Body.Close() // Do I even have to close?

			err = json.NewDecoder(resp.Body).Decode(target)
			if err != nil {
				errorChannel <- err
				wg.Done()
				return
			}
			wg.Done()
			return
		}

		repoURL := "https://api.github.com/repos/" + user + "/" + repo
		languagesURL := repoURL + "/languages"
		contributorsURL := repoURL + "/contributors"

		errorChannel := make(chan error)
		wg := &sync.WaitGroup{}
		wg.Add(3)

		go UnmarshalJSON(githubRepo, repoURL, wg, errorChannel)
		go UnmarshalJSON(&(githubRepo.Languages), languagesURL, wg, errorChannel)
		go UnmarshalJSON(&(githubRepo.Contributors), contributorsURL, wg, errorChannel)

		wg.Wait()

		close(errorChannel)
		for err := range errorChannel {
			return nil, err
		}
	}

	return GitRepositoryOut{
		Repository: githubRepo.Name,
		Owner:      githubRepo.Owner.Name,
		Committer:  githubRepo.Contributors[0].Name,
		Commits:    githubRepo.Contributors[0].Contributions,
		Languages: (func() []string {

			langs := make([]string, 0, len(githubRepo.Languages))
			for key := range githubRepo.Languages {
				langs = append(langs, key)
			}
			return langs
		}()),
	}, nil
}

func badRequestHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
}

func service(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "application/json")
	vars := mux.Vars(r)
	info, err := ParseGitRepository(vars["user"], vars["repo"])

	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(info)
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/projectinfo/v1/github.com/", badRequestHandler)
	router.HandleFunc("/projectinfo/v1/github.com/{user}", badRequestHandler)
	router.HandleFunc("/projectinfo/v1/github.com/{user}/{repo}", service)
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), router))
}
