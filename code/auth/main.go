package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"math/rand"

	v31 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	auth "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	"github.com/gogo/googleapis/google/rpc"
	"google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc"
)

// empty struct because this isn't a fancy example
type AuthorizationServer struct{}

// inject a header that can be used for future rate limiting
func (a *AuthorizationServer) Check(ctx context.Context, req *auth.CheckRequest) (*auth.CheckResponse, error) {
	num := rand.Intn(20)

	tenant := "default"
	if num < 10 {
		tenant = fmt.Sprintf("tenant-%d", num)
	}

	return &auth.CheckResponse{
		Status: &status.Status{Code: int32(rpc.OK)},
		HttpResponse: &auth.CheckResponse_OkResponse{
			OkResponse: &auth.OkHttpResponse{
				Headers: []*v31.HeaderValueOption{
					{
						Header: &v31.HeaderValue{
							Key:   "x-ext-auth-tenant-id",
							Value: tenant,
						},
					},
				},
			},
		},
	}, nil

}

func main() {
	lis, err := net.Listen("tcp", ":9010")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("listening on %s", lis.Addr())

	grpcServer := grpc.NewServer()
	authServer := &AuthorizationServer{}
	auth.RegisterAuthorizationServer(grpcServer, authServer)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
