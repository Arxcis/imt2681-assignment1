package main

import(
    "net/http"
    "log"
    "io/ioutil"
    "encoding/json"
    "sync"
    "fmt"
)

const BASE                          = "https://api.github.com"
const CURRENT_REPO                  = "FylkesmannenIKT/ordbase"
const CURRENT_REPO_URL              = BASE + "/repos/" + CURRENT_REPO
const CURRENT_REPO_LANGUAGES_URL    = BASE + "/repos/" + CURRENT_REPO + "/languages"
const CURRENT_REPO_CONTRIBUTORS_URL = BASE + "/repos/" + CURRENT_REPO + "/contributors"

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
        panic(err)
    }

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        panic(err)
    }
    resp.Body.Close()
    return body
}

func main() {
    messages := make(chan string)
    var wg sync.WaitGroup
    wg.Add(3) // 3 Functions have to signal finish before wait will pass

    githubRepo      := GithubRepo{}
    serviceResponse := ServiceResponse{}

    // 1. repos/<username>/<reponame>
    go func() {
        defer wg.Done()

        body := getRequestBody(CURRENT_REPO_URL)    
        if err := json.Unmarshal(body, &githubRepo); err != nil {
            panic(err)
        }
        serviceResponse.Repository = githubRepo.Name // @TODO - This should be on the form github.com/user/reponame
        serviceResponse.Owner      = githubRepo.Owner.Name

        messages <- "Thread1: " + CURRENT_REPO_URL + " is done!\n" 
    }()

    // 2. repos/<username>/<reponame>/languages
    go func() {
        defer wg.Done()

        body := getRequestBody(CURRENT_REPO_LANGUAGES_URL)
        if err := json.Unmarshal(body, &githubRepo.Languages); err != nil {
            panic(err)
        }   

        for k,_ := range githubRepo.Languages {
            serviceResponse.Languages = append(serviceResponse.Languages, k)
        }
        messages <- "Thread2: " + CURRENT_REPO_LANGUAGES_URL + " is done!\n" 
    }()

    // 3. repos/<username>/<reponame>/contributors
    go func() {
        defer wg.Done()

        body := getRequestBody(CURRENT_REPO_CONTRIBUTORS_URL)
        if err := json.Unmarshal(body, &githubRepo.Contributors); err != nil {
            panic(err)
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

        messages <- "Thread3: " + CURRENT_REPO_CONTRIBUTORS_URL + " is done!\n" 

    }()
    go func() {
        for i := range messages {
            log.Println(i)
        }
    }()

    wg.Wait()

    // 5. Return to user
    data, err := json.Marshal(serviceResponse)
    if err != nil {
        panic(err)
    }
    log.Println(string(data))
    return
}