package main

import (
	"context"
	"io/ioutil"
	"log"
	"net"

	"github.com/mjudeikis/grpc-video-streaming/proto"
	"google.golang.org/grpc"
)

const (
	port = ":4040"
)

type server struct {
}

func main() {
	log.Println("starting server")
	conn, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	proto.RegisterStreamServiceServer(s, &server{})

	log.Println("serve :4040")
	if err := s.Serve(conn); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (s *server) GetFile(ctx context.Context, request *proto.Request) (*proto.Response, error) {
	fileName := request.GetFilename()
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal("File Reading error: ", err.Error())
	}

	return &proto.Response{Content: data}, nil
}
