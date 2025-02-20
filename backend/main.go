package main

import (
	"net/http"

	"github.com/ZakkBob/AskDave/gocommon/url"
"github.com/gin-gonic/gin"
)

var _, _ = url.ParseAbs("") //keep import

type album struct {
	ID   string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

var albums = []album{
	{ID: "1", Title: "W album", Artist: "John Coltrane", Price: 56.99},
	{ID: "2", Title: "Best Album", Artist: "Gerry Mulligan", Price: 17.99},
	{ID: "3", Title: "Even better album", Artist: "Sarah Vaughan", Price: 39.99},
}

func getAlbums(context *gin.Context) {
	context.IndentedJSON(http.StatusOK, albums)
}

func getAlbumByID(context *gin.Context) {
	id := context.Param("id")

	for _, album := range albums {
		if album.ID == id {
			context.IndentedJSON(http.StatusOK, album)
			return
		}
	}
	context.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
}

func postAlbums(context *gin.Context) {
	var newAlbum album

	if err := context.BindJSON(&newAlbum); err != nil {
		return
	}

	albums = append(albums, newAlbum)
	context.IndentedJSON(http.StatusCreated, newAlbum)
}

func deleteAlbum(context *gin.Context) {
	id := context.Param("id")

	for i, album := range albums {
		if album.ID == id {
			deletedAlbum := albums[i]
			albums[i] = albums[len(albums)-1]
			albums = albums[:len(albums)-1]
			context.IndentedJSON(http.StatusOK, deletedAlbum)
			return
		}
	}

	context.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
}

func main() {
	router := gin.Default()
	router.GET("/albums", getAlbums)
	router.GET("/albums/:id", getAlbumByID)
	router.POST("/albums", postAlbums)
	router.DELETE("/albums/:id", deleteAlbum)

	router.Run("localhost:8080")
}
