package grpc

import (
	"testing"

	"github.com/stretchr/testify/suite"
	etcdv3 "go.etcd.io/etcd/client/v3"
)

type KratosTestSutie struct {
	suite.Suite
	client *etcdv3.Client
}

func TestKratosSuite(t *testing.T) {
	suite.Run(t, new(KratosTestSutie))
}
func (s *KratosTestSutie) TestKratosServer() {
	//grpcSvc := grpc.NewServer()

}
func (s *KratosTestSutie) TestKratosClient() {

}
