package response

import (
	"encoding/json"
	"errors"
)

var (
	SUCCESS = errors.New("SUCCESS")
	FAILURE = errors.New("FAILURE")
)

type StdResp interface {
	String() string
}

type StdResponse struct {
	data interface{}
}

func NewStdResponse(data interface{}) StdResp {
	return &StdResponse{data}
}

func (s *StdResponse) String() string {
	bytes, err := json.Marshal(s.data)
	if err != nil {
		return ""
	}
	return string(bytes)
}
