package kafka

import (
	"context"

	"github.com/IBM/sarama"
)

func (s *Service) APIVersion(ctx context.Context) (*sarama.ApiVersionsResponse, error) {
	brokers, _, _ := s.kAdmin.DescribeCluster()
	return brokers[0].ApiVersions(&sarama.ApiVersionsRequest{
		ClientSoftwareName:    "kafkalypse",
		ClientSoftwareVersion: "0.1.0", // TODO: get version from somewhere
	})
}
