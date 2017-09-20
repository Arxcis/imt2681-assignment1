package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
)

// UnmarshalJSONHttp ...
func UnmarshalJSONHttp(target interface{}, url string, wg *sync.WaitGroup, errorChannel chan error) {

	log.Println("Requesting URL: " + url)

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
	}
	wg.Done()
	return
}

// UnmarshalJSONFile ...
func UnmarshalJSONFile(target interface{}, filepath string, wg *sync.WaitGroup, errorChannel chan error) {

	log.Println("Reading FILE: " + filepath)

	file, err := ioutil.ReadFile(filepath)
	if err != nil {
		errorChannel <- err
		wg.Done()
		return
	}

	err = json.Unmarshal(file, target)
	if err != nil {
		errorChannel <- err
	}

	wg.Done()
	return
}

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
		Languages map[string]interface{}
	}

	githubRepo := &GitRepositoryIn{}
	{
		errorChannel := make(chan error)
		wg := &sync.WaitGroup{}
		wg.Add(3)

		devenv, _ := strconv.ParseBool(os.Getenv("DEVENV"))
		if devenv {
			repoFile := "json/repo.json"
			languagesFile := "json/languages.json"
			contributorsFile := "json/contributors.json"

			go UnmarshalJSONFile(githubRepo, repoFile, wg, errorChannel)
			go UnmarshalJSONFile(&(githubRepo.Languages), languagesFile, wg, errorChannel)
			go UnmarshalJSONFile(&(githubRepo.Contributors), contributorsFile, wg, errorChannel)

		} else {
			repoURL := "https://api.github.com/repos/" + user + "/" + repo
			languagesURL := repoURL + "/languages"
			contributorsURL := repoURL + "/contributors"

			go UnmarshalJSONHttp(githubRepo, repoURL, wg, errorChannel)
			go UnmarshalJSONHttp(&(githubRepo.Languages), languagesURL, wg, errorChannel)
			go UnmarshalJSONHttp(&(githubRepo.Contributors), contributorsURL, wg, errorChannel)
		}
		wg.Wait()

		close(errorChannel)
		for err := range errorChannel {
			return nil, err
		}
	}

	type GitRepositoryOut struct {
		Repository string   `json:"repository"` // e.g. Ordbase
		Owner      string   `json:"owner"`      // e.g. FylkesmannenIKT
		Committer  string   `json:"committer"`  // e.g. Arxcis
		Commits    uint     `json:"commits"`    // e.g. 115
		Languages  []string `json:"languages"`  // e.g. [shell, java, scala, ...]
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

// BadRequestHandler ...
func BadRequestHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
}

// GitRepositoryHandler ...
func GitRepositoryHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "application/json")
	vars := mux.Vars(r)
	parsedRepository, err := ParseGitRepository(vars["user"], vars["repo"])

	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(parsedRepository)
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/projectinfo/v1/github.com/", BadRequestHandler)
	router.HandleFunc("/projectinfo/v1/github.com/{user}", BadRequestHandler)
	router.HandleFunc("/projectinfo/v1/github.com/{user}/{repo}", GitRepositoryHandler)
	log.Println(http.ListenAndServe(":"+os.Getenv("PORT"), router))
}
