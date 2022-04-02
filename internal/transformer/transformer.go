package transformer

import "bytes"

type TransformerHandler interface {
	Transform(src *bytes.Buffer) error
}

type Transformer func(next TransformerFunc) TransformerFunc

type TransformerFunc func(data *bytes.Buffer) error

func (tf TransformerFunc) Transform(data *bytes.Buffer) error {
	return tf(data)
}

// ApplyTransformers apply TransformerChain in reversed order so that `tf` will handle it in the end
func ApplyTransformers(tf TransformerFunc, chains ...Transformer) TransformerFunc {
	for i := len(chains) - 1; i >= 0; i-- {
		tf = chains[i](tf)
	}
	return tf
}
