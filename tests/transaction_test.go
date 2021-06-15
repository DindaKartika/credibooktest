package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"credibooktest/router"

	"github.com/gavv/httpexpect"
)

type Login struct {
	Username string `form:"username"`
	Password string `form:"password"`
}

func TestGetTransaction(t *testing.T) {
	// create http.Handler
	handler := router.New()

	// run server using httptest
	server := httptest.NewServer(handler)

	defer server.Close()

	// create httpexpect instance
	e := httpexpect.New(t, server.URL)

	r := e.POST("/login").WithForm(Login{"admin", "admin"}).
		Expect().
		Status(http.StatusOK).JSON().String()

	token := r.Raw()

	// without query string
	obj := e.GET("/transaction").WithHeader("Authorization", "Bearer "+token).
		Expect().
		Status(http.StatusOK).JSON().Object()
	obj.Keys().Contains("data")
	obj.Value("data").Array().Element(0).Object()

	// with page set
	obj = e.GET("/transaction").WithQuery("page", "1").WithHeader("Authorization", "Bearer "+token).
		Expect().
		Status(http.StatusOK).JSON().Object()
	obj.Keys().Contains("data")

	// with perpage set
	obj = e.GET("/transaction").WithQuery("perpage", "10").WithHeader("Authorization", "Bearer "+token).
		Expect().
		Status(http.StatusOK).JSON().Object()
	obj.Keys().Contains("data")

	// with filter type
	obj = e.GET("/transaction").WithQuery("type", "expense").WithHeader("Authorization", "Bearer "+token).
		Expect().
		Status(http.StatusOK).JSON().Object()
	obj.Keys().Contains("data")
	obj.Value("data").Array().Element(0).Object().Value("type").String().Contains("expense")

	// with filter min_amount
	obj = e.GET("/transaction").WithQuery("min_amount", 5000).WithHeader("Authorization", "Bearer "+token).
		Expect().
		Status(http.StatusOK).JSON().Object()
	obj.Keys().Contains("data")

	// with filter max_amount
	obj = e.GET("/transaction").WithQuery("max_amount", 1000000).WithHeader("Authorization", "Bearer "+token).
		Expect().
		Status(http.StatusOK).JSON().Object()
	obj.Keys().Contains("data")
}

func TestAddTransaction(t *testing.T) {
	// create http.Handler
	handler := router.New()

	// run server using httptest
	server := httptest.NewServer(handler)

	defer server.Close()

	// create httpexpect instance
	e := httpexpect.New(t, server.URL)

	r := e.POST("/login").WithForm(Login{"admin", "admin"}).
		Expect().
		Status(http.StatusOK).JSON().String()

	token := r.Raw()

	payload := make(map[string]interface{})

	// normal add new
	payload = map[string]interface{}{
		"amount": 150000,
		"notes":  "Penjualan",
		"type":   "income",
	}
	obj := e.POST("/transaction").
		WithJSON(payload).WithHeader("Authorization", "Bearer "+token).
		Expect().
		Status(http.StatusCreated).JSON().Object()
	obj.ContainsKey("amount").ValueEqual("amount", 150000)
	obj.ContainsKey("notes").ValueEqual("notes", "Penjualan")
	obj.ContainsKey("type").ValueEqual("type", "income")

	// failed because amount empty
	payload = map[string]interface{}{
		"notes": "Penjualan",
		"type":  "income",
	}
	obj = e.POST("/transaction").
		WithJSON(payload).WithHeader("Authorization", "Bearer "+token).
		Expect().
		Status(http.StatusBadRequest).JSON().Object()
	obj.Value("Validation error").Object().Value("amount").Array().Element(0).String().Equal("The amount field is required")
}

func TestUpdateTransaction(t *testing.T) {
	//create http handler
	handler := router.New()

	//run server using httptest
	server := httptest.NewServer(handler)
	defer server.Close()

	//create httpexpect instance
	e := httpexpect.New(t, server.URL)

	r := e.POST("/login").WithForm(Login{"admin", "admin"}).
		Expect().
		Status(http.StatusOK).JSON().String()

	token := r.Raw()

	inputData := make(map[string]interface{})

	//success update transaction
	inputData = map[string]interface{}{
		"amount": 130000,
	}
	obj := e.PUT("/transaction/{id}").WithPath("id", "1").WithJSON(inputData).WithHeader("Authorization", "Bearer "+token).Expect().Status(http.StatusOK).JSON().Object()
	obj.ContainsKey("id").ValueEqual("id", 1)
	obj.ContainsKey("transaction").ValueEqual("amount", 130000)

	//error update transaction
	e.PUT("/transaction/{id}").WithPath("id", "dua").WithJSON(inputData).WithHeader("Authorization", "Bearer "+token).Expect().Status(http.StatusNotFound)

	//error update transaction
	e.PUT("/transaction/{id}").WithPath("id", "1000000").WithJSON(inputData).WithHeader("Authorization", "Bearer "+token).Expect().Status(http.StatusNotFound)
}

func TestDeleteTransaction(t *testing.T) {
	//create http handler
	handler := router.New()

	//run server using httptest
	server := httptest.NewServer(handler)
	defer server.Close()

	//create httpexpect instance
	e := httpexpect.New(t, server.URL)

	r := e.POST("/login").WithForm(Login{"admin", "admin"}).
		Expect().
		Status(http.StatusOK).JSON().String()

	token := r.Raw()

	//success delete transaction
	e.DELETE("/transaction/{id}").WithPath("id", "1").WithHeader("Authorization", "Bearer "+token).Expect().Status(http.StatusOK)

	//error delete transaction
	e.DELETE("/transaction/{id}").WithPath("id", "dua").WithHeader("Authorization", "Bearer "+token).Expect().Status(http.StatusNotFound)

	//error delete transaction
	e.DELETE("/transaction/{id}").WithPath("id", "1000000").WithHeader("Authorization", "Bearer "+token).Expect().Status(http.StatusNotFound)
}
