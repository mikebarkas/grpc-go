package main

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"

	"github.com/mikebarkas/grpc-go/blog/blogpb"
)

func main() {

	opts := grpc.WithInsecure()

	cc, err := grpc.Dial("localhost:50051", opts)
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer cc.Close()

	c := blogpb.NewBlogServiceClient(cc)

	blog := &blogpb.Blog{
		AuthorId: "mike",
		Title:    "Third Blog",
		Content:  "Content of my third blog post",
	}
	res, err := c.CreateBlog(context.Background(), &blogpb.CreateBlogRequest{Blog: blog})
	if err != nil {
		log.Fatalf("error")
	}
	fmt.Printf("Blog request sent %v \n", res)
}
