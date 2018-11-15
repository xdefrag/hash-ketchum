package api

import (
	"context"
	"crypto/sha512"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"github.com/xdefrag/hash-ketchum/api/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// ClientConfig for client connection to grpc server.
type ClientConfig struct {
	Host string
	Port int
}

// Client struct with dependencies.
type Client struct {
	cfg    ClientConfig
	cred   Credentials
	logger *log.Logger
}

// NewClient created new client with dependencies.
func NewClient(cfg ClientConfig, cred Credentials, logger *log.Logger) Client {
	return Client{cfg, cred, logger}
}

// Run client with requests number.
func (c Client) Run(reqs int) error {
	conn, err := grpc.Dial(
		fmt.Sprintf("%s:%d", c.cfg.Host, c.cfg.Port),
		grpc.WithInsecure(),
		grpc.WithPerRPCCredentials(c.cred),
		grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor(
			grpc_retry.WithCodes(codes.Unauthenticated),
		)),
	)
	if err != nil {
		return err
	}
	defer conn.Close()

	gc := pb.NewHashKetchumClient(conn)

	var wg sync.WaitGroup
	wg.Add(reqs)

	for i := 0; i < reqs; i++ {
		go c.makeRequest(gc, &wg)
	}

	wg.Wait()

	return nil
}

func (c Client) makeRequest(gc pb.HashKetchumClient, wg *sync.WaitGroup) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	defer wg.Done()

	rh := generateRandomSha512()

	res, err := gc.Submit(ctx, &pb.HashRequest{Hash: rh}, grpc_retry.WithMax(10))
	if err != nil {
		c.log(err.Error())
		return
	}

	if res.GetResult() == "Success" {
		c.log(fmt.Sprintf("Submition succeeded: %s", rh))
	} else {
		c.log(fmt.Sprintf("Submition rejected: %s", rh))
	}

}

func (c Client) log(log string) {
	if c.logger != nil {
		c.logger.Println(log)
	}
}

func generateRandomSha512() string {
	data := make([]byte, 10)
	for i := range data {
		data[i] = byte(rand.Intn(512))
	}

	return fmt.Sprintf("%x", sha512.Sum512(data))
}
