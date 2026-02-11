package scaffold

import (
	"fmt"
	"path/filepath"
	"strings"
)

func writeScalarFiles(root string, opts Options) error {
	apiRoot := filepath.Join(root, "apps", "api")
	module := opts.ProjectName + "/apps/api"

	files := map[string]string{
		filepath.Join(apiRoot, "internal", "docs", "docs.go"):    apiScalarHandlerGo(opts),
		filepath.Join(apiRoot, "internal", "docs", "openapi.go"): apiOpenAPISpecGo(opts),
	}

	for path, content := range files {
		content = strings.ReplaceAll(content, "{{MODULE}}", module)
		if err := writeFile(path, content); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
	}

	return nil
}

func apiScalarHandlerGo(opts Options) string {
	return fmt.Sprintf(`package docs

import (
	"net/http"

	"github.com/gin-gonic/gin"
	scalar "github.com/MarceloPetrucio/go-scalar-api-reference"
)

// RegisterRoutes mounts the API documentation endpoints.
//   - GET /docs          → Scalar interactive API reference
//   - GET /docs/openapi.json → OpenAPI 3.0 JSON spec
func RegisterRoutes(r *gin.Engine) {
	r.GET("/docs/openapi.json", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.String(http.StatusOK, OpenAPISpec())
	})

	r.GET("/docs", func(c *gin.Context) {
		htmlContent, err := scalar.ApiReferenceHTML(&scalar.Options{
			SpecURL: "/docs/openapi.json",
			CustomOptions: scalar.CustomOptions{
				PageTitle: "%s — API Reference",
			},
			DarkMode: true,
		})
		if err != nil {
			c.String(http.StatusInternalServerError, "Error generating docs: "+err.Error())
			return
		}
		c.Header("Content-Type", "text/html")
		c.String(http.StatusOK, htmlContent)
	})
}
`, opts.ProjectName)
}

