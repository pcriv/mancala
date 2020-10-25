package restapi

import (
	"bytes"
	"html/template"
	"net/http"
	"path"

	"github.com/labstack/echo/v4"
)

// SwaggerUIConfig configures the SwaggerUI middleware
type SwaggerUIConfig struct {
	// BasePath for the UI path, defaults to: /
	BasePath string
	// Path combines with BasePath for the full UI path, defaults to: docs
	Path string
	// SpecURL the url to find the spec for
	SpecURL string

	// The three components needed to embed swagger-ui
	SwaggerURL       string
	SwaggerPresetURL string
	SwaggerStylesURL string

	Favicon32 string
	Favicon16 string

	// Title for the documentation site, default to: API documentation
	Title string
}

var DefaultSwaggerUIConfig = SwaggerUIConfig{
	BasePath:         "/",
	Path:             "docs",
	SpecURL:          "/swagger.json",
	SwaggerURL:       swaggerLatest,
	SwaggerPresetURL: swaggerPresetLatest,
	SwaggerStylesURL: swaggerStylesLatest,
	Favicon16:        swaggerFavicon16Latest,
	Favicon32:        swaggerFavicon32Latest,
	Title:            "API documentation",
}

// EnsureDefaults in case some options are missing
func (r *SwaggerUIConfig) EnsureDefaults() {
	if r.BasePath == "" {
		r.BasePath = DefaultSwaggerUIConfig.BasePath
	}
	if r.Path == "" {
		r.Path = DefaultSwaggerUIConfig.Path
	}
	if r.SpecURL == "" {
		r.SpecURL = DefaultSwaggerUIConfig.SpecURL
	}
	if r.SwaggerURL == "" {
		r.SwaggerURL = DefaultSwaggerUIConfig.SwaggerURL
	}
	if r.SwaggerPresetURL == "" {
		r.SwaggerPresetURL = DefaultSwaggerUIConfig.SwaggerPresetURL
	}
	if r.SwaggerStylesURL == "" {
		r.SwaggerStylesURL = DefaultSwaggerUIConfig.SwaggerStylesURL
	}
	if r.Favicon16 == "" {
		r.Favicon16 = DefaultSwaggerUIConfig.Favicon16
	}
	if r.Favicon32 == "" {
		r.Favicon32 = DefaultSwaggerUIConfig.Favicon32
	}
	if r.Title == "" {
		r.Title = DefaultSwaggerUIConfig.Title
	}
}

func SwaggerUI() echo.MiddlewareFunc {
	return SwaggerUIWithConfig(DefaultSwaggerUIConfig)
}

// SwaggerUI creates a middleware to serve a documentation site for a swagger spec.
// This allows for altering the spec before starting the http listener.
func SwaggerUIWithConfig(config SwaggerUIConfig) echo.MiddlewareFunc {
	config.EnsureDefaults()

	pth := path.Join(config.BasePath, config.Path)
	tmpl := template.Must(template.New("swaggerui").Parse(swaggerUITemplate))

	buf := bytes.NewBuffer(nil)
	_ = tmpl.Execute(buf, &config)
	b := buf.Bytes()

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if c.Request().URL.Path == pth {
				c.Response().Header().Set("Content-Type", "text/html; charset=utf-8")
				c.Response().WriteHeader(http.StatusOK)

				_, _ = c.Response().Write(b)
			}
			return next(c)
		}
	}
}

const (
	swaggerLatest          = "https://unpkg.com/swagger-ui-dist/swagger-ui-bundle.js"
	swaggerPresetLatest    = "https://unpkg.com/swagger-ui-dist/swagger-ui-standalone-preset.js"
	swaggerStylesLatest    = "https://unpkg.com/swagger-ui-dist/swagger-ui.css"
	swaggerFavicon32Latest = "https://unpkg.com/swagger-ui-dist/favicon-32x32.png"
	swaggerFavicon16Latest = "https://unpkg.com/swagger-ui-dist/favicon-16x16.png"
	swaggerUITemplate      = `
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8">
		<title>{{ .Title }}</title>
    <link rel="stylesheet" type="text/css" href="{{ .SwaggerStylesURL }}" >
    <link rel="icon" type="image/png" href="{{ .Favicon32 }}" sizes="32x32" />
    <link rel="icon" type="image/png" href="{{ .Favicon16 }}" sizes="16x16" />
    <style>
      html
      {
        box-sizing: border-box;
        overflow: -moz-scrollbars-vertical;
        overflow-y: scroll;
      }
      *,
      *:before,
      *:after
      {
        box-sizing: inherit;
      }
      body
      {
        margin:0;
        background: #fafafa;
      }
    </style>
  </head>
  <body>
    <div id="swagger-ui"></div>
    <script src="{{ .SwaggerURL }}"> </script>
    <script src="{{ .SwaggerPresetURL }}"> </script>
    <script>
    window.onload = function() {
      // Begin Swagger UI call region
      const ui = SwaggerUIBundle({
        url: '{{ .SpecURL }}',
        dom_id: '#swagger-ui',
        deepLinking: true,
        presets: [
          SwaggerUIBundle.presets.apis,
          SwaggerUIStandalonePreset
        ],
        plugins: [
          SwaggerUIBundle.plugins.DownloadUrl
        ],
        layout: "StandaloneLayout"
      })
      // End Swagger UI call region
      window.ui = ui
    }
  </script>
  </body>
</html>
`
)
