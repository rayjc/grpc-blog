package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/rayjc/grpc-blog/blogpb"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
)

var collection *mongo.Collection

type server struct{}

type blogItem struct {
	Id       primitive.ObjectID `bson:"_id,omiempty"`
	AuthorId string             `bson:"author_id"`
	Content  string             `bson:"content"`
	Title    string             `bson:"title"`
}

func main() {
	// log file name and line number upon crash
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	fmt.Println("Blog service started.")

	fmt.Println("Connecting to mongodb...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
	if err != nil {
		log.Fatalf("Failed to connect to mongodb: %v", err)
	}

	collection = client.Database("devdb").Collection("blog")

	listener, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	opts := []grpc.ServerOption{}
	s := grpc.NewServer(opts...)
	blogpb.RegisterBlogServiceServer(s, &server{})

	go func() {
		fmt.Println("Staring server...")
		if err := s.Serve(listener); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	// Wait for exit signal (ctrl+c)
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	// Block until signal is received then pull from channel
	<-ch
	fmt.Println("Stopping server...")
	s.Stop()
	fmt.Println("Closing listener...")
	listener.Close()
	fmt.Println("Closing MongoDB connection...")
	client.Disconnect(ctx)
	fmt.Println("Exited")
}
