package data

import (
	"github.com/snakeice/kafkalypse/internal/pkg/tui/components/table"
	"github.com/snakeice/kafkalypse/internal/pkg/tui/shortcut"
)

type Action int

const (
	Insert Action = iota
	Delete
)

type SampleDatasource struct {
	data [][]string
}

func NewSampleDatasource() SampleDatasource {
	return SampleDatasource{
		data: [][]string{
			{"a", "b", "c"},
			{"d", "e", "f"},
			{"g", "h", "i"},
			{"j", "k", "l"},
			{"m", "n", "o"},
			{"p", "q", "r"},
			{"s", "t", "u"},
			{"v", "w", "x"},
			{"y", "z", "1"},
			{"2", "3", "4"},
			{"5", "6", "7"},
			{"8", "9", "0"},
		},
	}
}

func (s SampleDatasource) Len() int {
	return len(s.data)
}

func (s SampleDatasource) At(i int) []string {
	return s.data[i]
}

func (s SampleDatasource) Shortcuts() []shortcut.Action {
	return []shortcut.Action{
		{Action: Insert, Shortcuts: []string{"i"}, Description: "insert"},
		{Action: Delete, Shortcuts: []string{"d"}, Description: "delete"},
	}
}

func (s SampleDatasource) Cols() []table.ColHead {
	return []table.ColHead{
		{Name: "Name", Perc: .5},
		{Name: "Value", Perc: .2},
		{Name: "Type", Perc: .3},
	}
}
