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
	lastPn := p.getLastPipelineFn()
	p.pns = p.pns[:l-1] // remove last PipelineFn
	for _, fn := range p.pns {
		d := new(bytes.Buffer)
		fn(s, d)
		s = d
	}
	lastPn(s, dest)
	return nil
}

func (p *Pipelines) getLastPipelineFn() PipelineFn {
	if len(p.pns) == 0 {
		return nil
	}

	return p.pns[len(p.pns)-1]
}
