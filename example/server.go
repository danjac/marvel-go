package main

import (
	"github.com/danjac/marvel-go/marvel"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"log"
	"net/http"
	"os"
)

func main() {

	public_key := os.Getenv("MARVEL_PUBLIC_KEY")

	if public_key == "" {
		log.Fatal("MARVEL_PUBLIC_KEY must be set in environment")
	}

	private_key := os.Getenv("MARVEL_PRIVATE_KEY")

	if private_key == "" {
		log.Fatal("MARVEL_PRIVATE_KEY must be set in environment")
	}

	client := marvel.NewClient(public_key, private_key)
	m := martini.Classic()

	m.Use(render.Renderer(render.Options{IndentJSON: true}))

	m.Get("/api/comics", func(r render.Render) {
		params := &marvel.ComicQueryParams{Limit: 20}
		comics, err := client.GetComics(params)
		if err != nil {
			log.Fatal(err)
		}
		r.JSON(http.StatusOK, comics)
	})

	m.Get("/api/comics/:comicId", func(params martini.Params, r render.Render) {
		comic, err := client.GetComic(params["comicId"])
		if err != nil {
			log.Fatal(err)
		}
		if comic == nil {
			r.JSON(http.StatusNotFound, nil)
		}
		r.JSON(http.StatusOK, comic)
	})

	m.Run()
}
