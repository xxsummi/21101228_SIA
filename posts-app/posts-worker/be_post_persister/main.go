package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	libp2p "github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/network"
)

// Post struct to represent a post
type Post struct {
	UserID  string `json:"userID"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

const graphqlEndpoint = "http://localhost:4002/graphql"

// handleStream handles incoming streams and processes the post data
func handleStream(stream network.Stream) {
	fmt.Println("Got a new stream!")

	// Read the incoming data from the stream
	buf := make([]byte, 4096)
	n, err := stream.Read(buf)
	if err != nil {
		log.Println("Error reading stream:", err)
		return
	}

	// Deserialize the post from JSON
	var post Post
	err = json.Unmarshal(buf[:n], &post)
	if err != nil {
		log.Println("Error unmarshaling post:", err)
		return
	}

	// Print the received post
	fmt.Printf("Received Post - UserID: %s, Title: %s, Content: %s\n", post.UserID, post.Title, post.Content)

	// Prepare the GraphQL mutation
	mutation := fmt.Sprintf(`
	mutation {
		createPost(title: "%s", content: "%s", userId: "%s") {
			id
			title
			content
			userId
		}
	}`, escapeString(post.Title), escapeString(post.Content), escapeString(post.UserID))

	// Create the request body
	payload := map[string]string{
		"query": mutation,
	}
	payloadBytes, _ := json.Marshal(payload)

	// Make the HTTP request to the GraphQL endpoint
	resp, err := http.Post(graphqlEndpoint, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		log.Println("Error sending GraphQL request:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println("GraphQL request failed with status:", resp.Status)
	} else {
		log.Println("Post successfully inserted via GraphQL!")
	}
}

// escapeString escapes quotes and newlines in strings for safe GraphQL use
func escapeString(input string) string {
	escaped, _ := json.Marshal(input)
	return string(escaped[1 : len(escaped)-1]) // remove leading and trailing quotes
}

func main() {
	// Create a new libp2p host
	h, err := libp2p.New()
	if err != nil {
		log.Fatal(err)
	}

	// Print host information
	fmt.Println("Host ID:", h.ID())
	fmt.Println("Listening addresses:")
	for _, addr := range h.Addrs() {
		fmt.Println(addr)
	}

	// Set the stream handler
	h.SetStreamHandler("/my-protocol/1.0.0", handleStream)

	// Block forever
	select {}
}



// package main

// import (
// 	"bytes"
// 	"encoding/json"
// 	"fmt"
// 	"log"
// 	"net/http" // To send HTTP requests
// 	"github.com/libp2p/go-libp2p"
// 	"github.com/libp2p/go-libp2p/core/network"
// )

// // Post struct to represent a post
// type Post struct {
// 	UserID  string `json:"userID"`
// 	Title   string `json:"title"`
// 	Content string `json:"content"`
// }

// // GraphQL mutation to insert post into the database
// const graphqlEndpoint = "http://localhost:4002/graphql"

// // GraphQLRequest struct represents the structure of a GraphQL request
// type GraphQLRequest struct {
// 	Query     string `json:"query"`
// 	Variables struct {
// 		UserID  string `json:"userID"`
// 		Title   string `json:"title"`
// 		Content string `json:"content"`
// 	} `json:"variables"`
// }

// // insertPost sends a GraphQL mutation to insert the post into the database
// func insertPost(post Post) error {
// 	// Construct the GraphQL mutation
// 	mutation := `
// 		mutation InsertPost($userID: String!, $title: String!, $content: String!) {
// 			insertPost(userID: $userID, title: $title, content: $content) {
// 				id
// 			}
// 		}
// 	`

// 	// Create the GraphQL request
// 	reqBody := GraphQLRequest{
// 		Query: mutation,
// 	}
// 	reqBody.Variables.UserID = post.UserID
// 	reqBody.Variables.Title = post.Title
// 	reqBody.Variables.Content = post.Content

// 	// Marshal the request body to JSON
// 	reqJSON, err := json.Marshal(reqBody)
// 	if err != nil {
// 		return fmt.Errorf("error marshalling GraphQL request: %v", err)
// 	}

// 	// Send the HTTP request to the GraphQL endpoint
// 	resp, err := http.Post(graphqlEndpoint, "application/json", bytes.NewBuffer(reqJSON))
// 	if err != nil {
// 		return fmt.Errorf("error sending request to GraphQL endpoint: %v", err)
// 	}
// 	defer resp.Body.Close()

// 	// Check the response status code
// 	if resp.StatusCode != http.StatusOK {
// 		return fmt.Errorf("received non-200 status code: %d", resp.StatusCode)
// 	}

// 	fmt.Printf("Successfully inserted post into the database. UserID: %s, Title: %s\n", post.UserID, post.Title)
// 	return nil
// }

// // handleStream handles incoming streams and processes the post data
// func handleStream(stream network.Stream) {
// 	fmt.Println("Got a new stream!")

// 	// Read the incoming data from the stream
// 	buf := make([]byte, 1024)
// 	n, err := stream.Read(buf)
// 	if err != nil {
// 		log.Fatal("Error reading stream:", err)
// 	}

// 	// Deserialize the post from JSON
// 	var post Post
// 	err = json.Unmarshal(buf[:n], &post)
// 	if err != nil {
// 		log.Fatal("Error unmarshaling post:", err)
// 	}

// 	// Print the received post (here you can insert into the database)
// 	fmt.Printf("Received Post - UserID: %s, Title: %s, Content: %s\n", post.UserID, post.Title, post.Content)

// 	// Call GraphQL mutation to insert the post into the database
// 	if err := insertPost(post); err != nil {
// 		log.Fatal("Error inserting post into database:", err)
// 	}
// }

// func main() {
// 	// Create a new libp2p host
// 	h, err := libp2p.New()
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// Print host information
// 	fmt.Println("Host ID:", h.ID())
// 	fmt.Println("Listening addresses:")
// 	for _, addr := range h.Addrs() {
// 		fmt.Println(addr)
// 	}

// 	// Set the stream handler to process incoming streams
// 	h.SetStreamHandler("/my-protocol/1.0.0", handleStream)

// 	// Block forever
// 	select {}
// }




// package main

// import (
// 	"fmt"
// 	"log"
// 	"encoding/json"  // Added to deserialize the post from JSON
// 	// "bytes"          // To handle the stream data

// 	libp2p "github.com/libp2p/go-libp2p"
// 	"github.com/libp2p/go-libp2p/core/network"
// )

// // Post struct to represent a post
// type Post struct {
// 	UserID  string `json:"userID"`
// 	Title   string `json:"title"`
// 	Content string `json:"content"`
// }

// // handleStream handles incoming streams and processes the post data
// func handleStream(stream network.Stream) {
// 	fmt.Println("Got a new stream!")

// 	// Read the incoming data from the stream
// 	buf := make([]byte, 1024)
// 	n, err := stream.Read(buf)
// 	if err != nil {
// 		log.Fatal("Error reading stream:", err)
// 	}

// 	// Deserialize the post from JSON
// 	var post Post
// 	err = json.Unmarshal(buf[:n], &post)
// 	if err != nil {
// 		log.Fatal("Error unmarshaling post:", err)
// 	}

// 	// Print the received post (here you can insert into the database)
// 	fmt.Printf("Received Post - UserID: %s, Title: %s, Content: %s\n", post.UserID, post.Title, post.Content)

// 	// Here you would call your GraphQL mutation to insert the post into the database
// 	// For now, we're just printing it
// }

// func main() {
// 	// Create a new libp2p host
// 	h, err := libp2p.New()
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// Print host information
// 	fmt.Println("Host ID:", h.ID())
// 	fmt.Println("Listening addresses:")
// 	for _, addr := range h.Addrs() {
// 		fmt.Println(addr)
// 	}

// 	// Set the stream handler to process incoming streams
// 	h.SetStreamHandler("/my-protocol/1.0.0", handleStream)

// 	// Block forever
// 	select {}
// }



// package main

// import (
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"log"
// 	"time"

// 	// "github.com/machinebox/graphql"
// 	"nhooyr.io/websocket"
// 	"nhooyr.io/websocket/wsjson"
// )

// const (
// 	wsEndpoint   = "ws://localhost:4002/graphql"
// 	httpEndpoint = "http://localhost:4002/graphql"
// )

// type Post struct {
// 	ID      string `json:"id"`
// 	Title   string `json:"title"`
// 	Content string `json:"content"`
// 	UserID  string `json:"userId"`
// }

// type SubscriptionResponse struct {
// 	Data struct {
// 		PostAdded Post `json:"postAdded"`
// 	} `json:"data"`
// }

// func main() {
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Hour)
// 	defer cancel()

