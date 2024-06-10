// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import (
	"github.com/swaggo/swag"
	"io/ioutil"
)

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "API Support",
            "url": "http://www.swagger.io/support",
            "email": "support@swagger.io"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/login": {
            "post": {
                "description": "get string by ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "accounts"
                ],
                "summary": "Show an account",
                "parameters": [
                    {
                        "description": "Login Data",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/controllers.LoginInput"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/controllers.LoginResponse"
                        }
                    }
                }
            }
        },
        "/p/user": {
            "get": {
                "description": "get string by ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User Authentication"
                ],
                "summary": "Get current app user",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/controllers.UserResponse"
                        }
                    }
                }
            }
        },
        "/register": {
            "post": {
                "description": "Register new users",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Register",
                "parameters": [
                    {
                        "description": "Register Data",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/controllers.RegisterByPhoneInput"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/controllers.UserResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "controllers.LoginInput": {
            "type": "object",
            "required": [
                "password",
                "phone"
            ],
            "properties": {
                "password": {
                    "type": "string"
                },
                "phone": {
                    "type": "string"
                }
            }
        },
        "controllers.LoginResponse": {
            "type": "object",
            "properties": {
                "access_token": {
                    "type": "string"
                },
                "expired_date": {
                    "type": "string"
                },
                "refresh_token": {
                    "type": "string"
                },
                "type": {
                    "type": "string"
                }
            }
        },
        "controllers.RegisterByPhoneInput": {
            "type": "object",
            "required": [
                "fio",
                "iin",
                "password",
                "phone"
            ],
            "properties": {
                "fio": {
                    "type": "string"
                },
                "iin": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                },
                "phone": {
                    "type": "string"
                }
            }
        },
        "controllers.UserResponse": {
            "type": "object",
            "properties": {
                "fio": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "iin": {
                    "type": "string"
                },
                "phone": {
                    "type": "string"
                },
                "type": {
                    "type": "string"
                }
            }
        }
    },
    "securityDefinitions": {
        "BearerAuth": {
            "type": "basic"
        }
    },
    "externalDocs": {
        "description": "OpenAPI",
        "url": "https://swagger.io/resources/open-api/"
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/api",
	Schemes:          []string{},
	Title:            "Swagger Example API",
	Description:      "This is a sample server celler server.",
	InfoInstanceName: "swagger",
	//SwaggerTemplate:  docTemplate,
	//LeftDelim:  "{{",
	//RightDelim: "}}",
}

func init() {
	doc, err := ioutil.ReadFile("./docs/swagger.yaml")
	if err != nil {
		panic(err.Error())
	}

	SwaggerInfo.SwaggerTemplate = string(doc)
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}