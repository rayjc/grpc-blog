package main

import (
	"context"
	"fmt"
	"log"

	"github.com/rayjc/grpc-blog/blogpb"
	"github.com/rayjc/grpc-blog/config"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Client connected.")

	opts := grpc.WithInsecure()
	url := fmt.Sprintf("localhost:%v", config.Port)
	conn, err := grpc.Dial(url, opts)
	if err != nil {
		log.Fatalf("Cannot connect: %v", err)
	}
	defer conn.Close()

	client := blogpb.NewBlogServiceClient(conn)
	// fmt.Printf("Client created: %f", client)

	blog := &blogpb.Blog{
		AuthorId: "Batman",
		Title:    "Batman's Identity",
		Content:  "Underneath the mask...",
	}
	res, err := client.CreateBlog(context.Background(), &blogpb.CreateBlogRequest{Blog: blog})
	if err != nil {
		log.Fatalf("Failed to create blog: %v", err)
	}
	fmt.Printf("Blog created: %v\n", res)
}
