# GitHub PR Metrics

A Golang tool to fetch and analyze pull request metrics for specified repositories using GitHub's GraphQL API.

## Technical Overview

This project is a single Golang file that interacts with the GitHub GraphQL API to fetch data for pull request metrics. It retrieves information for a list of specified repositories and shows the following data:

1. For each contributor (provided by GitHub username), the date of their first, second, up to the 10th pull request.
2. For the pull requests mentioned above, the time taken to approve and merge them.

### Usage

Replace the `repos` array with the desired repositories in the format "owner/repo" and provide your GitHub API token in the `token` variable. The code will fetch the first 10 pull requests for each repository and print the author's username, the created date, and the merged date (if available) for each pull request.

### Dependencies

- Go 1.16 or higher
- GitHub API token with appropriate permissions

