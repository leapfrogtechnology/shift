package github

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/leapfrogtechnology/shift/cli/utils/http"
)

var baseURL = "https://api.github.com"

var userURL = baseURL + "/user"
var organizationsURL = baseURL + "/user/orgs"
var createTokenURL = baseURL + "/authorizations"
var userReposURL = baseURL + "/user/repos?affiliation=owner"

// GitCredentials describe the credentials required for Github
type GitCredentials struct {
	Username string
	Password string
}

type createPersonalTokenBody struct {
	Note   string   `json:"note"`
	Scopes []string `json:"scopes"`
}

type personalTokenSuccess struct {
	Token string `json:"token"`
}

type organization struct {
	Login string `json:"login"`
}

type user struct {
	Login string `json:"login"`
}

type repo struct {
	Name     string `json:"name"`
	CloneURL string `json:"clone_url"`
}

// CreatePersonalToken creates a personal access token in github
func CreatePersonalToken(credentials *GitCredentials) (string, error) {
	fmt.Print("Connecting to Github... ")
	data := &createPersonalTokenBody{
		Note:   "Shift CLI:" + time.Now().String(),
		Scopes: []string{"repo", "read:org"},
	}

	basicAuthString := base64.StdEncoding.EncodeToString([]byte(credentials.Username + ":" + credentials.Password))
	jsonData, _ := json.Marshal(data)

	response := &personalTokenSuccess{}
	_, err := http.Client.R().
		SetHeader("Authorization", "basic "+basicAuthString).
		SetBody(jsonData).
		SetResult(response).
		Post(createTokenURL)

	if err != nil {
		fmt.Println(err)

		return "", err
	}

	return response.Token, nil
}

// FetchOrganizations fetches user's organizations
func FetchOrganizations(personalToken string) ([]string, error) {
	response := []organization{}

	_, err := http.Client.R().
		SetHeader("Authorization", "token "+personalToken).
		SetResult(&response).
		Get(organizationsURL)

	if err != nil {
		fmt.Println(err)

		return []string{}, err
	}

	organizations := []string{}
	for _, organization := range response {
		organizations = append(organizations, organization.Login)
	}

	return organizations, nil
}

// FetchUser fetches user's information
func FetchUser(personalToken string) (string, error) {
	response := user{}

	_, err := http.Client.R().
		SetHeader("Authorization", "token "+personalToken).
		SetResult(&response).
		Get(userURL)

	if err != nil {
		fmt.Println(err)

		return "", err
	}

	return response.Login, nil
}

// FetchOrgRepos fetches user's repositories
func FetchOrgRepos(personalToken string, organization string) ([]string, map[string]string, error) {
	response := []repo{}
	repos := []string{}
	repoURL := map[string]string{}

	_, err := http.Client.R().
		SetHeader("Authorization", "token "+personalToken).
		SetResult(&response).
		Get(baseURL + "/orgs/" + organization + "/repos")

	if err != nil {
		fmt.Println("Error:")
		fmt.Println(err)

		return repos, repoURL, err
	}

	for _, repo := range response {
		repos = append(repos, repo.Name)
		repoURL[repo.Name] = repo.CloneURL
	}

	return repos, repoURL, nil
}

// FetchUserRepos fetches user's repositories
func FetchUserRepos(personalToken string) ([]string, map[string]string, error) {
	response := []repo{}
	repos := []string{}
	repoURL := map[string]string{}

	_, err := http.Client.R().
		SetHeader("Authorization", "token "+personalToken).
		SetResult(&response).
		Get(userReposURL)

	if err != nil {
		fmt.Println("Error:")
		fmt.Println(err)

		return repos, repoURL, err
	}

	for _, repo := range response {
		repos = append(repos, repo.Name)
		repoURL[repo.Name] = repo.CloneURL
	}

	return repos, repoURL, nil
}
