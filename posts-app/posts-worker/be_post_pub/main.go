package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"
	"encoding/json"  // Added to serialize the post into JSON

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/peer"
    "github.com/libp2p/go-libp2p/core/network"
	"github.com/multiformats/go-multiaddr"
)

const graphqlEndpoint = "http://localhost:4002/graphql"

// randomString generates a random string of length n
func randomString(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// randomUserID generates a random 3-digit user ID
func randomUserID() string {
	digits := []rune("0123456789")
	b := make([]rune, 3) // 3-digit numeric ID
	for i := range b {
		b[i] = digits[rand.Intn(len(digits))]
	}
	return string(b)
}

// randomTitle generates a random title
func randomTitle() string {
	return randomString(10) // Random title with 10 characters
}

// Post struct to represent a post
type Post struct {
	UserID  string `json:"userID"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

// publishPost serializes a post and sends it over the libp2p stream
func publishPost(stream network.Stream) {
	// Simulating a post with random data
	userID := randomUserID()
	title := randomTitle()
	content := randomString(20)

	// Creating a post
	post := Post{
		UserID:  userID,
		Title:   title,
		Content: content,
	}

	// Serializing the post to JSON
	postJSON, err := json.Marshal(post)
	if err != nil {
		fmt.Println("Error marshaling post:", err)
		return
	}

	// Send the post over the stream
	_, err = stream.Write(postJSON)
	if err != nil {
		fmt.Println("Error sending post over stream:", err)
	}
	fmt.Printf("Published post from User %s - Title: %s - Content: %s\n", userID, title, content)
}

func main() {
	rand.Seed(time.Now().UnixNano())

	// Create a new libp2p host
	ctx := context.Background()

	h, err := libp2p.New()
	if err != nil {
		log.Fatal(err)
	}
	defer h.Close()

	// Print host information
	fmt.Println("Host ID:", h.ID())
	fmt.Println("Listening addresses:")
	for _, addr := range h.Addrs() {
		fmt.Println(addr)
	}

	// Example of connecting to a peer
	peerID, err := peer.Decode("12D3KooWDjqL3pKoEag7pZZHkbfr5Ue3Xfb3B368ZiSZwzRLDRdk") // your peer ID
	if err != nil {
		log.Fatal(err)
	}

	// Construct the multiaddr to connect
	maddr, err := multiaddr.NewMultiaddr("/ip4/127.0.0.1/tcp/60841/p2p/" + peerID.String())
	if err != nil {
		log.Fatal(err)
	}

	// Get the peer information from the multiaddr
	info, err := peer.AddrInfoFromP2pAddr(maddr)
	if err != nil {
		log.Fatal(err)
	}

	// Connect to the peer
	if err := h.Connect(ctx, *info); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to:", peerID)

	// Set up a ticker to publish posts every second
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	// Loop to publish posts every second
	for {
		select {
		case <-ticker.C:
			// Open a new stream to send posts to the peer
			stream, err := h.NewStream(ctx, peerID, "/my-protocol/1.0.0")
			if err != nil {
				log.Fatal(err)
			}

			// Publish post over the stream
			publishPost(stream)
		}
	}
}



// package main

// import (
// 	"context"
// 	"fmt"
// 	"math/rand"
// 	"time"

// 	"github.com/machinebox/graphql"
// )

// const graphqlEndpoint = "http://localhost:4002/graphql"

// func randomString(n int) string {
// 	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
// 	b := make([]rune, n)
// 	for i := range b {
// 		b[i] = letters[rand.Intn(len(letters))]
// 	}
// 	return string(b)
// }

// func randomUserID() string {
// 	digits := []rune("0123456789")
// 	b := make([]rune, 3) // 8-digit numeric ID
// 	for i := range b {
// 		b[i] = digits[rand.Intn(len(digits))]
// 	}
// 	return string(b)
// }

// func main() {
// 	rand.Seed(time.Now().UnixNano())
// 	client := graphql.NewClient(graphqlEndpoint)

// 	for {
// 		// Add timestamps to ensure uniqueness
// 		title := fmt.Sprintf("%s-%d", randomString(5), time.Now().UnixNano())
// 		content := fmt.Sprintf("%s-%d", randomString(20), time.Now().UnixNano())
// 		userId := randomUserID()

// 		req := graphql.NewRequest(`
// 			mutation($title: String!, $content: String!, $userId: String!) {
// 				createPost(title: $title, content: $content, userId: $userId) {
// 					id
// 				}
// 			}
// 		`)
// 		req.Var("title", title)
// 		req.Var("content", content)
// 		req.Var("userId", userId)

// 		ctxTimeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		
// 		var respData map[string]interface{}
// 		err := client.Run(ctxTimeout, req, &respData)
// 		cancel()

// 		if err != nil {
// 			fmt.Println("❌ Error creating post:", err)
// 		} else {
// 			fmt.Println("📝 Created synthetic post:", title, "by user", userId)
// 		}

// 		time.Sleep(1 * time.Second)
// 	}
// }