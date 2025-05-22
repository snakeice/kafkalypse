package kafka

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/snakeice/kafkalypse/internal/config"
	"github.com/snakeice/kafkalypse/internal/tools"
)

type KafkaConnectionReq struct{}
type KafkaConnectionRes struct {
	KafkaService *Service
	Err          error
}

func Connect(ctxConfig *config.KafkaContext) tea.Cmd {
	service, err := NewKafkaClient(ctxConfig)
	return tools.WrapCmd(KafkaConnectionRes{
		KafkaService: service,
		Err:          err,
	})

}

func GetConn() tea.Cmd { return tools.WrapCmd(KafkaConnectionReq{}) }
