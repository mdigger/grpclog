package grpclog_test

import (
	"context"
	"io"

	"github.com/mdigger/grpclog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Example() {
	conn, err := grpc.Dial("localhost:3000",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	stream, err := grpclog.Receiver(context.Background(), conn)
	if err != nil {
		panic(err)
	}

	for {
		logMsg, err := stream()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}

		grpclog.Print(logMsg)
	}
}
