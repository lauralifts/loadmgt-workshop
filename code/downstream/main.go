package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "google.golang.org/grpc/examples/helloworld/helloworld"
)

// config - why isn't updatedFlag getting reset
// metrics
// add metrics to upstream as well
// instructions to run it - curl config etc

type config struct {
	grpc_rate            int
	grpc_max_parallelism int
	http_rate            int
	http_max_parallelism int
	updatedFlag          bool
}

var grpc_server = ""
var http_server = ""
var conf = config{}
var confLock = sync.RWMutex{}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	http_server = os.Getenv("HTTP_SERVER")
	grpc_server = os.Getenv("GRPC_SERVER")
	config_port := os.Getenv("CONFIG_PORT")

	http.HandleFunc("/config", configUpdate)
	log.Printf("Listening on port %s\n", config_port)
	go http.ListenAndServe(":"+config_port, nil)

	doRequests()
}

func doRequests() {
	for {
		confLock.Lock()
		grpcLimiter := rate.NewLimiter(rate.Limit(conf.grpc_rate), conf.grpc_rate)
		grpcParallelism := conf.grpc_max_parallelism
		httpLimiter := rate.NewLimiter(rate.Limit(conf.http_rate), conf.http_rate)
		httpParallelism := conf.http_max_parallelism
		stop := make(chan bool, 1)
		conf.updatedFlag = false
		confLock.Unlock()

		// kick off grpcParallelism and httpParallelism threads, params are limiter and the channel
		for i := 0; i < grpcParallelism; i++ {
			go doGRPCReqsWorker(stop, grpcLimiter)
		}
		for i := 0; i < httpParallelism; i++ {
			go doHTTPReqsWorker(stop, httpLimiter)
		}

		updateNeeded := false
		for !updateNeeded {
			time.Sleep(time.Second)
			confLock.RLock()
			updateNeeded = conf.updatedFlag
			confLock.RUnlock()
		}
		stop <- true
	}
}

func doGRPCReqsWorker(stop chan bool, rl *rate.Limiter) {
	// create grpc client

	conn, err := grpc.NewClient(grpc_server, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewGreeterClient(conn)

	for len(stop) == 0 {
		rl.Wait(context.TODO())
		r, err := client.SayHello(context.TODO(), &pb.HelloRequest{Name: "load generator"})
		if err != nil {
			log.Printf("could not greet: %v", err)
		}
		log.Printf("Greeting: %s\n", r.GetMessage())
	}
}

func doHTTPReqsWorker(stop chan bool, rl *rate.Limiter) {
	for len(stop) == 0 {
		rl.Wait(context.TODO())
		res, err := http.Get(http_server)
		if err != nil {
			log.Printf("Http request to %s errored - %+v", http_server, err)
		} else {
			log.Printf("Http request to %s done, result code %d\n", http_server, res.StatusCode)
		}
	}
}

func configUpdate(w http.ResponseWriter, req *http.Request) {
	newHTTPParallelism := conf.http_max_parallelism
	newGRPCParallelism := conf.grpc_max_parallelism
	newHTTPRate := conf.http_rate
	newGRPCRate := conf.grpc_rate

	updateSeen := false

	val, ok := getVal(req, "grpc_max_parallelism")
	if ok {
		updateSeen = true
		newGRPCParallelism = val
	}

	val, ok = getVal(req, "http_max_parallelism")
	if ok {
		updateSeen = true
		newHTTPParallelism = val
	}

	val, ok = getVal(req, "grpc_rate")
	if ok {
		updateSeen = true
		newGRPCRate = val
	}

	val, ok = getVal(req, "http_rate")
	if ok {
		updateSeen = true
		newHTTPRate = val
	}

	if updateSeen {
		confLock.Lock()
		conf.updatedFlag = true
		conf.grpc_max_parallelism = newGRPCParallelism
		conf.grpc_rate = newGRPCRate
		conf.http_max_parallelism = newHTTPParallelism
		conf.http_rate = newHTTPRate
		confLock.Unlock()
	}

	confLock.RLock()
	confStr := fmt.Sprintf("%+v", conf)
	confLock.RUnlock()

	fmt.Fprintf(w, "%s", confStr)
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