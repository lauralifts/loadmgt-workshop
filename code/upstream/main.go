package main

import (
	"context"
	"fmt"
	"log"
	"math/rand/v2"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"github.com/paulbellamy/ratecounter"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/sync/semaphore"
	"google.golang.org/grpc"

	"google.golang.org/grpc/codes"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
	"google.golang.org/grpc/status"
)

var latency_msec = 0
var parallelism = 1
var errorRate = 0.0
var confLock sync.RWMutex
var sem = semaphore.NewWeighted(int64(parallelism))
var gradient = false
var counter = ratecounter.NewRateCounter(10 * time.Second)

var (
	healthchecks = promauto.NewCounter(prometheus.CounterOpts{
		Name: "healthchecks_total",
		Help: "The total number of HTTP healthchecks completed",
	})
	http_requests_in = promauto.NewCounter(prometheus.CounterOpts{
		Name: "http_requests_in_total",
		Help: "The total number of HTTP requests received",
	})
	grpc_requests_in = promauto.NewCounter(prometheus.CounterOpts{
		Name: "grpc_requests_in_total",
		Help: "The total number of gRPC requests received",
	})
	http_requests_complete = promauto.NewCounter(prometheus.CounterOpts{
		Name: "http_requests_completed_total",
		Help: "The total number of HTTP requests served",
	})
	grpc_requests_complete = promauto.NewCounter(prometheus.CounterOpts{
		Name: "grpc_requests_completed_total",
		Help: "The total number of gRPC requests served",
	})
)

type helloServer struct {
	pb.UnimplementedGreeterServer
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	gl_str := os.Getenv("GRADIENT_LATENCY")
	if gl_str == "true" {
		gradient = true
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

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/config", configUpdate)
	http.HandleFunc("/", hello)
	http.HandleFunc("/hipri", hello)
	http.HandleFunc("/health", healthcheck)

	fmt.Printf("Listening on port %s\n", port)
	go http.ListenAndServe(":"+port, nil)

	lis, err := net.Listen("tcp", ":"+grpc_port)
	if err != nil {
		log.Fatalln("Failed to listen:", err)
	}

	// Create a gRPC server object
	s := grpc.NewServer()
	// Attach the Greeter service to the server
	pb.RegisterGreeterServer(s, NewServer())
	// Serve gRPC Server
	log.Println("Serving gRPC on port %s", grpc_port)
	log.Fatal(s.Serve(lis))
}

func NewServer() *helloServer {
	return &helloServer{}
}

func (s *helloServer) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	grpc_requests_in.Inc()

	confLock.RLock()
	semRef := sem
	wait := time.Duration(latency_msec) * time.Millisecond
	confLock.RUnlock()

	semRef.Acquire(context.TODO(), 1)
	time.Sleep(wait)
	semRef.Release(1)
	grpc_requests_complete.Inc()

	doErr := rand.Float64()
	if doErr < errorRate {
		return nil, status.Error(codes.Internal, "oops I did it again")
	} else {
		fmt.Printf("Sending greeting")
		return &pb.HelloReply{Message: in.Name + " hello"}, nil
	}
}

func healthcheck(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "ok\n")
	healthchecks.Inc()
}

func hello(w http.ResponseWriter, req *http.Request) {
	counter.Incr(1)
	http_requests_in.Inc()

	confLock.RLock()
	semRef := sem
	wait := time.Duration(latency_msec) * time.Millisecond
	if gradient {
		if wait == 0 {
			wait = 100 * time.Millisecond
		}

		// starts to increase latency after 20 qps
		grad := float64(counter.Rate()) / float64(20) // multiply by half current qps
		log.Printf("Grad is %f\n", grad)
		if grad < 1 {
			grad = 1
		} else if grad > 10 {
			grad = 10
		}
		log.Printf("Adjusted grad is %f\n", grad)

		wait *= time.Duration(grad)
		log.Printf("Gradient wait is %d\n", wait.Milliseconds())
	}

	confLock.RUnlock()

	semRef.Acquire(context.TODO(), 1)
	time.Sleep(wait)
	semRef.Release(1)

	doErr := rand.Float64()
	if doErr < errorRate {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		fmt.Fprintf(w, "hello\n")
	}
	http_requests_complete.Inc()
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

func configUpdate(w http.ResponseWriter, req *http.Request) {
	confLock.Lock()
	defer confLock.Unlock()
	newPll, ok := getVal(req, "parallelism")
	if ok {
		parallelism = newPll
		sem = semaphore.NewWeighted(int64(parallelism))
	}

	newLatency, ok := getVal(req, "latency")
	if ok {
		latency_msec = newLatency
	}

	newErrorRate, ok := getValFloat(req, "error_rate")
	if ok {
		errorRate = newErrorRate
	}

	fmt.Fprintf(w, "Parallelism %d, latency %d milliseconds, error rate %v", parallelism, latency_msec, errorRate)
}

func getVal(req *http.Request, param string) (int, bool) {
	params, _ := url.ParseQuery(req.URL.RawQuery)
	if len(params[param]) != 1 {
		return 0, false
	}

	val, err := strconv.Atoi(params[param][0])

	if err != nil {
		log.Printf("Can't parse new %s - %v", param, err)
		return 0, false
	}
	return val, true
}

func getValFloat(req *http.Request, param string) (float64, bool) {
	params, _ := url.ParseQuery(req.URL.RawQuery)
	if len(params[param]) != 1 {
		return 0, false
	}

	val, err := strconv.ParseFloat(params[param][0], 64)

	if err != nil {
		log.Printf("Can't parse new %s - %v", param, err)
		return 0, false
	}
	return val, true
}
