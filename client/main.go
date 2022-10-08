package main

import (
	"context"
	"flag"
	"log"
	"os"
	"time"

	pb "github.com/hopkiw/fdna/fdna"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	defaultName = "world"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
)

func main() {
	flag.Parse()
	if len(os.Args) == 0 {
		log.Fatalf("Specify action")
	}
	log.Printf("Action %s", os.Args[1])
	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewFdnaClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if os.Args[1] == "init" {
		var records []*pb.Record
		records = append(records, &pb.Record{
			Svc:      "my-service",
			Zone:     "us-west",
			Endpoint: "127.0.0.1:8008",
			State:    pb.State_STATE_HEALTHY})
		log.Printf("Will send records: %v", records)

		r, err := c.Gossip(ctx, &pb.GossipRequest{Records: records})
		if err != nil {
			log.Fatalf("could not gossip: %v", err)
		}
		log.Printf("Records: %v", r.GetRecords())
		os.Exit(0)
	}

	if os.Args[1] == "heartbeat" {
		r, err := c.Heartbeat(ctx, &pb.HeartbeatRequest{Record: &pb.Record{
			Svc:      "fdna",
			Zone:     "us-east",
			Endpoint: os.Args[2],
			State:    pb.State_STATE_HEALTHY}})
		if err != nil {
			log.Fatalf("Could not heartbeat: %v", err)
		}
		log.Printf("Heartbeat result: %v", r.GetResult())
		os.Exit(0)
	}

	if os.Args[1] == "get" {
		r, err := c.Get(ctx, &pb.GetRequest{})
		if err != nil {
			log.Fatalf("Could not get records: %v", err)
		}
		for _, record := range r.GetRecords() {
			log.Printf("Records: %v", record)
		}
		os.Exit(0)
	}

	log.Printf("Done.")
}
