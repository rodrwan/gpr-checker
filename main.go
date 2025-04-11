package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/sashabaranov/go-openai"
)

func mustGet(name, fallback string) string {
	if fallback != "" {
		return fallback
	}
	env := os.Getenv(name)
	if env == "" {
		panic(fmt.Sprintf("Missing required value: %s", name))
	}
	return env
}

func downloadDiff(repo string, prNumber int, githubToken string) string {
	url := fmt.Sprintf("https://api.github.com/repos/%s/pulls/%d", repo, prNumber)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic("failed to create request: " + err.Error())
	}
	req.Header.Set("Authorization", "Bearer "+githubToken)
	req.Header.Set("Accept", "application/vnd.github.v3.diff")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic("failed to download diff: " + err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		panic(fmt.Sprintf("GitHub diff download error (%d): %s", resp.StatusCode, string(body)))
	}

	diff, err := io.ReadAll(resp.Body)
	if err != nil {
		panic("failed to read diff: " + err.Error())
	}
	return string(diff)
}

func reviewWithOpenAI(diff, apiKey string) string {
	client := openai.NewClient(apiKey)
	resp, err := client.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest{
		Model: openai.GPT4,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: "You are a senior software engineer reviewing a GitHub Pull Request. Provide constructive feedback in plain English.",
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: fmt.Sprintf("Please review the following GitHub diff and provide comments:\n\n%s", diff),
			},
		},
		Temperature: 0.3,
	})
	if err != nil {
		panic("OpenAI error: " + err.Error())
	}
	return resp.Choices[0].Message.Content
}

func postGitHubComment(token, repo string, prNumber int, body string) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/issues/%d/comments", repo, prNumber)
	payload := map[string]string{"body": body}
	jsonPayload, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		panic(fmt.Sprintf("failed to create GitHub request: %v", err))
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(fmt.Sprintf("GitHub request failed: %v", err))
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		respBody, _ := io.ReadAll(resp.Body)
		panic(fmt.Sprintf("GitHub API error: %s", string(respBody)))
	}

	fmt.Println("✅ Comment posted successfully!")
}

func validateGitHubToken(token, repo string) {
	url := fmt.Sprintf("https://api.github.com/repos/%s", repo)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(fmt.Sprintf("Failed to create token validation request: %v", err))
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(fmt.Sprintf("GitHub token validation failed: %v", err))
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		panic("❌ Repository not found — check the --repo argument and GitHub token access.")
	}

	if resp.StatusCode == 401 || resp.StatusCode == 403 {
		panic("❌ Invalid or unauthorized GitHub token. Make sure it has 'repo' scope and access to this repository.")
	}

	fmt.Println("✅ GitHub token validated successfully.")
}

func main() {
	repo := flag.String("repo", "", "GitHub repository (owner/repo)")
	prStr := flag.String("pr", "", "Pull request number")
	openaiKey := flag.String("openai-key", "", "OpenAI API key")
	githubToken := flag.String("github-token", "", "GitHub token")
	flag.Parse()

	validateGitHubToken(*githubToken, *repo)

	repoVal := mustGet("GITHUB_REPOSITORY", *repo)
	prVal := mustGet("PR_NUMBER", *prStr)
	openaiVal := mustGet("OPENAI_API_KEY", *openaiKey)
	githubVal := mustGet("GITHUB_TOKEN", *githubToken)

	prNum, err := strconv.Atoi(prVal)
	if err != nil {
		panic("invalid PR number: " + err.Error())
	}

	diff := downloadDiff(repoVal, prNum, githubVal)
	review := reviewWithOpenAI(diff, openaiVal)
	postGitHubComment(githubVal, repoVal, prNum, review)
}
