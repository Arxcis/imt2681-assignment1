package main

import (
    "net/http"
    "log"     
    "github.com/gorilla/mux"
    "fmt"
    "io/ioutil"
    "encoding/json"
    "sync" 
    "os"
)
const BASE = "https://api.github.com"

type GithubOwner struct {
    Name string `json:"login"`
    Id   uint   `json:"id"`
}

type GithubContributor struct {
    Name          string `json:"login"`
    Contributions uint   `json:"contributions"`    
}

type GithubRepo struct {
    Name         string      `json:"name"`
    Owner        GithubOwner `json:"owner"`
    Contributors []GithubContributor
    Languages    map[string]uint
}

type ServiceResponse struct {
    Repository string   `json:"repository"`
    Owner      string   `json:"owner"`
    Committer  string   `json:"commiter"`
    Commits    uint     `json:"commits"`
    Languages  []string `json:"languages"`
}

func getRequestBody(url string) []byte {
    resp, err := http.Get(url);
    if err != nil {
        log.Println(err)
    }

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        log.Println(err)
    }
    resp.Body.Close()
    return body
}

func getProjectinfo(user string, repo string) string {
    CURRENT_REPO                  := user + "/" + repo
    CURRENT_REPO_URL              := BASE + "/repos/" + CURRENT_REPO
    CURRENT_REPO_LANGUAGES_URL    := BASE + "/repos/" + CURRENT_REPO + "/languages"
    CURRENT_REPO_CONTRIBUTORS_URL := BASE + "/repos/" + CURRENT_REPO + "/contributors"

    // -1. Create a go-routine waitgroup
    var wg sync.WaitGroup
    wg.Add(3)
    
    // 0. Create global data
    githubRepo      := GithubRepo{}
    serviceResponse := ServiceResponse{}

    // 1. repos/<username>/<reponame>
    go func() { 
        defer wg.Done()

        body := getRequestBody(CURRENT_REPO_URL)    
        if err := json.Unmarshal(body, &githubRepo); err != nil {
            log.Println(string(body))
            log.Println(err)
        }
        serviceResponse.Repository = githubRepo.Name // @TODO - This should be on the form github.com/user/reponame
        serviceResponse.Owner      = githubRepo.Owner.Name

    }()

    // 2. repos/<username>/<reponame>/languages
    go func() {
        defer wg.Done()

        body := getRequestBody(CURRENT_REPO_LANGUAGES_URL)
        if err := json.Unmarshal(body, &githubRepo.Languages); err != nil {
            log.Println(string(body))
            log.Println(err)
        }   

        for k,_ := range githubRepo.Languages {
            serviceResponse.Languages = append(serviceResponse.Languages, k)
        }
    }()

    // 3. repos/<username>/<reponame>/contributors
    go func() {
        defer wg.Done()
       
        body := getRequestBody(CURRENT_REPO_CONTRIBUTORS_URL)
        if err := json.Unmarshal(body, &githubRepo.Contributors); err != nil {
            log.Println(string(body))
            log.Println(err)
        }

        // 4. Find the most valuable contributor
        maxContributor := GithubContributor{}
        for _, c := range githubRepo.Contributors {
            if maxContributor.Contributions < c.Contributions {
                maxContributor = c
            } 
        }
        serviceResponse.Committer = maxContributor.Name
        serviceResponse.Commits   = maxContributor.Contributions
    }()

    wg.Wait()

    // 5. Return to user
    data, err := json.Marshal(serviceResponse)
    if err != nil {
        log.Println(err)
    }
    log.Println(string(data))
    return string(data)
}


func DefaultHandler(w http.ResponseWriter, r *http.Request) {

    defaultBase := "projectinfo/v1/github.com/"
    defaultUser := "FylkesmannenIKT"
    defaultRepo := "OrdBase"

    fmt.Fprintln(w, "Url:  " + r.URL.Path)
    fmt.Fprintln(w, "Base: " + defaultBase)
    fmt.Fprintln(w, "User: " + defaultUser)
    fmt.Fprintln(w, "Repo: " + defaultRepo)

    fmt.Fprintln(w, "----------------- ")
    fmt.Fprintln(w, "Projectinfo: " + getProjectinfo(defaultUser, defaultRepo))
}

func GithubProjectinfo(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "Github - projectinfo")

    vars := mux.Vars(r)
    fmt.Fprintln(w, "Url:  " + r.URL.Path)
    fmt.Fprintln(w, "Base: " + vars["base"])
    fmt.Fprintln(w, "User: " + vars["user"])
    fmt.Fprintln(w, "Repo: " + vars["repo"])

    fmt.Fprintln(w, "----------------- ")
    fmt.Fprintln(w, "Projectinfo: " + getProjectinfo(vars["user"], vars["repo"]))
}

//
// MAIN
//
func main() {
    router := mux.NewRouter().StrictSlash(true)

    router.HandleFunc("/projectinfo/v1/{base}/{user}/{repo}", GithubProjectinfo)
    router.HandleFunc("/", DefaultHandler)
    log.Fatal(http.ListenAndServe(":" + os.Getenv("PORT"), router))
}
