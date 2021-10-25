package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

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

	// CreateBlog
	createRes, err := client.CreateBlog(context.Background(), &blogpb.CreateBlogRequest{Blog: blog})
	if err != nil {
		log.Fatalf("Failed to create blog: %v\n", err)
	}
	fmt.Printf("Blog created: %v\n\n", createRes)

	// ReadBlog
	_, readErr := client.ReadBlog(context.Background(), &blogpb.ReadBlogRequest{BlogId: "not_found"})
	if readErr != nil {
		fmt.Printf("Failed to read blog: %v\n", err)
	}
	readRes, readErr := client.ReadBlog(
		context.Background(),
		&blogpb.ReadBlogRequest{BlogId: createRes.GetBlog().GetId()},
	)
	if readErr != nil {
		log.Fatalf("Failed to read blog: %v\n", err)
	}
	fmt.Printf("Got blog: %v\n\n", readRes)

	// UpdateBlog
	newBlog := &blogpb.Blog{
		Id:       createRes.Blog.Id,
		AuthorId: "Batman",
		Title:    fmt.Sprintf("Batman's Identity (updated:%v)", time.Now().Format(time.RFC850)),
		Content:  "Bruce? Is that you?",
	}
	updateRes, updateErr := client.UpdateBlog(context.Background(), &blogpb.UpdateBlogRequest{Blog: newBlog})
	if updateErr != nil {
		log.Fatalf("Error updating blog: %v\n", updateErr)
	}
	fmt.Printf("Blog updated: %v\n\n", updateRes)

	// DeleteBlog
	deleteRes, deleteErr := client.DeleteBlog(context.Background(), &blogpb.DeleteBlogRequest{BlogId: createRes.Blog.Id})
	if deleteErr != nil {
		log.Fatalf("Error deleting blog: %v\n", deleteErr)
	}
	fmt.Printf("Blog deleted: %v\n\n", deleteRes)

	// ListBlog
	stream, err := client.ListBlog(context.Background(), &blogpb.ListBlogRequest{})
	if err != nil {
		log.Fatalf("Error calling ListBlog RPC: %v", err)
	}
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Unexpected error occured during server streaming... %v\n", err)
		}
		// parse response
		fmt.Println(res.GetBlog())
	}
}
