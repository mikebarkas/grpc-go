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
		Title:    "Blog Five",
		Content:  "Content of my fith blog post",
	}

	// Create a blog
	res, err := c.CreateBlog(context.Background(), &blogpb.CreateBlogRequest{Blog: blog})
	if err != nil {
		log.Fatalf("error")
	}
	fmt.Printf("Blog request sent %v \n", res)

	// Read blog
	blogID := res.GetBlog().GetId()
	readReq := &blogpb.ReadBlogRequest{BlogId: blogID}
	readRes, readErr := c.ReadBlog(context.Background(), readReq)
	if readErr != nil {
		log.Fatalf("error")
	}
	fmt.Printf("Blog read response: %v \n", readRes)

	// Update blog
	updatedBlog := &blogpb.Blog{
		AuthorId: "mike",
		Title:    "Blog Five Updated",
		Content:  "This content was updated",
		Id:       blogID,
	}
	upRes, upErr := c.UpdateBlog(context.Background(), &blogpb.UpdateBlogRequest{Blog: updatedBlog})
	if upErr != nil {
		log.Fatalf("error")
	}
	fmt.Printf("Blog update response: %v \n", upRes)

}
