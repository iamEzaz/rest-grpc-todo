//grpc client code

package main

import (
	"context"
	"fmt"
	"net/http"

	"google.golang.org/grpc"

	"log"

	"github.com/gorilla/mux"
	pb "github.com/iamEzaz/grpc-client"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {

	router := mux.NewRouter()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		dialOption := grpc.WithTransportCredentials(insecure.NewCredentials())
		conn, err := grpc.Dial("localhost:8080", dialOption)
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}

		c := pb.NewTodoServiceClient(conn)

		log.Println("CreateTodo")

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

	})

	fmt.Println("http running on 3005")
	if err := http.ListenAndServe(":3005", router); err != nil {
		panic(err)
	}

}
