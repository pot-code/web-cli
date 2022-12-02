package task

import (
	"bytes"
	"io"
)

type PipelineFn func(src io.Reader, dest io.Writer) error

type Pipeline struct {
	pns []PipelineFn
}

func NewPipeline() *Pipeline {
	return &Pipeline{}
}

func (p *Pipeline) AddPipelineFn(pn PipelineFn) *Pipeline {
	p.pns = append(p.pns, pn)
	return p
}

func (p *Pipeline) Apply(src io.Reader, dest io.Writer) error {
	l := len(p.pns)
	if l == 0 {
		return nil
	}

	s := src
	lastPn := p.pns[l-1]
	p.pns = p.pns[:l-1]
	for _, fn := range p.pns {
		d := new(bytes.Buffer)
		fn(s, d)
		s = d
	}
	lastPn(s, dest)
	return nil
}
