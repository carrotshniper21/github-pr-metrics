package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

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

func main() {
	repos := []string{"owner/repo1", "owner/repo2"} // Replace with the desired repos in the format "owner/repo"
	token := "your_github_token"                   // Replace with your GitHub API token

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
		err = json.Unmarshal(respBody, &data)
		if err != nil {
			fmt.Println("Error unmarshaling response:", err)
			return
		}

		fmt.Println("Repo:", repo)
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
	return fmt.Sprintf(`{
		"query": "query {
			repository(owner: \"%s\", name: \"%s\") {
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
		}"
	}`, owner, name)
}
