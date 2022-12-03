package pipeline

import (
	"bytes"
	"io"
)

type PipelineFn func(src io.Reader, dest io.Writer) error

type Pipelines struct {
	pns []PipelineFn
}

func NewPipelines(pns ...PipelineFn) *Pipelines {
	return &Pipelines{pns}
}

func (p *Pipelines) AddPipelineFn(pn PipelineFn) *Pipelines {
	p.pns = append(p.pns, pn)
	return p
}

func (p *Pipelines) Apply(src io.Reader, dest io.Writer) error {
	l := len(p.pns)
	if l == 0 {
		return nil
	}

	s := src
	lastPn := p.pop()
	for _, fn := range p.pns {
		d := new(bytes.Buffer)
		fn(s, d)
		s = d
	}
	lastPn(s, dest)
	return nil
}

func (p *Pipelines) pop() PipelineFn {
	if len(p.pns) == 0 {
		return nil
	}

	lastIdx := len(p.pns) - 1
	last := p.pns[lastIdx]
	p.pns = p.pns[:lastIdx]
	return last
}
