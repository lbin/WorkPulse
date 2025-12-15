package router

import (
	"net/http"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gin-gonic/gin"

	"workpulse/internal/middleware"
)

type RouteSpec struct {
	Method     string
	Path       string
	Permission string
	Summary    string
	Tag        string
	Handler    gin.HandlerFunc
}

func registerRoutes(api *gin.RouterGroup, specs []RouteSpec) {
	for _, spec := range specs {
		handlers := []gin.HandlerFunc{}
		if spec.Permission != "" {
			handlers = append(handlers, middleware.RequirePermission(spec.Permission))
		}
		handlers = append(handlers, spec.Handler)

		switch strings.ToUpper(spec.Method) {
		case http.MethodGet:
			api.GET(spec.Path, handlers...)
		case http.MethodPost:
			api.POST(spec.Path, handlers...)
		case http.MethodPut:
			api.PUT(spec.Path, handlers...)
		case http.MethodDelete:
			api.DELETE(spec.Path, handlers...)
		default:
			api.Any(spec.Path, handlers...)
		}
	}
}

func openAPISpec(basePath string, specs []RouteSpec) *openapi3.T {
	paths := openapi3.Paths{}
	components := openapi3.Components{
		SecuritySchemes: openapi3.SecuritySchemes{
			"bearerAuth": &openapi3.SecuritySchemeRef{
				Value: &openapi3.SecurityScheme{
					Type:         "http",
					Scheme:       "bearer",
					BearerFormat: "JWT",
				},
			},
		},
	}

	doc := &openapi3.T{
		OpenAPI: "3.0.3",
		Info: &openapi3.Info{
			Title:       "WorkPulse API",
			Description: "Auto-generated OpenAPI spec derived from registered Gin routes.",
			Version:     "1.0.0",
		},
		Servers:    openapi3.Servers{{URL: basePath}},
		Paths:      &paths,
		Components: &components,
		Security:   openapi3.SecurityRequirements{{"bearerAuth": []string{}}},
	}

	for _, spec := range specs {
		path := normalizePath(spec.Path)
		pathItem := doc.Paths.Value(path)
		if pathItem == nil {
			pathItem = &openapi3.PathItem{}
		}

		responses := openapi3.Responses{}
		responses.Set("200", &openapi3.ResponseRef{Value: openapi3.NewResponse().WithDescription("OK")})

		op := &openapi3.Operation{
			Summary:     spec.Summary,
			Description: spec.Summary,
			Tags:        []string{spec.Tag},
			Responses:   &responses,
		}

		security := openapi3.SecurityRequirements{{"bearerAuth": []string{}}}
		op.Security = &security

		switch strings.ToUpper(spec.Method) {
		case http.MethodGet:
			pathItem.Get = op
		case http.MethodPost:
			pathItem.Post = op
		case http.MethodPut:
			pathItem.Put = op
		case http.MethodDelete:
			pathItem.Delete = op
		}

		doc.Paths.Set(path, pathItem)
	}

	return doc
}

func serveOpenAPI(basePath string, specs []RouteSpec) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, openAPISpec(basePath, specs))
	}
}

func serveGraphQLSDL(specs []RouteSpec) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.String(http.StatusOK, generateSDL(specs))
	}
}

func normalizePath(path string) string {
	parts := strings.Split(path, "/")
	for i, p := range parts {
		if strings.HasPrefix(p, ":") {
			parts[i] = "{" + strings.TrimPrefix(p, ":") + "}"
		}
	}
	cleaned := strings.Join(parts, "/")
	if !strings.HasPrefix(cleaned, "/") {
		cleaned = "/" + cleaned
	}
	return cleaned
}

func generateSDL(specs []RouteSpec) string {
	var queries []string
	var mutations []string
	added := map[string]struct{}{}

	for _, spec := range specs {
		field := routeToFieldName(spec)
		if field == "" {
			continue
		}
		line := "  " + field + ": JSON"
		if _, exists := added[spec.Method+field]; exists {
			continue
		}
		added[spec.Method+field] = struct{}{}

		if strings.ToUpper(spec.Method) == http.MethodGet {
			queries = append(queries, line)
		} else {
			mutations = append(mutations, line)
		}
	}

	if len(queries) == 0 {
		queries = append(queries, "  _empty: JSON")
	}
	if len(mutations) == 0 {
		mutations = append(mutations, "  _empty: JSON")
	}

	return "scalar JSON\n\ntype Query {\n" + strings.Join(queries, "\n") + "\n}\n\ntype Mutation {\n" + strings.Join(mutations, "\n") + "\n}\n"
}

func routeToFieldName(spec RouteSpec) string {
	cleaned := strings.Trim(spec.Path, "/")
	if cleaned == "" {
		cleaned = "root"
	}
	parts := strings.Split(cleaned, "/")
	var builder []string
	for _, p := range parts {
		if strings.HasPrefix(p, ":") {
			continue
		}
		builder = append(builder, strings.Title(strings.ReplaceAll(p, "-", "_")))
	}

	if len(builder) == 0 {
		builder = append(builder, "Root")
	}

	name := strings.Join(builder, "")
	prefix := ""
	switch strings.ToUpper(spec.Method) {
	case http.MethodPost:
		prefix = "Create"
	case http.MethodPut:
		prefix = "Update"
	case http.MethodDelete:
		prefix = "Delete"
	}

	full := prefix + name
	return strings.ToLower(full[:1]) + full[1:]
}
