package main

import (
	"context"
	"io"
	"log"
	"time"

	pb "github.com/rlevidev/studying-grpc-in-go/pb/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Em ambiente local podemos utilizar "insecure.NewCredentials", mas em produção podemos utilizar TLS
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Falha ao se conectar: %v", err)
	}
	defer conn.Close()

	client := pb.NewTaskManagerClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tasks := []struct {
		title string
		desc string
	}{
		{"Estudar gRPC", "Aprender mais sobre gRPC"},
		{"Criar microserviços", "Usar gRPC para comunicação serviço para serviço"},
	}

	for _, t := range tasks {
		resp, err := client.CreateTask(ctx, &pb.CreateTaskRequest{
			Title: t.title,
			Description: t.desc,
		})
		if err != nil {
			log.Fatalf("Falha no CreateTask: %v", err)
		}
		log.Printf("Criado: %s - %s", resp.Id, resp.Title)
	}

	log.Println("\nListando todas as tasks (transmissão de servidor):")

	stream, err := client.ListTasks(ctx, &pb.ListTasksRequest{})
	if err != nil {
		log.Fatalf("Falha no ListTasks: %v", err)
	}

	for {
		task, err := stream.Recv()
		if err == io.EOF {
			break	// Server fecha enviando todas as tasks
		}
		if err != nil {
			log.Fatalf("Erro ao receber a task: %v", err)
		}
		log.Printf("-> [%s] %s: %s", task.Id, task.Title, task.Description)
	}

	log.Println("\nDone!")
}
