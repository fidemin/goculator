package goculator

type TokenType string

const (
	TokenTypeNUM    TokenType = "NUM"
	TokenTypeVAR    TokenType = "VAR"
	TokenTypePLUS   TokenType = "PLUS"
	TokenTypeMINUS  TokenType = "MINUS"
	TokenTypeMULTI  TokenType = "MULTI"
	TokenTypeDIV    TokenType = "DIV"
	TokenTypeEOF    TokenType = "EOF"
	TokenTypeNONE   TokenType = "NONE"
	TokenTypeLPARAN TokenType = "LPARAN"
	TokenTypeRPARAN TokenType = "RPARAN"
)

type Token struct {
	Type  TokenType
	Value string
}