// 	// Connect to WebSocket with subprotocol
// 	conn, _, err := websocket.Dial(ctx, wsEndpoint, &websocket.DialOptions{
// 		Subprotocols: []string{"graphql-transport-ws"},
// 	})
// 	if err != nil {
// 		log.Fatal("WebSocket dial error:", err)
// 	}
// 	defer conn.Close(websocket.StatusNormalClosure, "closing")

// 	// Send connection_init with empty payload
// 	initMessage := map[string]interface{}{
// 		"type": "connection_init",
// 		"payload": map[string]interface{}{},
// 	}
// 	err = wsjson.Write(ctx, conn, initMessage)
// 	if err != nil {
// 		log.Fatal("connection_init failed:", err)
// 	}

// 	// Start subscription
// 	subscribeMessage := map[string]interface{}{
// 		"id": "1",
// 		"type": "subscribe",
// 		"payload": map[string]interface{}{
// 			"query": `
// 				subscription {
// 					postAdded {
// 						id
// 						title
// 						content
// 						userId
// 					}
// 				}
// 			`,
// 		},
// 	}
// 	err = wsjson.Write(ctx, conn, subscribeMessage)
// 	if err != nil {
// 		log.Fatal("subscribe failed:", err)
// 	}

// 	// client := graphql.NewClient(httpEndpoint)

