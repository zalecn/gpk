package protocol

import (
	"encoding/base64"
)

// Represent an identification Token
// any []byte is ok. usually it is a fixed size []byte. 
type Token []byte

//FormatURL format this token in a string suitable for url encoding. Mainly it uses base64.URLEncoding instead of std one
func (t *Token) FormatURL() string {
	return base64.URLEncoding.EncodeToString(*t)
}

//FormatStd format this token in a string
func (t *Token) FormatStd() string {
	return base64.StdEncoding.EncodeToString(*t)
}
//ParseURLToken read a Token in a string
func ParseURLToken(v string) (t  *Token, err error) {
	b, err := base64.URLEncoding.DecodeString(v)
	token := Token(b)
	return &token, err
}
//ParseStdToken read a Token in a string using base64.URLEncoding
func ParseStdToken(v string) (t  *Token, err error) {
	b, err := base64.StdEncoding.DecodeString(v)
	token := Token(b)
	return &token, err
}