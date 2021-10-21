package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/rayjc/grpc-blog/blogpb"
	"google.golang.org/grpc"
)

type server struct{}

func main() {
	// log file name and line number upon crash
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	
	fmt.Println("Blog service started.")

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
	fmt.Println("Exited")
}
