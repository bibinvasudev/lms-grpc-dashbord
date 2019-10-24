// websockets.go
package main

import (
    "fmt"
    "net/http"
    "log"
    "os"
    "time"

    "github.com/gorilla/websocket"

    "golang.org/x/net/context"
    "google.golang.org/grpc"
    pb "github.com/bibinvasudev/grpc-recipe/grpc-logmonitoring/logmonitoring"

)

const (
	address     = "localhost:50051"
	defaultQuery = "world"
)


var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
}


func communicate_with_grpc_server() string {

	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewLoggerClient(conn)

	// Contact the server and print out its response.
	query := defaultQuery
	if len(os.Args) > 1 {
		query = os.Args[1]
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.DashBoardLogManagement(ctx, &pb.LogRequest{Query: query})
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	log.Printf("Logs: %s", r.GetMessage())
        return string(r.GetMessage())
}


func main() {
    // Create a simple file server
    fs := http.FileServer(http.Dir("../public"))
    http.Handle("/", fs)
    http.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
        conn, _ := upgrader.Upgrade(w, r, nil) // error ignored for sake of simplicity

        for {
            // Read message from browser
            msgType, msg, err := conn.ReadMessage()
            if err != nil {
                return
            }

            // Print the message to the console
            fmt.Printf("%s sent: %s\n", conn.RemoteAddr(), string(msg))
            grpcmessage := communicate_with_grpc_server()
            msg = []byte(grpcmessage)
            fmt.Printf(grpcmessage)
            // Write message back to browser
            if err = conn.WriteMessage(msgType, msg); err != nil {
                return
            }
        }
    })

    http.ListenAndServe(":8080", nil)
}
