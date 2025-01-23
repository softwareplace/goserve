package service

import "github.com/softwareplace/http-utils/example/gen"

type Service interface {
	IsWorkingV2() gen.BaseResponse
	IsWorking() gen.BaseResponse
}

type _service struct {
}

func NewService() Service {
	return &_service{}
}

func (s *_service) IsWorkingV2() gen.BaseResponse {
	message := "Test v2 it's working"
	code := 200
	success := true
	timestamp := 1625867200

	return gen.BaseResponse{
		Message:   &message,
		Code:      &code,
		Success:   &success,
		Timestamp: &timestamp,
	}
}

func (s *_service) IsWorking() gen.BaseResponse {
	message := "It's working"
	code := 200
	success := true
	timestamp := 1625867200

	return gen.BaseResponse{
		Message:   &message,
		Code:      &code,
		Success:   &success,
		Timestamp: &timestamp,
	}
}