// 	fmt.Println("🚀 Listening for new posts...")

// 	for {
// 		var msg map[string]interface{}
// 		err = wsjson.Read(ctx, conn, &msg)
// 		if err != nil {
// 			log.Fatal("WebSocket read error:", err)
// 		}

// 		// Check if it's a "next" message (data from subscription)
// 		if msgType, ok := msg["type"].(string); ok && msgType == "next" {
// 			payload := msg["payload"].(map[string]interface{})
// 			dataBytes, _ := json.Marshal(payload)

// 			var subResp SubscriptionResponse
// 			err := json.Unmarshal(dataBytes, &subResp)
// 			if err != nil {
// 				log.Println("Failed to unmarshal subscription response:", err)
// 				continue
// 			}

// 			post := subResp.Data.PostAdded
// 			fmt.Printf("✨ New Post: ID=%s Title=%s Content=%s UserID=%s\n", post.ID, post.Title, post.Content, post.UserID)

// 			// Save post via HTTP mutation if needed
// 			// req := graphql.NewRequest(`
// 			// 	mutation ($title: String!, $content: String!, $userId: String!) {
// 			// 		createPost(title: $title, content: $content, userId: $userId) {
// 			// 			id
// 			// 		}
// 			// 	}
// 			// `)
// 			// req.Var("title", post.Title)
// 			// req.Var("content", post.Content)
// 			// req.Var("userId", post.UserID)

// 			// ctxHTTP := context.Background()
// 			// if err := client.Run(ctxHTTP, req, nil); err != nil {
// 			// 	log.Println("Failed to save post:", err)
// 			// }
// 		}
// 	}
// }
