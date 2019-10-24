
package main

import (
	"context"
	"log"
	"net"
        "fmt"
        "io/ioutil"
	"google.golang.org/grpc"
       pb "github.com/bibinvasudev/grpc-recipe/grpc-logmonitoring/logmonitoring"
)

const (
	port = ":50051"
)

// server is used to implement logmonitorning.LoggerServer.
type server struct {
	pb.UnimplementedLoggerServer
}

// DashBoardLogManagement implements logmonitorning.LoggerServer
func (s *server) DashBoardLogManagement(ctx context.Context, in *pb.LogRequest) (*pb.LogReply, error) {
	log.Printf("Received: %v", in.GetQuery())
        data, err := ioutil.ReadFile("access.log")
        if err != nil {
           fmt.Println("File reading error", err)
        }
        fmt.Println("Contents of file:", string(data))

	return &pb.LogReply{Message: "Log " +string(data)}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterLoggerServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