func apiOpenAPISpecGo(opts Options) string {
	return fmt.Sprintf(`package docs

import "fmt"

// OpenAPISpec returns the OpenAPI 3.0 JSON specification for the API.
func OpenAPISpec() string {
	return fmt.Sprintf(%[1]s{
  "openapi": "3.0.3",
  "info": {
    "title": "%%s API",
    "description": "REST API built with Grit — Go + React meta-framework. Powered by Gin + GORM.",
    "version": "1.0.0",
    "contact": {
      "name": "API Support"
    }
  },
  "servers": [
    {
      "url": "http://localhost:8080",
      "description": "Local development"
    }
  ],
  "tags": [
    { "name": "Health", "description": "Service health check" },
    { "name": "Authentication", "description": "Register, login, token refresh, and password reset" },
    { "name": "Users", "description": "User management (admin)" },
    { "name": "Uploads", "description": "File uploads and management" },
    { "name": "AI", "description": "AI text generation (Claude / OpenAI)" },
    { "name": "System", "description": "Background jobs, cron tasks, and admin tools" }
  ],
  "paths": {
    "/api/health": {
      "get": {
        "tags": ["Health"],
        "summary": "Health check",
        "description": "Returns the API status and version.",
        "responses": {
          "200": {
            "description": "Service is healthy",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": {
                    "status": { "type": "string", "example": "ok" },
                    "version": { "type": "string", "example": "0.1.0" }
                  }
                }
              }
            }
          }
        }
      }
    },
    "/api/auth/register": {
      "post": {
        "tags": ["Authentication"],
        "summary": "Register a new user",
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": { "$ref": "#/components/schemas/RegisterRequest" }
            }
          }
        },
        "responses": {
          "201": {
            "description": "User created",
            "content": {
              "application/json": {
                "schema": { "$ref": "#/components/schemas/AuthResponse" }
              }
            }
          },
          "400": { "$ref": "#/components/responses/BadRequest" },
          "422": { "$ref": "#/components/responses/ValidationError" }
        }
      }
    },
    "/api/auth/login": {
      "post": {
        "tags": ["Authentication"],
        "summary": "Login with credentials",
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": { "$ref": "#/components/schemas/LoginRequest" }
            }
          }
        },
        "responses": {
          "200": {
            "description": "Login successful",
            "content": {
              "application/json": {
                "schema": { "$ref": "#/components/schemas/AuthResponse" }
              }
            }
          },
          "401": { "$ref": "#/components/responses/Unauthorized" }
        }
      }
    },
    "/api/auth/refresh": {
      "post": {
        "tags": ["Authentication"],
        "summary": "Refresh access token",
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "type": "object",
                "required": ["refresh_token"],
                "properties": {
                  "refresh_token": { "type": "string" }
                }
              }
            }
          }
        },
        "responses": {
          "200": {
            "description": "Token refreshed",
            "content": {
              "application/json": {
                "schema": { "$ref": "#/components/schemas/AuthResponse" }
              }
            }
          },
          "401": { "$ref": "#/components/responses/Unauthorized" }
        }
      }
    },
    "/api/auth/forgot-password": {
      "post": {
        "tags": ["Authentication"],
        "summary": "Request password reset",
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "type": "object",
                "required": ["email"],
                "properties": {
                  "email": { "type": "string", "format": "email" }
                }
              }
            }
          }
        },
        "responses": {
          "200": {
            "description": "Reset email sent (always returns success for security)"
          }
        }
      }
    },
    "/api/auth/reset-password": {
      "post": {
        "tags": ["Authentication"],
        "summary": "Reset password with token",
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "type": "object",
                "required": ["token", "password"],
                "properties": {
                  "token": { "type": "string" },
                  "password": { "type": "string", "minLength": 8 }
                }
              }
            }
          }
        },
        "responses": {
          "200": { "description": "Password reset successful" },
          "400": { "$ref": "#/components/responses/BadRequest" }
        }
      }
    },
    "/api/auth/me": {
      "get": {
        "tags": ["Authentication"],
        "summary": "Get current user",
        "security": [{ "bearerAuth": [] }],
        "responses": {
          "200": {
            "description": "Current user profile",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": {
                    "data": { "$ref": "#/components/schemas/User" }
                  }
                }
              }
            }
          },
          "401": { "$ref": "#/components/responses/Unauthorized" }
        }
      }
    },
    "/api/auth/logout": {
      "post": {
        "tags": ["Authentication"],
        "summary": "Logout (invalidate tokens)",
        "security": [{ "bearerAuth": [] }],
        "responses": {
          "200": { "description": "Logged out successfully" }
        }
      }
    },
    "/api/users": {
      "get": {
        "tags": ["Users"],
        "summary": "List users (admin)",
        "description": "Paginated list with search, sort, and filter. Requires admin role.",
        "security": [{ "bearerAuth": [] }],
        "parameters": [
          { "name": "page", "in": "query", "schema": { "type": "integer", "default": 1 } },
          { "name": "page_size", "in": "query", "schema": { "type": "integer", "default": 20 } },
          { "name": "search", "in": "query", "schema": { "type": "string" } },
          { "name": "sort", "in": "query", "schema": { "type": "string", "example": "created_at" } },
          { "name": "order", "in": "query", "schema": { "type": "string", "enum": ["asc", "desc"] } }
        ],
        "responses": {
          "200": {
            "description": "Paginated user list",
            "content": {
              "application/json": {
                "schema": { "$ref": "#/components/schemas/UserListResponse" }
              }
            }
          },
          "401": { "$ref": "#/components/responses/Unauthorized" },
          "403": { "$ref": "#/components/responses/Forbidden" }
        }
      }
    },
    "/api/users/{id}": {
      "get": {
        "tags": ["Users"],
        "summary": "Get user by ID",
        "security": [{ "bearerAuth": [] }],
        "parameters": [
          { "name": "id", "in": "path", "required": true, "schema": { "type": "integer" } }
        ],
        "responses": {
          "200": {
            "description": "User details",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": { "data": { "$ref": "#/components/schemas/User" } }
                }
              }
            }
          },
          "404": { "$ref": "#/components/responses/NotFound" }
        }
      },
      "put": {
        "tags": ["Users"],
        "summary": "Update user",
        "security": [{ "bearerAuth": [] }],
        "parameters": [
          { "name": "id", "in": "path", "required": true, "schema": { "type": "integer" } }
        ],
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": { "$ref": "#/components/schemas/UpdateUserRequest" }
            }
          }
        },
        "responses": {
          "200": {
            "description": "User updated",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": { "data": { "$ref": "#/components/schemas/User" } }
                }
              }
            }
          },
          "404": { "$ref": "#/components/responses/NotFound" }
        }
      },
      "delete": {
        "tags": ["Users"],
        "summary": "Delete user (admin)",
        "security": [{ "bearerAuth": [] }],
        "parameters": [
          { "name": "id", "in": "path", "required": true, "schema": { "type": "integer" } }
        ],
        "responses": {
          "200": { "description": "User deleted" },
          "404": { "$ref": "#/components/responses/NotFound" }
        }
      }
    },
    "/api/uploads": {
      "post": {
        "tags": ["Uploads"],
        "summary": "Upload a file",
        "security": [{ "bearerAuth": [] }],
        "requestBody": {
          "required": true,
          "content": {
            "multipart/form-data": {
              "schema": {
                "type": "object",
                "properties": {
                  "file": { "type": "string", "format": "binary" }
                }
              }
            }
          }
        },
        "responses": {
          "201": {
            "description": "File uploaded",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": { "data": { "$ref": "#/components/schemas/Upload" } }
                }
              }
            }
          },
          "400": { "$ref": "#/components/responses/BadRequest" }
        }
      },
      "get": {
        "tags": ["Uploads"],
        "summary": "List uploads",
        "security": [{ "bearerAuth": [] }],
        "parameters": [
          { "name": "page", "in": "query", "schema": { "type": "integer", "default": 1 } },
          { "name": "page_size", "in": "query", "schema": { "type": "integer", "default": 20 } }
        ],
        "responses": {
          "200": {
            "description": "Paginated upload list",
            "content": {
              "application/json": {
                "schema": { "$ref": "#/components/schemas/UploadListResponse" }
              }
            }
          }
        }
      }
    },
    "/api/uploads/{id}": {
      "get": {
        "tags": ["Uploads"],
        "summary": "Get upload by ID",
        "security": [{ "bearerAuth": [] }],
        "parameters": [
          { "name": "id", "in": "path", "required": true, "schema": { "type": "integer" } }
        ],
        "responses": {
          "200": {
            "description": "Upload details",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": { "data": { "$ref": "#/components/schemas/Upload" } }
                }
              }
            }
          },
          "404": { "$ref": "#/components/responses/NotFound" }
        }
      },
      "delete": {
        "tags": ["Uploads"],
        "summary": "Delete upload",
        "security": [{ "bearerAuth": [] }],
        "parameters": [
          { "name": "id", "in": "path", "required": true, "schema": { "type": "integer" } }
        ],
        "responses": {
          "200": { "description": "Upload deleted" },
          "404": { "$ref": "#/components/responses/NotFound" }
        }
      }
    },
    "/api/ai/complete": {
      "post": {
        "tags": ["AI"],
        "summary": "Generate text completion",
        "security": [{ "bearerAuth": [] }],
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "type": "object",
                "required": ["prompt"],
                "properties": {
                  "prompt": { "type": "string" },
                  "max_tokens": { "type": "integer", "default": 1024 }
                }
              }
            }
          }
        },
        "responses": {
          "200": {
            "description": "AI completion",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": {
                    "data": {
                      "type": "object",
                      "properties": {
                        "text": { "type": "string" }
                      }
                    }
                  }
                }
              }
            }
          }
        }
      }
    },
    "/api/ai/chat": {
      "post": {
        "tags": ["AI"],
        "summary": "Chat with AI",
        "security": [{ "bearerAuth": [] }],
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "type": "object",
                "required": ["messages"],
                "properties": {
                  "messages": {
                    "type": "array",
                    "items": {
                      "type": "object",
                      "properties": {
                        "role": { "type": "string", "enum": ["user", "assistant", "system"] },
                        "content": { "type": "string" }
                      }
                    }
                  },
                  "max_tokens": { "type": "integer", "default": 1024 }
                }
              }
            }
          }
        },
        "responses": {
          "200": {
            "description": "Chat response",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": {
                    "data": {
                      "type": "object",
                      "properties": {
                        "text": { "type": "string" }
                      }
                    }
                  }
                }
              }
            }
          }
        }
      }
    },
    "/api/ai/stream": {
      "post": {
        "tags": ["AI"],
        "summary": "Stream AI response (SSE)",
        "description": "Returns a Server-Sent Events stream of text chunks.",
        "security": [{ "bearerAuth": [] }],
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "type": "object",
                "required": ["prompt"],
                "properties": {
                  "prompt": { "type": "string" },
                  "max_tokens": { "type": "integer", "default": 1024 }
                }
              }
            }
          }
        },
        "responses": {
          "200": {
            "description": "SSE stream of text chunks",
            "content": {
              "text/event-stream": {
                "schema": { "type": "string" }
              }
            }
          }
        }
      }
    },
    "/api/admin/jobs/stats": {
      "get": {
        "tags": ["System"],
        "summary": "Get job queue statistics (admin)",
        "security": [{ "bearerAuth": [] }],
        "responses": {
          "200": {
            "description": "Queue stats",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": {
                    "data": { "type": "object" }
                  }
                }
              }
            }
          },
          "403": { "$ref": "#/components/responses/Forbidden" }
        }
      }
    },
    "/api/admin/jobs/{status}": {
      "get": {
        "tags": ["System"],
        "summary": "List jobs by status (admin)",
        "security": [{ "bearerAuth": [] }],
        "parameters": [
          { "name": "status", "in": "path", "required": true, "schema": { "type": "string", "enum": ["active", "pending", "completed", "failed", "retry"] } }
        ],
        "responses": {
          "200": {
            "description": "List of jobs",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": { "data": { "type": "array", "items": { "type": "object" } } }
                }
              }
            }
          },
          "403": { "$ref": "#/components/responses/Forbidden" }
        }
      }
    },
    "/api/admin/jobs/{id}/retry": {
      "post": {
        "tags": ["System"],
        "summary": "Retry a failed job (admin)",
        "security": [{ "bearerAuth": [] }],
        "parameters": [
          { "name": "id", "in": "path", "required": true, "schema": { "type": "string" } }
        ],
        "responses": {
          "200": { "description": "Job queued for retry" },
          "403": { "$ref": "#/components/responses/Forbidden" }
        }
      }
    },
    "/api/admin/cron/tasks": {
      "get": {
        "tags": ["System"],
        "summary": "List cron tasks (admin)",
        "security": [{ "bearerAuth": [] }],
        "responses": {
          "200": {
            "description": "Cron task list",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": { "data": { "type": "array", "items": { "type": "object" } } }
                }
              }
            }
          },
          "403": { "$ref": "#/components/responses/Forbidden" }
        }
      }
    }
  },
  "components": {
    "securitySchemes": {
      "bearerAuth": {
        "type": "http",
        "scheme": "bearer",
        "bearerFormat": "JWT",
        "description": "JWT access token obtained from /api/auth/login"
      }
    },
    "schemas": {
      "User": {
        "type": "object",
        "properties": {
          "id": { "type": "integer", "example": 1 },
          "name": { "type": "string", "example": "John Doe" },
          "email": { "type": "string", "format": "email", "example": "john@example.com" },
          "role": { "type": "string", "enum": ["user", "editor", "admin"], "example": "user" },
          "avatar": { "type": "string", "example": "" },
          "active": { "type": "boolean", "example": true },
          "created_at": { "type": "string", "format": "date-time" },
          "updated_at": { "type": "string", "format": "date-time" }
        }
      },
      "Upload": {
        "type": "object",
        "properties": {
          "id": { "type": "integer" },
          "filename": { "type": "string" },
          "key": { "type": "string" },
          "content_type": { "type": "string" },
          "size": { "type": "integer" },
          "url": { "type": "string" },
          "user_id": { "type": "integer" },
          "created_at": { "type": "string", "format": "date-time" }
        }
      },
      "RegisterRequest": {
        "type": "object",
        "required": ["name", "email", "password"],
        "properties": {
          "name": { "type": "string", "minLength": 2, "example": "John Doe" },
          "email": { "type": "string", "format": "email", "example": "john@example.com" },
          "password": { "type": "string", "minLength": 8, "example": "password123" }
        }
      },
      "LoginRequest": {
        "type": "object",
        "required": ["email", "password"],
        "properties": {
          "email": { "type": "string", "format": "email", "example": "admin@example.com" },
          "password": { "type": "string", "example": "password" }
        }
      },
      "UpdateUserRequest": {
        "type": "object",
        "properties": {
          "name": { "type": "string" },
          "email": { "type": "string", "format": "email" },
          "role": { "type": "string", "enum": ["user", "editor", "admin"] },
          "active": { "type": "boolean" }
        }
      },
      "AuthResponse": {
        "type": "object",
        "properties": {
          "data": {
            "type": "object",
            "properties": {
              "user": { "$ref": "#/components/schemas/User" },
              "tokens": { "$ref": "#/components/schemas/TokenPair" }
            }
          },
          "message": { "type": "string", "example": "Login successful" }
        }
      },
      "TokenPair": {
        "type": "object",
        "properties": {
          "access_token": { "type": "string" },
          "refresh_token": { "type": "string" },
          "expires_at": { "type": "integer", "description": "Unix timestamp" }
        }
      },
      "UserListResponse": {
        "type": "object",
        "properties": {
          "data": {
            "type": "array",
            "items": { "$ref": "#/components/schemas/User" }
          },
          "meta": { "$ref": "#/components/schemas/PaginationMeta" }
        }
      },
      "UploadListResponse": {
        "type": "object",
        "properties": {
          "data": {
            "type": "array",
            "items": { "$ref": "#/components/schemas/Upload" }
          },
          "meta": { "$ref": "#/components/schemas/PaginationMeta" }
        }
      },
      "PaginationMeta": {
        "type": "object",
        "properties": {
          "total": { "type": "integer", "example": 100 },
          "page": { "type": "integer", "example": 1 },
          "page_size": { "type": "integer", "example": 20 },
          "pages": { "type": "integer", "example": 5 }
        }
      },
      "ErrorResponse": {
        "type": "object",
        "properties": {
          "error": {
            "type": "object",
            "properties": {
              "code": { "type": "string", "example": "VALIDATION_ERROR" },
              "message": { "type": "string", "example": "Email is required" },
              "details": { "type": "object" }
            }
          }
        }
      }
    },
    "responses": {
      "BadRequest": {
        "description": "Bad request",
        "content": {
          "application/json": {
            "schema": { "$ref": "#/components/schemas/ErrorResponse" }
          }
        }
      },
      "Unauthorized": {
        "description": "Unauthorized — invalid or missing token",
        "content": {
          "application/json": {
            "schema": { "$ref": "#/components/schemas/ErrorResponse" }
          }
        }
      },
      "Forbidden": {
        "description": "Forbidden — insufficient permissions",
        "content": {
          "application/json": {
            "schema": { "$ref": "#/components/schemas/ErrorResponse" }
          }
        }
      },
      "NotFound": {
        "description": "Resource not found",
        "content": {
          "application/json": {
            "schema": { "$ref": "#/components/schemas/ErrorResponse" }
          }
        }
      },
      "ValidationError": {
        "description": "Validation failed",
        "content": {
          "application/json": {
            "schema": { "$ref": "#/components/schemas/ErrorResponse" }
          }
        }
      }
    }
  }
}%[1]s, "%[2]s")
}
`, "`", opts.ProjectName)
}
