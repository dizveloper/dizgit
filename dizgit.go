package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type CliInput struct {
	gitUrl       string
	releases     bool
	pullRequests bool
}

type Pull struct {
	Title  string `json:"title"`
	Number int    `json:"number"`
}

type Release struct {
	Name    string `json:"name"`
	Version string `json:"tag_name"`
}

func getInfo(ghUrl string) []byte {
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	req, err := http.NewRequest("GET", ghUrl, nil)
	req.Header.Add("Accept", "application/vnd.github.v3+json")

	if err != nil {
		log.Fatal(err)
	}

	resp, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	return body
}

func getReleases(ghUrl string) []Release {
	repoCutout := parseGhUrlForRepo(ghUrl) // owner/repo
	listReleasesApi := "https://api.github.com/repos/" + repoCutout + "/releases"

	var resp = []Release{}
	json.Unmarshal(getInfo(listReleasesApi), &resp)

	if len(resp) > 3 {
		return resp[:3]
	} else {
		return resp
	}
}

func getPullRequests(ghUrl string) []Pull {
	repoCutout := parseGhUrlForRepo(ghUrl) // owner/repo
	listReleasesApi := "https://api.github.com/repos/" + repoCutout + "/pulls"

	var resp = []Pull{}
	json.Unmarshal(getInfo(listReleasesApi), &resp)

	if len(resp) > 3 {
		return resp[:3]
	} else {
		return resp
	}
}

func parseGhUrlForRepo(ghUrl string) string {
	splitter := strings.Split(ghUrl, "/")

	return splitter[len(splitter)-2] + "/" + strings.ReplaceAll(splitter[len(splitter)-1], ".git", "")
}

func main() {
	var releases = flag.Bool("releases", false, "Get Latest GitHub Repository Releases")
	var pullRequests = flag.Bool("pullrequests", false, "Get Latest GitHub Repository Pull Requests")

	var repo = flag.String("repo", "", "GitHub Repository URL")

	flag.Parse()

	if len(os.Args) < 2 {
		fmt.Println("No arguments were passed")
		flag.PrintDefaults()
		os.Exit(0)
	}

	if len(*repo) == 0 || (!*releases && !*pullRequests) {
		/*
			Repo and at least one of these flags is required
		*/
		fmt.Println("")
		fmt.Println("REQUIRED: '--repository GIT_URL' along with '--pullrequests' and/or '--release'")
		fmt.Println("")

		flag.PrintDefaults()
		os.Exit(1)
	}

	cliArgs := CliInput{
		gitUrl:       *repo,
		releases:     *releases,
		pullRequests: *pullRequests,
	}

	/*

		This GitHub URL regex was based on the one here: https://github.com/jonschlinkert/is-git-url/blob/master/test.js

		I modified it slightly to only work with GitHub hostname, and no longer accepts basic auth as part of URL and those must be passed as arguments

	*/

	gitUrlRegex := regexp.MustCompile("(?:git|ssh|https?|git@[-\\w.]+):(\\/\\/)?(github.com\\/.*?)(\\.git)(\\/?|\\#[-\\d\\w._]+?)$")
	validGitUrl := gitUrlRegex.Match([]byte(cliArgs.gitUrl))

	if validGitUrl {
		if cliArgs.releases {
			rels := getReleases(cliArgs.gitUrl)
			if len(rels) > 0 {
				fmt.Println("")
				fmt.Println("--- Last " + strconv.Itoa(len(rels)) + " Releases ---")
				for _, rel := range rels {
					fmt.Println("-")
					fmt.Println("Release Name: " + rel.Name)
					fmt.Println("Release Version: " + rel.Version)
				}
			} else {
				fmt.Println("")
				fmt.Println("--- The provided repository doesn't have any releases ---")
				fmt.Println("")
			}
		}

		if cliArgs.pullRequests {
			prs := getPullRequests(cliArgs.gitUrl)
			if len(prs) > 0 {
				fmt.Println("")
				fmt.Println("--- Last " + strconv.Itoa(len(prs)) + " Pull Requests ---")
				for _, pr := range prs {
					fmt.Println("-")
					fmt.Println("PR Title: " + pr.Title)
					fmt.Println("PR Number: " + strconv.Itoa(pr.Number))
				}
			} else {
				fmt.Println("")
				fmt.Println("--- The provided repository doesn't have any Pull Requests ---")
				fmt.Println("")
			}
		}
	} else {
		fmt.Println("The GitHub repository URL provided is not properly formatted. Example: https://github.com/user/repo.git")
		os.Exit(0)
	}
}
