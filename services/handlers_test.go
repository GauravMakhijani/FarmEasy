package services

import (
	"FarmEasy/mocks"
	"net/http"
	"reflect"
	"testing"

	"github.com/stretchr/testify/suite"
)

type HandlerTestSuite struct {
	suite.Suite
	service *mocks.Service
}

func TestHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(HandlerTestSuite))
}

func (suite *HandlerTestSuite) SetupTest() {
	suite.service = &mocks.Service{}
}

func (suite *HandlerTestSuite) TearDownTest() {
	suite.service.AssertExpectations(suite.T())
}

func (s *HandlerTestSuite) Test_registerHandler() {
	t := s.T()
	type args struct {
		deps dependencies
	}
	tests := []struct {
		name string
		args args
		want http.HandlerFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := registerHandler(tt.args.deps); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("registerHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}
