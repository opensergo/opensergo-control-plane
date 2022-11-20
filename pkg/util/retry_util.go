package util

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TimeoutCondition(err error) bool {
	s, _ := status.FromError(err)
	if s.Code() == codes.DeadlineExceeded {
		return true
	}
	return false
}
