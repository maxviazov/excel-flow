package pipelines

import "context"

type Issue struct {
	Level  string // info|warn|error
	Code   string
	RowNum int
	Detail string
}

// Row — «типизированная» карта: после normalize тут будут string/float64/Time.
type Row = map[string]any

// Frame — пакет данных, который проходит через все стейджи.
type Frame struct {
	Raw    []map[string]string // как прочитали из Excel (строки—строки)
	Rows   []Row               // после select+normalize+filter
	Issues []Issue
	Meta   map[string]any
}

type Stage interface {
	Name() string
	Run(ctx context.Context, in *Frame) (*Frame, error)
}
