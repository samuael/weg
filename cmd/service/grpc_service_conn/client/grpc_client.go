package client

import (
	"bufio"
	"context"
	"encoding/json"
	"log"
	"net"
	"strings"
	"time"

	pb "github.com/samuael/Project/Weg/cmd/service/grpc_service_conn/proto"
	"github.com/samuael/Project/Weg/internal/pkg/entity"
	"google.golang.org/grpc"
)

type GrpcClient struct {
	Servers []string
}

func NewGrpcClient() *GrpcClient {
	return &GrpcClient{
		Servers: []string{},
	}
}

// I Will Create a socket Connection to the IP Broadcasting unit to get a list of hosts
func (client *GrpcClient) UpdatePeerServers() {
	connected := false
	conn, er := net.Dial("tcp", "localhost:3000")
	if er != nil {
		println("ERROR : WHILE Connecting the Server ")
	} else {
		connected = true
		// defer conn.Close()
	}
	for {
		if !connected {
			time.Sleep(time.Duration(1000))
			// println("Trying to Reconnect ... ")
			conn, er = net.Dial("tcp", "localhost:3000")
			if er != nil {
				// println("ERROR : WHILE Connecting the Server ")
				connected = false
				continue
			}
		}
		defer conn.Close()
		// println("Writing the PORT to the server")
		writer := bufio.NewWriter(conn)

		_, ers := writer.Write(append([]byte(entity.GRPC_PORT), '\n'))
		writer.Flush()
		if ers != nil {
			connected = false
			continue
		}
		println("CONNECTED ... ")
		connected = true
		reader := bufio.NewReader(conn)
		for {
			message, er := reader.ReadBytes('\n')

			if er != nil {
				println("Error whle reading the message ....")
				connected = false
				break
			}
			val := strings.Trim(string(string(message)), "\n ")
			message = []byte(val)
			println("Executing  ... ")
			er = json.Unmarshal(message, &client.Servers)
			if er != nil {
				println("ERROR WHILE PARSING THE MESSAGE ... ", er.Error())
				continue
			}
			client.RemoveMyHostFromGRPCHostsList()
			for _, ee := range client.Servers {
				println(ee)
			}
		}
	}
}

func (client *GrpcClient) RemoveMyHostFromGRPCHostsList() {
	ind := -1
	for in, val := range client.Servers {
		if strings.Split(val, ":")[1] == entity.GRPC_PORT {
			ind = in
		}
	}
	if ind == 0 {
		client.Servers = client.Servers[1:]
	}
	if ind != -1 {
		client.Servers = append(client.Servers[0:ind], client.Servers[ind+1:]...)
	}
}

func (client *GrpcClient) BroadcastEEMessage(eebin entity.EEMBinary) bool {
	// Set up a connection to the server.
	success := false
	for _, host := range client.Servers {
		conn, err := grpc.Dial(host, grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		defer conn.Close()
		c := pb.NewMessageServiceClient(conn)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		r, err := c.HandleEEMessage(ctx, &pb.EEBinary{UserID: eebin.UserID, Data: string(eebin.Data)})
		if err != nil {
			// log.Fatalf("could not greet: %v", err)
			continue
		}
		if err != nil || !(r.Success) {
			return success
		}
		success = true
	}
	return success
}
