package main

import (
	"flag"
	"fmt"
	"net"

	"golang.org/x/net/context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"

	pb "github.com/slowteetoe/tidechecker/server"
	"github.com/slowteetoe/tidechecker/tides"
)

var (
	tls      = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	certFile = flag.String("cert_file", "testdata/server1.pem", "The TLS cert file")
	keyFile  = flag.String("key_file", "testdata/server1.key", "The TLS key file")
	port     = flag.Int("port", 10000, "The server port to listen on")
)

type tideCheckerServer struct {
	data *tides.ObservationHolder
}

// GetPrediction returns the nearest tide prediction for the station id requested (use 9410230) on 20170328
func (s *tideCheckerServer) GetPrediction(ctx context.Context, req *pb.LocationRequest) (*pb.PredictionResponse, error) {

	resp := &pb.PredictionResponse{}

	if req.StationId == "" {
		fmt.Println("no station requested")
		return resp, nil
	}

	loc, ok := s.data.Locations[req.StationId]

	if !ok {
		fmt.Printf("unable to locate station[%s]", req.StationId)
		return resp, nil
	}

	pred := loc.FindNearestPrediction("2017/03/28")

	fmt.Printf("%v\n", pred)

	return resp, nil
}

func newServer() (*tideCheckerServer, error) {
	holder := tides.ObservationHolder{Locations: make(map[string]*tides.Location)}

	err := holder.LoadDataStore("data")
	if err != nil {
		return nil, fmt.Errorf("failed to load obs data: %v", err)
	}

	for index := 0; index < 5; index++ {
		fmt.Printf("%v\n", holder.Locations["9410230"].Items[index])
	}
	return &tideCheckerServer{data: &holder}, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		grpclog.Fatalf("failed to listen on port: %v", err)
	}
	var opts []grpc.ServerOption
	if *tls {
		creds, err := credentials.NewServerTLSFromFile(*certFile, *keyFile)
		if err != nil {
			grpclog.Fatalf("failed to use tls creds: %v", err)
		}
		opts = []grpc.ServerOption{grpc.Creds(creds)}
	}
	grpcServer := grpc.NewServer(opts...)
	myServer, err := newServer()
	if err != nil {
		grpclog.Fatalf("failed to create prediction server: %v", err)
	}
	pb.RegisterTideCheckerServer(grpcServer, myServer)
	grpcServer.Serve(lis)
}
