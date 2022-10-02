//grpc client code

package main

import (
	"context"

	"google.golang.org/grpc"

	"log"

	pb "github.com/iamEzaz/grpc-client"
)

func main() {

	log.Println("Dialng to port 8080")
	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewTodoServiceClient(conn)

	//call CreateTodo
	log.Println("CreateTodo")
	r, err := c.CreateTodo(context.Background(), &pb.CreateTodoRequest{
		Title: "Laundry",
		Text:  "Do laundry",
	})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("\nCreated todo: %s\n %s\n%s", r.GetId(), r.GetTitle(), r.GetText())

	rr, err := c.CreateTodo(context.Background(), &pb.CreateTodoRequest{
		Title: "Study",
		Text:  "Do study",
	})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	//log response
	log.Printf("\nCreated todo: %s\n %s\n%s", rr.GetId(), rr.GetTitle(), rr.GetText())

	//use getToDo by id 1
	log.Println("\n\n FEtching all todos")
	rrr, err := c.GetAllTodos(context.Background(), &pb.GetAllTodosRequest{})
	if err != nil {
		log.Fatalf("could not found todos due to %v", err)
	}
	log.Printf("\n All todos are : %s", rrr.Todos)

	log.Println("\n Get all StreamTodos")

	//implement StreamTodos
	stream, err := c.StreamTodos(context.Background(), &pb.GetAllTodosRequest{})
	if err != nil {
		log.Fatalf("could not found todos due to %v", err)
	}

	//read from stream
	for {
		todo, err := stream.Recv()
		if err != nil {
			log.Fatalf("could not found todos due to %v", err)
		}
		log.Printf("\n Stream todo: %s\n %s\n %s", todo.GetId(), todo.GetTitle(), todo.GetText())
	}

}
