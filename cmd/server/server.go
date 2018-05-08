package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// Member is part of an album
type Member struct {
	Name  string
	IsDir bool
}

// List is a list of a directory
type List struct {
	Members []Member
}

var port = flag.Int("port", 1323, "port")
var root = flag.String("root", ".", "directory root")
var db = flag.String("db", "/var/tmp/sqlite.db", "path to database")

func list(c echo.Context) error {
	path := c.Param("path")
	// TODO(tstromberg): Make less of a security hole
	files, err := ioutil.ReadDir(filepath.Join(*root, path))
	if err != nil {
		return err
	}
	l := List{}
	for _, f := range files {
		l.Members = append(l.Members, Member{Name: f.Name(), IsDir: f.IsDir()})
	}

	return c.JSON(http.StatusOK, l)
}

func main() {
	flag.Parse()

	e := echo.New()
	// Reveal error messages.
	e.Debug = true
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.GET("/album/list/*", list)
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", *port)))
}
