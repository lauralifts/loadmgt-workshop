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
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	pb "google.golang.org/grpc/examples/helloworld/helloworld"
)

type config struct {
	grpc_rate            int
	grpc_max_parallelism int
	http_rate            int
	http_max_parallelism int
	hipri                bool
	updatedFlag          bool
}

var grpc_server = ""
var http_server = ""
var conf = config{}
var confLock = sync.RWMutex{}

var (
	http_requests_made = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "http_requests_made_total",
		Help: "The total number of HTTP requests made",
	}, []string{"code", "priority"})
	grpc_requests_made = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "grpc_requests_made_total",
		Help: "The total number of gRPC requests made",
	}, []string{"code"})
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	http.DefaultTransport.(*http.Transport).MaxIdleConnsPerHost = 10000

	http_server = os.Getenv("HTTP_SERVER")
	grpc_server = os.Getenv("GRPC_SERVER")
	config_port := os.Getenv("CONFIG_PORT")

	hrstr := os.Getenv("HTTP_RATE")
	if hrstr != "" {
		res, err := strconv.Atoi(hrstr)
		if err == nil {
			conf.http_rate = res
		} else {
			log.Fatal(fmt.Sprintf("%v", err))
		}
	}

	hrpll := os.Getenv("HTTP_MAX_PARALLELISM")
	if hrpll != "" {
		res, err := strconv.Atoi(hrpll)
		if err == nil {
			conf.http_max_parallelism = res
		} else {
			log.Fatal(fmt.Sprintf("%v", err))
		}
	}

	grstr := os.Getenv("GRPC_RATE")
	if grstr != "" {
		res, err := strconv.Atoi(grstr)
		if err == nil {
			conf.grpc_rate = res
		} else {
			log.Fatal(fmt.Sprintf("%v", err))
		}
	}

	grpll := os.Getenv("GRPC_MAX_PARALLELISM")
	if grpll != "" {
		res, err := strconv.Atoi(grpll)
		if err == nil {
			conf.grpc_max_parallelism = res
		} else {
			log.Fatal(fmt.Sprintf("%v", err))
		}
	}

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/config", configUpdate)
	log.Printf("Listening on port %s\n", config_port)
	go http.ListenAndServe(":"+config_port, nil)

	log.Printf("Config is %+v", conf)

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
		req := pb.HelloRequest{Name: "load generator"}
		r, err := client.SayHello(context.TODO(), &req)

		code := codes.OK

		if err != nil {
			log.Printf("could not greet: %v", err)
			status, ok := status.FromError(err)
			if ok {
				code = status.Code()
				grpc_requests_made.With(prometheus.Labels{"code": fmt.Sprintf("%d", code)}).Inc()
			} else {
				log.Printf("Can't parse %v as grpc status", err)
				// todo inc a metric
			}
		}

		log.Printf("Greeting: %s\n", r.GetMessage())
	}
}

func doHTTPReqsWorker(stop chan bool, rl *rate.Limiter) {
	for len(stop) == 0 {
		rl.Wait(context.TODO())

		url := http_server
		confLock.RLock()
		hipri := conf.hipri
		confLock.RUnlock()

		if hipri {
			url += "/hipri"
		}

		res, err := http.Get(url)
		if err != nil {
			log.Printf("Http request to %s errored - %+v", url, err)
			// todo inc a metric
		} else {
			priority := "default"
			if hipri {
				priority = "high"
			}
			log.Printf("Http request to %s done at priority %s, result code %d\n", url, priority, res.StatusCode)
			http_requests_made.With(prometheus.Labels{"code": fmt.Sprintf("%d", res.StatusCode), "priority": priority}).Inc()
		}

		if res != nil {
			res.Body.Close()
		}
	}
}

func configUpdate(w http.ResponseWriter, req *http.Request) {
	newHTTPParallelism := conf.http_max_parallelism
	newGRPCParallelism := conf.grpc_max_parallelism
	newHTTPRate := conf.http_rate
	newGRPCRate := conf.grpc_rate
	newHipri := conf.hipri

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

	bval, ok := getValBool(req, "hipri")
	if ok {
		updateSeen = true
		newHipri = bval
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
		conf.hipri = newHipri
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

func getValBool(req *http.Request, param string) (bool, bool) {
	params, _ := url.ParseQuery(req.URL.RawQuery)
	if len(params[param]) != 1 {
		return false, false
	}

	if params[param][0] == "true" {
		return true, true
	}

	return false, true
}
