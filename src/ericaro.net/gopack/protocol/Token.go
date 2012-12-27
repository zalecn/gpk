package protocol

import (
	"encoding/base64"
)

// Represent an identification Token
type Token []byte

func (t *Token) Format() string {
	return base64.URLEncoding.EncodeToString(*t)
}

func DecodeString(v string) (t  *Token, err error) {
	b, err := base64.URLEncoding.DecodeString(v)
	token := Token(b)
	return &token, err
}
