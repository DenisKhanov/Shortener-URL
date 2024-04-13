package handlers // Package main is the main package for the application.

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
)

// ExampleGetShortURL demonstrates the usage of GetShortURL in Godoc.
// This example shows how to make a POST request to shorten a URL.
// Output: Status Code: 201
//
//		Response Body: <shortened_url>
//		Expected Result: A successful response with status code 201 and the shortened URL
//	    or HTTP status 409 Conflict if the URL is already shortened and the shortened later URL
func ExampleHandlers_GetShortURL() {
	// Create a test server
	router := gin.Default()

	myHandler := Handlers{}

	publicRoutes := router.Group("/")
	publicRoutes.Use(myHandler.MiddlewareAuthPublic())
	publicRoutes.Use(myHandler.MiddlewareLogging())
	publicRoutes.Use(myHandler.MiddlewareCompress())
	publicRoutes.POST("/", myHandler.GetShortURL)

	// Prepare a test request with valid data
	requestBody := bytes.NewBufferString("https://example.com/long-url")
	req, _ := http.NewRequest("POST", "http:/localhost:8080/", requestBody)
	req.Header.Set("Content-Type", "text/plain")

	// Simulate the request
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Print the result
	fmt.Println("Status Code:", w.Code)
	fmt.Println("Response Body:", w.Body.String())
}
