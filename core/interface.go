package core

type Generator interface {
	Gen() error
	Cleanup() error
}
