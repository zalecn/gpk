package got

import (
	"encoding/json"
	"io"
)

type Decoder struct {
	decoder *json.Decoder
}

func Decode(p *Project, r io.Reader) error {
	return NewDecoder(r).Decode(p)
}
func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{
		decoder: json.NewDecoder(r),
	}
}

// Overwrite fields of p with values from the reader 
func (dec *Decoder) Decode(p *Project) error {
	return dec.decoder.Decode(p)
}
