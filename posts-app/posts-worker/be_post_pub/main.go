package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/machinebox/graphql"
)

const graphqlEndpoint = "http://localhost:4002/graphql"

func randomString(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func randomUserID() string {
	digits := []rune("0123456789")
	b := make([]rune, 3) // 8-digit numeric ID
	for i := range b {
		b[i] = digits[rand.Intn(len(digits))]
	}
	return string(b)
}

func main() {
	rand.Seed(time.Now().UnixNano())
	client := graphql.NewClient(graphqlEndpoint)

	for {
		// Add timestamps to ensure uniqueness
		title := fmt.Sprintf("%s-%d", randomString(5), time.Now().UnixNano())
		content := fmt.Sprintf("%s-%d", randomString(20), time.Now().UnixNano())
		userId := randomUserID()

		req := graphql.NewRequest(`
			mutation($title: String!, $content: String!, $userId: String!) {
				createPost(title: $title, content: $content, userId: $userId) {
					id
				}
			}
		`)
		req.Var("title", title)
		req.Var("content", content)
		req.Var("userId", userId)

		ctxTimeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		
		var respData map[string]interface{}
		err := client.Run(ctxTimeout, req, &respData)
		cancel()

		if err != nil {
			fmt.Println("❌ Error creating post:", err)
		} else {
			fmt.Println("📝 Created synthetic post:", title, "by user", userId)
		}

		time.Sleep(1 * time.Second)
	}
}