package kafka

import (
	"github.com/IBM/sarama"
	tea "github.com/charmbracelet/bubbletea"
)

type KafkaConnectionMsg struct {
	Svc *Service
	Err error
}

type GetConnMsg struct{}

func GetConn() tea.Cmd { return func() tea.Msg { return GetConnMsg{} } }

type Service struct {
	kAdmin  sarama.ClusterAdmin
	kClient sarama.Client
}

func NewKafkaClient() (*Service, error) {
	kClient, err := sarama.NewClient([]string{"localhost:9093"}, nil)
	if err != nil {
		return nil, err
	}

	kAdmin, err := sarama.NewClusterAdminFromClient(kClient)
	if err != nil {
		return nil, err
	}

	return &Service{
		kAdmin:  kAdmin,
		kClient: kClient,
	}, nil
}

func (k *Service) Close() error {

	if err := k.kAdmin.Close(); err != nil {
		return err
	}

	if err := k.kClient.Close(); err != nil {
		return err
	}

	return nil
}
