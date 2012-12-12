package got

import (
	"encoding/json"
	"io"
)

type Encoder struct {
	encoder *json.Encoder
}

func Encode(p *Project, w io.Writer) error {
	return NewEncoder(w).Encode(p)
}
func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{
		encoder: json.NewEncoder(w),
	}
}

// Overwrite fields of p with values from the reader 
func (enc *Encoder) Encode(p *Project) error {
	return enc.encoder.Encode(p)
}
