//server code for todo app 

package main 

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	pb "github.com/iamEzaz/grpc-client"
	"log"
	"net"
	"math/rand"
)

type Server interface {
	CreateTodo(context.Context, *pb.CreateTodoRequest) (*pb.Todo, error)
	GetAllTodos(context.Context, *pb.GetAllTodosRequest) (*pb.GetAllTodosResponse, error)
	Run()
	StreamTodos(*pb.GetAllTodosRequest, pb.TodoService_StreamTodosServer) error
}


type ToDoServer struct {
	pb.UnimplementedTodoServiceServer  //for forward compatibility
	todo_list *pb.GetAllTodosResponse
}

//implement streamTodos
func (s *ToDoServer) StreamTodos(in *pb.GetAllTodosRequest, req pb.TodoService_StreamTodosServer) error {
	log.Println("StreamTodos")
	//stream todolist to client
	for _, todo := range s.todo_list.Todos {
		if err := req.Send(todo); err != nil {
			return err
		}
	}
	return nil
}


func NewToDoServer() *ToDoServer {
	//log new server init 
	log.Println("NewToDoServer init")
	return &ToDoServer{
		todo_list: &pb.GetAllTodosResponse{},
	}
}

func (s *ToDoServer) Run() error {
	//listen to por 8080
	fmt.Println("listening on port 8080")
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	//init server
	server := grpc.NewServer()
	//register service
	pb.RegisterTodoServiceServer(server, s)
	log.Println("server started")

	//start server
	if err := server.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	return nil
}

func (s *ToDoServer) GetAllTodos(ctx context.Context, req *pb.GetAllTodosRequest) (*pb.GetAllTodosResponse, error) {
	log.Println("GetAllTodos")
	//check for null todo_list
	if s.todo_list == nil {
		//return error
		return nil, fmt.Errorf("todo_list is null")
	}
	return s.todo_list, nil

}


// CreateTodo implements TodoService.CreateTodo
func (s *ToDoServer) CreateTodo(ctx context.Context, req *pb.CreateTodoRequest) (*pb.Todo, error) {
	log.Println("CreateTodo")

	log.Printf("\nCreateTodo: %s\n%s", req.GetTitle(), req.GetText())

	//create new todo
	todo := &pb.Todo{
		Id: rand.Int31(), //int32
		Title: req.GetTitle(),
		Text: req.GetText(),
	}
	//add todo to list
	log.Println("pushing todo to list")
	s.todo_list.Todos = append(s.todo_list.Todos, todo)
	return todo, nil
}


func main(){
	//init server
	server := NewToDoServer()
	//log todolist size 
	log.Printf("todo_list size: %d", len(server.todo_list.Todos))
	//start server
	if err := server.Run(); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	
}