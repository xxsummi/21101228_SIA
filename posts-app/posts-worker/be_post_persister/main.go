package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	// "github.com/machinebox/graphql"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

const (
	wsEndpoint   = "ws://localhost:4002/graphql"
	httpEndpoint = "http://localhost:4002/graphql"
)

type Post struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	UserID  string `json:"userId"`
}

type SubscriptionResponse struct {
	Data struct {
		PostAdded Post `json:"postAdded"`
	} `json:"data"`
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Hour)
	defer cancel()

	// Connect to WebSocket with subprotocol
	conn, _, err := websocket.Dial(ctx, wsEndpoint, &websocket.DialOptions{
		Subprotocols: []string{"graphql-transport-ws"},
	})
	if err != nil {
		log.Fatal("WebSocket dial error:", err)
	}
	defer conn.Close(websocket.StatusNormalClosure, "closing")

	// Send connection_init with empty payload
	initMessage := map[string]interface{}{
		"type": "connection_init",
		"payload": map[string]interface{}{},
	}
	err = wsjson.Write(ctx, conn, initMessage)
	if err != nil {
		log.Fatal("connection_init failed:", err)
	}

	// Start subscription
	subscribeMessage := map[string]interface{}{
		"id": "1",
		"type": "subscribe",
		"payload": map[string]interface{}{
			"query": `
				subscription {
					postAdded {
						id
						title
						content
						userId
					}
				}
			`,
		},
	}
	err = wsjson.Write(ctx, conn, subscribeMessage)
	if err != nil {
		log.Fatal("subscribe failed:", err)
	}

	// client := graphql.NewClient(httpEndpoint)

	fmt.Println("🚀 Listening for new posts...")

	for {
		var msg map[string]interface{}
		err = wsjson.Read(ctx, conn, &msg)
		if err != nil {
			log.Fatal("WebSocket read error:", err)
		}

		// Check if it's a "next" message (data from subscription)
		if msgType, ok := msg["type"].(string); ok && msgType == "next" {
			payload := msg["payload"].(map[string]interface{})
			dataBytes, _ := json.Marshal(payload)

			var subResp SubscriptionResponse
			err := json.Unmarshal(dataBytes, &subResp)
			if err != nil {
				log.Println("Failed to unmarshal subscription response:", err)
				continue
			}

			post := subResp.Data.PostAdded
			fmt.Printf("✨ New Post: ID=%s Title=%s Content=%s UserID=%s\n", post.ID, post.Title, post.Content, post.UserID)

			// Save post via HTTP mutation if needed
			// req := graphql.NewRequest(`
			// 	mutation ($title: String!, $content: String!, $userId: String!) {
			// 		createPost(title: $title, content: $content, userId: $userId) {
			// 			id
			// 		}
			// 	}
			// `)
			// req.Var("title", post.Title)
			// req.Var("content", post.Content)
			// req.Var("userId", post.UserID)

			// ctxHTTP := context.Background()
			// if err := client.Run(ctxHTTP, req, nil); err != nil {
			// 	log.Println("Failed to save post:", err)
			// }
		}
	}
}
