package core

type Generator interface {
	Executor
	Cleanup() error
}

type Executor interface {
	Run() error
}
