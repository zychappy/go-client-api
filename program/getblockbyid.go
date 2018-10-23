package main

import (
	"flag"
	"fmt"
	"github.com/zychappy/go-client-api/service"
	"log"
	"strings"
)

func main() {
	grpcAddress := flag.String("grpcAddress", "54.236.37.243:50051",
		"gRPC address: <IP:port> example: -grpcAddress localhost:50051")

	hash := flag.String("hash",
		"2b6f3bd46c01072d2888813980def8551a64a0dbe9c02f9dd13ae4c363f436dc",
		"hash: <block hash>")

	flag.Parse()

	if (strings.EqualFold("", *hash) && len(*hash) == 0) || (strings.EqualFold("", *grpcAddress) && len(*grpcAddress) == 0) {
		log.Fatalln("./get-block-by-id -grpcAddress localhost" +
			":50051 -hash <block hash>")
	}

	client := service.NewGrpcClient(*grpcAddress)
	client.Start()
	defer client.Conn.Close()

	block := client.GetBlockById(*hash)

	fmt.Printf("block: %v\n", block)
}
