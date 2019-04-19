package goculator

// TokenType is specific types of token.
type TokenType string

const (
	// TokenTypeNUM represents token with float string value
	TokenTypeNUM TokenType = "NUM"
	// TokenTypeVAR represents token with variable name
	TokenTypeVAR TokenType = "VAR"
	// TokenTypePLUS represents token with "+" character
	TokenTypePLUS TokenType = "PLUS"
	// TokenTypePLUS represents token with "-" character
	TokenTypeMINUS TokenType = "MINUS"
	// TokenTypePLUS represents token with "*" character
	TokenTypeMULTI TokenType = "MULTI"
	// TokenTypePLUS represents token with "/" character
	TokenTypeDIV TokenType = "DIV"
	// TokenTypePLUS represents EOF token.
	TokenTypeEOF TokenType = "EOF"
	// TokenTypeLPARAN represents token with "("
	TokenTypeLPARAN TokenType = "LPARAN"
	// TokenTypeLPARAN represents token with ")"
	TokenTypeRPARAN TokenType = "RPARAN"
	// TokenTypeNone represents token with value which cannot be tokenized.
	TokenTypeNONE TokenType = "NONE"
)

// Token is the token used for the calculator.
type Token struct {
	Type  TokenType
	Value string
}
