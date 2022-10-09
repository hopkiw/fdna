# Failure Detection Network Agent

FDNA is an agent for sharing health information of a distributed network of
service nodes using a gossip protocol. Applications regularly send health status
to their local agent, which shares the information across the network.
Applications can query the state of the network to select healthy nodes for a
given service.

## Details

The FDNA agent is written in go and uses [gRPC](https://grpc.io) to offer RPCs
for accepting local heartbeats or remote gossips. The agent periodically chooses
1 random healthy fdna node to gossip to. The target of a gossip also shares its
view of the network in response. Both nodes integrate the new records into the
known record list, updating incoming records with a newer timestamp and healthy
status, and adding unknown records. The agent periodically updates record state
where a record which hasn't been updated in 30 seconds is marked unhealthy, and
a record which hasn't been updated in 60 seconds is deleted.

## RPCs

The fdna package offers client and server code for the relevant RPCs:

### Heartbeat 

**Heartbeat** is the local mechanism for an application to advertise a healthy
service endpoint. The application regularly invokes this RPC and the FDNA
agent is responsible for sharing this across the network. If the application
stops heartbeating, the record will be marked unhealthy and eventually deleted
by all nodes.

    Heartbeat(Record) -> None

### Gossip 

**Gossip** is the bidirectional communication RPC that is used to update the
network. Gossip sends the entire list of known records from initiator to target,
and the initiator receives the records known by the target in return. Both
agents integrate updates into their view of the network.

    Gossip([]Record) -> []Record

### Get 

**Get** is like gossip with no input and is used for local apps to get a view of
the network for any purpose, such as selecting a healthy endpoint for a given
service.

    Get() -> []Record

## Building

The various components can easily be built with the go build toolchain.

    $ go run agent/main.go &
    $ go run client/main.go get
    2022/10/09 06:59:12 Action get
    2022/10/09 06:59:12 Records: svc:"fdna" zone:"us-west" endpoint:"172.17.0.2:50051" state:STATE_HEALTHY last_updated:{seconds:1665298746 nanos:647596351}             

The agent and client binaries can also easily be built into a Docker image using
the provided Dockerfile.

    $ docker build -t fnda .
    [...]
    $ container=$(docker run fdna /agent)
    $ docker exec $container /client get
    2022/10/09 06:59:12 Action get
    2022/10/09 06:59:12 Records: svc:"fdna" zone:"us-west" endpoint:"172.17.0.2:50051" state:STATE_HEALTHY last_updated:{seconds:1665298746 nanos:647596351}             
