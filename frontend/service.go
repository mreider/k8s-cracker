package main

import (
	"context"

	"performance-metrics/rpc"
)

type service struct {
	cracks chan<- CrackerID
}

func NewService(cracks chan<- CrackerID) rpc.FrontendServiceServer {
	return &service{
		cracks: cracks,
	}
}

func (s *service) NotifyCracked(_ context.Context, r *rpc.NotifyCrackedRequest) (*rpc.NotifyCrackedResponse, error) {
	s.cracks <- CrackerID(r.CrackerId)
	return &rpc.NotifyCrackedResponse{}, nil
}
