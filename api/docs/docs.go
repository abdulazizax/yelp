// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

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
        "/send-verification-code": {
            "post": {
                "description": "Sends a verification code to the user's email for verification purposes",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Send Verification Code",
                "parameters": [
                    {
                        "description": "Send Verification Code Request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/github_com_abdulazizax_yelp_internal_entity_user.SendVerificationCodeRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Verification code sent successfully",
                        "schema": {
                            "$ref": "#/definitions/github_com_abdulazizax_yelp_internal_entity_user.SendVerificationCodeResponse"
                        }
                    },
                    "400": {
                        "description": "Invalid user data",
                        "schema": {
                            "$ref": "#/definitions/github_com_abdulazizax_yelp_internal_entity.Error"
                        }
                    },
                    "404": {
                        "description": "User not found",
                        "schema": {
                            "$ref": "#/definitions/github_com_abdulazizax_yelp_internal_entity.Error"
                        }
                    },
                    "500": {
                        "description": "Failed to send verification code",
                        "schema": {
                            "$ref": "#/definitions/github_com_abdulazizax_yelp_internal_entity.Error"
                        }
                    }
                }
            }
        },
        "/sign-in": {
            "post": {
                "description": "Authenticates a user and returns a JWT token upon successful login",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "User Sign In",
                "parameters": [
                    {
                        "description": "Sign In Request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/github_com_abdulazizax_yelp_internal_entity_user.SignInRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "JWT Token",
                        "schema": {
                            "$ref": "#/definitions/gin.H"
                        }
                    },
                    "400": {
                        "description": "Invalid user data",
                        "schema": {
                            "$ref": "#/definitions/github_com_abdulazizax_yelp_internal_entity.Error"
                        }
                    },
                    "401": {
                        "description": "Incorrect password",
                        "schema": {
                            "$ref": "#/definitions/github_com_abdulazizax_yelp_internal_entity.Error"
                        }
                    },
                    "404": {
                        "description": "User not found",
                        "schema": {
                            "$ref": "#/definitions/github_com_abdulazizax_yelp_internal_entity.Error"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/github_com_abdulazizax_yelp_internal_entity.Error"
                        }
                    }
                }
            }
        },
        "/sign-up": {
            "post": {
                "description": "Registers a new user with the provided details and returns a confirmation message",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Create a new user account",
                "parameters": [
                    {
                        "description": "User registration details",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/github_com_abdulazizax_yelp_internal_entity_user.CreateUser"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "User successfully registered",
                        "schema": {
                            "$ref": "#/definitions/github_com_abdulazizax_yelp_internal_entity.Info"
                        }
                    },
                    "400": {
                        "description": "Invalid request payload",
                        "schema": {
                            "$ref": "#/definitions/github_com_abdulazizax_yelp_internal_entity.Error"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/github_com_abdulazizax_yelp_internal_entity.Error"
                        }
                    }
                }
            }
        },
        "/update-password": {
            "post": {
                "description": "Updates a user's password after validating the request",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Update User Password",
                "parameters": [
                    {
                        "description": "Update User Password Request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/github_com_abdulazizax_yelp_internal_entity_user.UpdateUserPassword"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Password updated successfully",
                        "schema": {
                            "$ref": "#/definitions/github_com_abdulazizax_yelp_internal_entity.Info"
                        }
                    },
                    "400": {
                        "description": "Invalid user data",
                        "schema": {
                            "$ref": "#/definitions/github_com_abdulazizax_yelp_internal_entity.Error"
                        }
                    },
                    "500": {
                        "description": "Failed to update user password",
                        "schema": {
                            "$ref": "#/definitions/github_com_abdulazizax_yelp_internal_entity.Error"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "gin.H": {
            "type": "object",
            "additionalProperties": {}
        },
        "github_com_abdulazizax_yelp_internal_entity.Error": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        },
        "github_com_abdulazizax_yelp_internal_entity.Info": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        },
        "github_com_abdulazizax_yelp_internal_entity_user.CreateSession": {
            "type": "object",
            "required": [
                "ip_address",
                "platform",
                "user_agent"
            ],
            "properties": {
                "ip_address": {
                    "type": "string"
                },
                "platform": {
                    "type": "string",
                    "enum": [
                        "web",
                        "mobile",
                        "admin_panel"
                    ]
                },
                "user_agent": {
                    "type": "string"
                }
            }
        },
        "github_com_abdulazizax_yelp_internal_entity_user.CreateUser": {
            "type": "object",
            "required": [
                "email",
                "gender",
                "name",
                "password"
            ],
            "properties": {
                "email": {
                    "type": "string",
                    "maxLength": 255
                },
                "gender": {
                    "type": "string",
                    "default": "male",
                    "enum": [
                        "male",
                        "female"
                    ]
                },
                "name": {
                    "type": "string",
                    "maxLength": 100
                },
                "password": {
                    "type": "string"
                }
            }
        },
        "github_com_abdulazizax_yelp_internal_entity_user.SendVerificationCodeRequest": {
            "type": "object",
            "required": [
                "email"
            ],
            "properties": {
                "email": {
                    "type": "string",
                    "maxLength": 255
                }
            }
        },
        "github_com_abdulazizax_yelp_internal_entity_user.SendVerificationCodeResponse": {
            "type": "object",
            "properties": {
                "duration": {
                    "type": "string"
                },
                "message": {
                    "type": "string"
                }
            }
        },
        "github_com_abdulazizax_yelp_internal_entity_user.SignInRequest": {
            "type": "object",
            "properties": {
                "create_session": {
                    "$ref": "#/definitions/github_com_abdulazizax_yelp_internal_entity_user.CreateSession"
                },
                "user": {
                    "$ref": "#/definitions/github_com_abdulazizax_yelp_internal_entity_user.SignInUser"
                }
            }
        },
        "github_com_abdulazizax_yelp_internal_entity_user.SignInUser": {
            "type": "object",
            "required": [
                "email",
                "password"
            ],
            "properties": {
                "email": {
                    "type": "string",
                    "maxLength": 255
                },
                "password": {
                    "type": "string"
                }
            }
        },
        "github_com_abdulazizax_yelp_internal_entity_user.UpdateUserPassword": {
            "type": "object",
            "required": [
                "email",
                "new_password",
                "verification_code"
            ],
            "properties": {
                "email": {
                    "type": "string",
                    "maxLength": 255
                },
                "new_password": {
                    "type": "string"
                },
                "verification_code": {
                    "type": "string"
                }
            }
        }
    },
    "securityDefinitions": {
        "BearerAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.03.67.83.145",
	Host:             "localhost:8080",
	BasePath:         "/",
	Schemes:          []string{"http", "https"},
	Title:            "# UdevsLab Homework3",
	Description:      "API Endpoints for MiniTwitter",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
