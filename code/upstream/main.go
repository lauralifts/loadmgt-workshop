package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"golang.org/x/sync/semaphore"
	"google.golang.org/grpc"

	"loadmgt-workshop/upstream/upstream_proto"
)

var latency_msec = 0
var parallelism = 1
var sem = semaphore.NewWeighted(int64(parallelism))

type helloServer struct {
	upstream_proto.UnimplementedGreeterServer
}

func NewServer() *helloServer {
	return &helloServer{}
}

func (s *helloServer) SayHello(ctx context.Context, in *upstream_proto.HelloRequest) (*upstream_proto.HelloReply, error) {
	sem.Acquire(context.TODO(), 1)
	time.Sleep(time.Duration(latency_msec) * time.Millisecond)
	return &upstream_proto.HelloReply{Message: in.Name + "hello"}, nil
}

func hello(w http.ResponseWriter, req *http.Request) {
	sem.Acquire(context.TODO(), 1)
	time.Sleep(time.Duration(latency_msec) * time.Millisecond)
	fmt.Fprintf(w, "hello\n")
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("HTTP_PORT")
	grpc_port := os.Getenv("GRPC_PORT")
	pll_str := os.Getenv("PARALLELISM")

	pll, ok := checkEnv(pll_str)
	if ok {
		parallelism = pll
		sem = semaphore.NewWeighted(int64(parallelism))
	}

	lat_str := os.Getenv("LATENCY_MSEC")
	lat, ok := checkEnv(lat_str)
	if ok {
		latency_msec = lat
	}

	http.HandleFunc("/", hello)
	fmt.Printf("Listening on port %s\n", port)
	go log.Fatal(http.ListenAndServe(":"+port, nil))

	lis, err := net.Listen("tcp", ":"+grpc_port)
	if err != nil {
		log.Fatalln("Failed to listen:", err)
	}

	// Create a gRPC server object
	s := grpc.NewServer()
	// Attach the Greeter service to the server
	upstream_proto.RegisterGreeterServer(s, NewServer())
	// Serve gRPC Server
	log.Println("Serving gRPC on port %s", grpc_port)
	log.Fatal(s.Serve(lis))
}

func checkEnv(in string) (int, bool) {
	if len(in) == 0 {
		return 0, false
	}

	conv, ok := strconv.Atoi(in)
	if ok != nil {
		return 0, false
	}

	return conv, true
}
