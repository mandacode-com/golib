package server

import "google.golang.org/grpc"

type GRPCRegisterer interface {
	Register(server *grpc.Server) error
}
