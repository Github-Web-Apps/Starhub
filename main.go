package main

import (
	"net/http"

	rice "github.com/GeertJohan/go.rice"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
)

func main() {
	// ts := oauth2.StaticTokenSource(
	// 	&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	// )
	// tc := oauth2.NewClient(oauth2.NoContext, ts)
	// client := github.NewClient(tc)
	// user := os.Args[1]
	// log.Println("Gathering data for", user)

	// followers, err := followers.Get(user, client)
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// log.Println("You have a total of", len(followers), "followers!")
	// for _, follower := range followers {
	// 	log.Println(*follower.Login)
	// }

	e := echo.New()
	// the file server for rice. "app" is the folder where the files come from.
	assetHandler := http.FileServer(rice.MustFindBox("static").HTTPBox())
	// serves the index.html from rice
	e.GET("/", standard.WrapHandler(assetHandler))

	// servers other static files
	e.GET("/static/*", standard.WrapHandler(http.StripPrefix("/static/", assetHandler)))
	e.Run(standard.New(":3000"))
}
