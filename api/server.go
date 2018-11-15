package api

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"
	"unicode"

	"github.com/xdefrag/hash-ketchum/api/pb"
	"github.com/xdefrag/hash-ketchum/pkg/types"
	"github.com/xdefrag/hash-ketchum/pkg/usecase"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

// Server struct with dependencies.
type Server struct {
	port   int
	auth   Authorizer
	hash   HashStorer
	logger *log.Logger
}

// Authorizer for authorizing login from request metadata.
type Authorizer interface {
	Authorize(login string) bool
}

// HashStorer for hash store usecase.
type HashStorer interface {
	Store(context.Context, types.Hash) error
}

// NewServer constructor for Server struct.
func NewServer(port int, auth Authorizer, hash HashStorer, logger *log.Logger) Server {
	return Server{port, auth, hash, logger}
}

// Run server, blocking operation.
func (s Server) Run() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		return err
	}

	gs := grpc.NewServer(
		grpc.UnaryInterceptor(s.uiAuth),
	)

	pb.RegisterHashKetchumServer(gs, s)
	reflection.Register(gs)

	s.log(fmt.Sprintf("Starting server on :%d", s.port))

	return gs.Serve(lis)
}

// Submit is HashStore handler.
func (s Server) Submit(ctx context.Context, in *pb.HashRequest) (*pb.HashResponse, error) {
	r := &pb.HashResponse{}

	err := s.hash.Store(ctx, types.Hash{
		Login:     s.login(ctx),
		Hash:      in.Hash,
		Timestamp: time.Now().Unix(),
	})

	if err == usecase.ErrNoLeadingZeroes {
		r.Result = err.Error()

		return r, nil
	}

	if err != nil {
		r.Error = err.Error()

		return r, status.Error(codes.Internal, err.Error())
	}

	r.Result = "Success"

	return r, nil
}

func (s Server) uiAuth(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if err := s.authorize(ctx); err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	return handler(ctx, req)
}

func (s Server) login(ctx context.Context) string {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if len(md["login"]) > 0 {
			return md["login"][0]
		}
	}

	return ""
}

func (s Server) authorize(ctx context.Context) error {
	login := s.login(ctx)

	if login == "" {
		return errEmptyLogin
	}

	if !s.validateLogin(login) {
		return errInvalidLogin
	}

	if !s.auth.Authorize(login) {
		return errAccessDenied
	}

	return nil

}

func (s Server) validateLogin(login string) bool {
	if len(login) > 255 {
		return false
	}

	for _, l := range login {
		if !unicode.IsLetter(l) && !unicode.IsDigit(l) {
			return false
		}
	}

	return true
}

func (s Server) log(log string) {
	if s.logger != nil {
		s.logger.Println(log)
	}
}
