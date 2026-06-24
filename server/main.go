package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	pb "github.com/rlevidev/studying-grpc-in-go/pb/proto"
	"google.golang.org/grpc"
)

type taskServer struct {
	pb.UnimplementedTaskManagerServer
	mu      sync.Mutex
	tasks   []*pb.Task
	counter int
}

func (t *taskServer) CreateTask(ctx context.Context, req *pb.CreateTaskRequest) (*pb.CreateTaskResponse, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.counter++
	id := fmt.Sprintf("task-%d", t.counter)

	task := &pb.Task{
		Id:          id,
		Title:       req.Title,
		Description: req.Description,
	}
	t.tasks = append(t.tasks, task)

	log.Printf("Task criada: %s", id)

	return &pb.CreateTaskResponse{
		Id:          id,
		Title:       req.Title,
		Description: req.Description,
	}, nil
}

func (t *taskServer) ListTasks(req *pb.ListTasksRequest, stream pb.TaskManager_ListTasksServer) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	for _, task := range t.tasks {
		if err := stream.Send(task); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	// porta 50051 é a porta converncional para gRPC
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Falha ao ouvir: %v", err)
	}

	// Criar novo gRPC server e registrar implementação TaskManager
	grpcServer := grpc.NewServer()
	pb.RegisterTaskManagerServer(grpcServer, &taskServer{})

	log.Println("Servidor gRPC está rodando em :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Falha no servidor: %v", err)
	}
}
