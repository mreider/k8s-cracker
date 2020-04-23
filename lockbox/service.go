package main

import (
	"context"

	"performance-metrics/rpc"
	"performance-metrics/utils"
)

const (
	PasswordSize = 5
)

type service struct{}

func NewService() rpc.LockboxServiceServer {
	return &service{}
}

func (s *service) GeneratePassword(context.Context, *rpc.GeneratePasswordRequest) (*rpc.GeneratePasswordResponse, error) {
	password, err := utils.GenerateRandomString(PasswordSize)
	if err != nil {
		return nil, err
	}
	return &rpc.GeneratePasswordResponse{
		Password: password,
	}, nil
}
