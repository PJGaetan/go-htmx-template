package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"text/template"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Template struct {
	templates *template.Template
}

type GlobalState struct {
	Count int
}

var global GlobalState

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	fmt.Print(data)
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	// Counter per session
	// https://templ.guide/server-side-rendering/example-counter-application
	// Using : https://github.com/alexedwards/scs

	// Echo instance
	e := echo.New()
	// serve static file instead of using cdn
	e.Static("/", "web/css")

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	t := &Template{
		templates: template.Must(template.ParseGlob("web/templates/*.html")),
	}
	e.Renderer = t

	// Routes
	// build in templating
	e.GET("/hello", hello)
	e.POST("/clicked", button)
	e.POST("/reset", reset)

	// go-templ
	e.GET("/templ", func(c echo.Context) error {
		component := IndexTempl(global.Count)
		h := templ.Handler(component)
		return h.Component.Render(context.Background(), c.Response().Writer)
	})
	e.POST("/button-tmpl-add", func(c echo.Context) error {
		global.Count += 1
		component := ButtonTempl(global.Count)
		h := templ.Handler(component)
		return h.Component.Render(context.Background(), c.Response().Writer)
	})

	e.POST("/button-tmpl-reset", func(c echo.Context) error {
		global.Count = 0
		component := ButtonTempl(global.Count)
		h := templ.Handler(component)
		return h.Component.Render(context.Background(), c.Response().Writer)
	})

	// Mix of both by calling /clicked and /reset
	// from buttonTempl and only rerender the full templ when needed

	// Start server
	e.Logger.Fatal(e.Start(":1234"))
}

// Handler
func hello(c echo.Context) error {
	return c.Render(http.StatusOK, "index.html", global)
}

func button(c echo.Context) error {
	global.Count += 1
	return c.String(http.StatusOK, strconv.Itoa(global.Count))
}

func reset(c echo.Context) error {
	global.Count = 0
	return c.String(http.StatusOK, strconv.Itoa(global.Count))
}
