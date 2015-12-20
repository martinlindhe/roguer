package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/martinlindhe/ravel/views"
	"github.com/martinlindhe/rogue"
)

func main() {

	//seed := time.Now().Unix()
	seed := int64(1450549167)

	log.Printf("Generating island with seed %d ...\n", seed)
	island := rogue.GenerateIsland(seed, 220, 140)
	island.FillWithCritters()
	log.Println("Done generating island")

	//	islandColImage := island.ColoredHeightMapAsImage()

	/*
		islandColImgFile, _ := os.Create("island_col.png")
		png.Encode(islandColImgFile, islandColImage)

		islandImage := island.HeightMapAsImage()
		islandImgFile, _ := os.Create("island.png")
		png.Encode(islandImgFile, islandImage)
	*/
	/*
		for i := 0; i < 10; i++ {
			island.Tick()
		}
	*/

	r := getRouter()

	// r.GET("/", views.Index()) // XXX: cant get this form to work with gorazor views
	r.GET("/", func(c *gin.Context) {
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(200, views.Index())
	})

	// XXX run http server in separate process?

	// listen and serve on 0.0.0.0:8080
	appPort := 3322
	listenAt := fmt.Sprintf(":%d", appPort)

	log.Printf("Starting http server on %s\n", listenAt)

	r.Run(listenAt)
}

func getRouter() *gin.Engine {

	// Creates a router without any middleware by default
	r := gin.New()

	// Global middleware
	r.Use(gin.Logger())

	//	r.Use(gzip.Gzip(gzip.DefaultCompression))

	r.GET("/ping", pingController)

	r.Static("/js", "./public/js")
	r.Static("/css", "./public/css")
	r.Static("/img", "./public/img")
	r.Static("/fonts", "./public/fonts")
	r.Static("/flags", "./public/flags")
	//r.LoadHTMLFiles("./public/index.html")
	/*
		// non authenticated api endpoints:
		r.POST("/api/auth/register", apiAuthRegister)
		r.POST("/api/auth/login", apiAuthLogin)
		//Route::post('auth/register', 'Api\Auth\AuthController@postRegister');
		//Route::post('auth/login', 'Api\Auth\AuthController@postLogin');

		// authenticated api endpoints: XXX check jwt token
		r.GET("/api/auth/logout", apiAuthLogout)
		//Route::get('auth/refresh-token', ['middleware' => 'jwt.refresh', 'uses' => 'Api\Auth\AuthController@refreshToken']);
		//Route::get('auth/logout', 'Api\Auth\AuthController@getLogout');
	*/
	return r
}

// curl -v "http://localhost:8080/ping"
func pingController(c *gin.Context) {
	c.JSON(200, gin.H{"pong": "now"})
}
