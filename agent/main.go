// package main
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"strings"
	"sync"
	"time"

	pb "github.com/hopkiw/fdna/fdna"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	port      = flag.Int("port", 50051, "The server port")
	recordMap map[string]*Record
	self      = "127.0.0.1:50051"
)

const (
	unhealthyThreshold = 30 * time.Second
	deadThreshold      = 60 * time.Second
)

type server struct {
	pb.UnimplementedFdnaServer
}

// Record is a record
type Record struct {
	*pb.Record
	LastUpdated time.Time
}

// UpdateStatus updates a records status
func (d *Record) UpdateStatus(t time.Time) {
	if d.GetEndpoint() == self {
		return
	}
	if t.Sub(d.LastUpdated) > unhealthyThreshold {
		log.Printf("marking %v unhealthy, it was last heard from at %v", d.Endpoint, d.LastUpdated)
		d.State = pb.State_STATE_UNHEALTHY
	}
}

// Dead returns whether a record is dead
func (d *Record) Dead() bool {
	return d.State == pb.State_STATE_UNHEALTHY && time.Now().Sub(d.LastUpdated) > deadThreshold
}

// Gossip trades records
func (s *server) Gossip(ctx context.Context, in *pb.GossipRequest) (*pb.GossipResponse, error) {
	log.Printf("Gossip()")
	// lock recordmap
	now := time.Now()
	var res []*pb.Record
	for _, record := range recordMap {
		record.UpdateStatus(now)
		if record.Dead() {
			delete(recordMap, record.GetEndpoint())
		} else {
			res = append(res, record.Record)
		}
	}
	updateRecordMap(in.GetRecords())
	return &pb.GossipResponse{Records: res}, nil
}

func updateRecordMap(records []*pb.Record) {
	now := time.Now()
	for _, record := range records {
		// skip records for my host
		if record.GetEndpoint() == self {
			continue
		}
		if existing, ok := recordMap[record.GetEndpoint()]; ok {
			if record.State == pb.State_STATE_HEALTHY {
				log.Printf("i heard %v is healthy", existing.Endpoint)
				existing.State = record.State
				existing.LastUpdated = now
			}
		} else {
			log.Printf("i learned %v is healthy", record.Endpoint)
			recordMap[record.GetEndpoint()] = &Record{record, now}
		}
	}
}

// Get gets records
func (s *server) Get(ctx context.Context, in *pb.GetRequest) (*pb.GetResponse, error) {
	var records []*pb.Record
	for _, record := range recordMap {
		records = append(records, record.Record)
	}
	return &pb.GetResponse{Records: records}, nil
}

// TODO: validation
func (s *server) Heartbeat(ctx context.Context, in *pb.HeartbeatRequest) (*pb.HeartbeatResponse, error) {
	return Heartbeat(in)
}

// Heartbeat receives a heartbeat for a service
func Heartbeat(in *pb.HeartbeatRequest) (*pb.HeartbeatResponse, error) {
	log.Printf("Heartbeat()")
	now := time.Now()
	if record, ok := recordMap[in.GetRecord().GetEndpoint()]; ok {
		record.State = pb.State_STATE_HEALTHY
		record.LastUpdated = now
	} else {
		recordMap[in.GetRecord().GetEndpoint()] = &Record{in.GetRecord(), now}
	}

	return &pb.HeartbeatResponse{Result: "Success"}, nil
}

// TODO: add heartbeat goroutine for service
func main() {
	flag.Parse()

	/*
		We start our record list with our own service.
		We start a goroutine with our agent RPCs - what if we get gossip now?
		We get a bootstrap node or list of nodes.
		We gossip with a bootstrap node to get our view of the network.
		They now have knowledge of my service, and might gossip with me.
		We start a goroutine for the local client RPCs.
		We start a goroutine that will heartbeat the service.
		For simplicity we could start with one RPC service.
	*/
	recordMap = make(map[string]*Record)

	ifaces, err := net.Interfaces()
	if err != nil {
		log.Fatalf("error getting interfaces: %v", err)
	}
	for _, iface := range ifaces {
		if iface.Name == "lo" {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			log.Fatalf("error getting interfaces: %v", err)
		}
		addr := strings.Split(addrs[0].String(), "/")[0]
		self = fmt.Sprintf("%s:%d", addr, *port)
	}

	wg := new(sync.WaitGroup)
	wg.Add(3)

	// start grpc server
	go func() {
		defer wg.Done()

		s := grpc.NewServer()
		pb.RegisterFdnaServer(s, &server{})

		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		log.Printf("server listening at %v", lis.Addr())
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// start heartbeat routine
	go func() {
		defer wg.Done()

		for {
			Heartbeat(&pb.HeartbeatRequest{Record: &pb.Record{
				Svc:      "fdna",
				Endpoint: self,
				Zone:     "us-west",
				State:    pb.State_STATE_HEALTHY,
			}})
			time.Sleep(5 * time.Second)
		}
	}()

	// update records and gossip
	go func() {
		defer wg.Done()

		for {
			log.Printf("updating records")
			t := time.Now()
			for _, record := range recordMap {
				record.UpdateStatus(t)
				if record.Dead() {
					delete(recordMap, record.GetEndpoint())
				}
			}
			gossip()
			time.Sleep(5 * time.Second)
		}
	}()

	wg.Wait()
}

func gossip() {
	var endpoints []string
	for endpoint, record := range recordMap {
		if endpoint == self {
			continue
		}
		if record.Svc == "fdna" &&
			record.State == pb.State_STATE_HEALTHY {
			endpoints = append(endpoints, endpoint)
		}
	}
	if len(endpoints) == 0 {
		log.Printf("no known fdna nodes, can't gossip")
		return
	}

	peer := endpoints[rand.Intn(len(endpoints))]
	log.Printf("gossipping with %v", peer)
	conn, err := grpc.Dial(peer, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewFdnaClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var records []*pb.Record
	for _, record := range recordMap {
		records = append(records, record.Record)
	}
	r, err := c.Gossip(ctx, &pb.GossipRequest{Records: records})
	if err != nil {
		log.Printf("could not gossip with %v: %v", peer, err)
		return
	}

	updateRecordMap(r.GetRecords())

	// TODO: errors, retries, logging, context, donechannel
}
