// Package docs contém a especificação Swagger 2.0 gerada automaticamente pelo
// swaggo/swag a partir das anotações nos handlers HTTP.
//
// Para regenerar o swagger.json após alterar as anotações, execute:
//
//	swag init -g cmd/server/main.go -o docs --parseDependency --parseInternal --outputTypes json
package docs

import _ "embed"

// SwaggerJSON é a especificação Swagger 2.0 embutida no binário.
//
//go:embed swagger.json
var SwaggerJSON []byte
