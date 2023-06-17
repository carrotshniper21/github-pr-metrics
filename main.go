package main

import (
	"bytes"
	"encoding/json"
	"fmt"
  	"strings"
	"io/ioutil"
	"net/http"
	"time"
)

type Config struct {
  Repos []string `json:"repos"`
  Token string `json:"token"`
}

type queryData struct {
	Query string `json:"query"`
}

type responseData struct {
	Data struct {
		Repository struct {
			PullRequests struct {
				Edges []struct {
					Node struct {
						Author struct {
							Login string `json:"login"`
						} `json:"author"`
						CreatedAt time.Time `json:"createdAt"`
						MergedAt  *time.Time `json:"mergedAt"`
					} `json:"node"`
				} `json:"edges"`
			} `json:"pullRequests"`
		} `json:"repository"`
	} `json:"data"`
}

func loadConfig() ([]string, string) {
  file, err := ioutil.ReadFile("config.json")
  if err != nil {
    fmt.Println("Error reading config file:", err)
  }

  var config Config

  if err = json.Unmarshal(file, &config); err != nil {
    fmt.Println("Error parsing config file:", err)
  }

  return config.Repos, config.Token
}

func main() {
  repos, token := loadConfig()

	client := &http.Client{}

	for _, repo := range repos {
		owner, name := parseRepo(repo)
		query := buildGraphQLQuery(owner, name)
		body := queryData{Query: query}

		reqBody, err := json.Marshal(body)
		if err != nil {
			fmt.Println("Error marshaling body:", err)
			return
		}

		req, err := http.NewRequest("POST", "https://api.github.com/graphql", bytes.NewBuffer(reqBody))
		if err != nil {
			fmt.Println("Error creating request:", err)
			return
		}

		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error making request:", err)
			return
		}

		defer resp.Body.Close()

		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading response:", err)
			return
		}

		var data responseData
		if err = json.Unmarshal(respBody, &data); err != nil {
			fmt.Println("Error unmarshaling response:", err)
		}

    fmt.Println("Repo:", name)
		for _, edge := range data.Data.Repository.PullRequests.Edges {
			pr := edge.Node
			fmt.Printf("Author: %s, CreatedAt: %s, MergedAt: %v\n", pr.Author.Login, pr.CreatedAt, pr.MergedAt)
		}
	}
}

func parseRepo(repo string) (string, string) {
	parts := strings.Split(repo, "/")
	return parts[0], parts[1]
}

func buildGraphQLQuery(owner, name string) string {
    query := `query {
      repository(owner: "%s", name: "%s") {
				pullRequests(first: 10, orderBy: {field: CREATED_AT, direction: ASC}) {
					edges {
						node {
							author {
								login
							}
							createdAt
							mergedAt
						}
					}
				}
			}
    }`
    return fmt.Sprintf(query, owner, name)
}
