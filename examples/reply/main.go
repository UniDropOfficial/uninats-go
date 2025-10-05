package main

import (
	"context"
	"fmt"

	v1 "github.com/unidropofficial/uninats-go/gen/unidrop/uninats/v1"
	"github.com/unidropofficial/uninats-go/uninats"
	"google.golang.org/protobuf/proto"
)

var (
	client = uninats.New("nats://localhost:4222")
)

func init() {
	client.Connect()
}

func main() {
	defer client.Close()

	_ = client.Reply("example.say.hello",
		func() proto.Message { return &v1.SayHelloRequest{} },
		func(ctx context.Context, req proto.Message) (proto.Message, error) {
			request := req.(*v1.SayHelloRequest)
			return &v1.SayHelloResponse{Message: fmt.Sprintf("Hello Reply %s", request.GetName())}, nil
		},
	)

	select {}
}
