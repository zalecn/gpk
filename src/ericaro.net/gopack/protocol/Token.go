package protocol

import (
	"encoding/base64"
)

// Represent an identification Token
type Token []byte

func (t *Token) FormatURL() string {
	return base64.URLEncoding.EncodeToString(*t)
}
func (t *Token) FormatStd() string {
	return base64.StdEncoding.EncodeToString(*t)
}
func ParseURLToken(v string) (t  *Token, err error) {
	b, err := base64.URLEncoding.DecodeString(v)
	token := Token(b)
	return &token, err
}

func ParseStdToken(v string) (t  *Token, err error) {
	b, err := base64.StdEncoding.DecodeString(v)
	token := Token(b)
	return &token, err
}