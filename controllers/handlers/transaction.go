package handlers

import (
	"credibooktest/config"
	"credibooktest/models"
	"net/http"
	"strconv"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/thedevsaddam/govalidator"
)

func GetAllTransaction(c echo.Context) error {
	db := config.App.DBConfig
	var transactions []models.Transaction
	var orderBy []string
	var pages, perpages int

	user := claimUser(c)

	if !user.IsAdmin {
		db = db.Where("user_id = ?", user.ID)
	}

	if types := c.FormValue("type"); types != "" {
		db = db.Where("type = ?", types)
	}

	if minAmount := c.FormValue("min_amount"); minAmount != "" {
		db = db.Where("amount >= ?", minAmount)
	}

	if maxAmount := c.FormValue("max_amount"); maxAmount != "" {
		db = db.Where("amount =< ?", maxAmount)
	}

	if page := c.FormValue("page"); page != "" {
		pages, _ = strconv.Atoi(page)
	}

	if perpage := c.FormValue("perpage"); perpage != "" {
		perpages, _ = strconv.Atoi(perpage)
	}

	if order := c.FormValue("order_by"); order != "" {
		orderBy = strings.Split(order, ",")
		for i, orderData := range orderBy {
			if strings.Contains(orderData, "date") {
				orderBy[i] = strings.Replace(orderData, "date", "created_at", 1)
			}
		}
	}

	db = db.Find(&transactions)

	p := Paginator{
		DB:      db,
		OrderBy: orderBy,
		Page:    pages,
		PerPage: perpages,
	}

	result := p.paginate(&transactions)

	return c.JSON(http.StatusOK, result)
}

func AddTransaction(c echo.Context) error {
	db := config.App.DBConfig
	var transaction models.Transaction

	rules := govalidator.MapData{
		"amount": []string{"required"},
		"notes":  []string{},
		"type":   []string{},
	}

	opts := govalidator.Options{
		Request: c.Request(),
		Data:    &transaction,
		Rules:   rules,
	}

	v := govalidator.New(opts)
	e := v.ValidateJSON()

	if len(e) > 0 {
		err := map[string]interface{}{"Validation error": e}
		return c.JSON(http.StatusBadRequest, err)
	}

	if transaction.Type != "income" {
		transaction.Type = "expense"
	}

	user := claimUser(c)

	transaction.UserID = user.ID

	db.Create(&transaction)

	return c.JSON(http.StatusCreated, transaction)
}

func UpdateTransaction(c echo.Context) error {
	db := config.App.DBConfig
	var transaction, existTransaction models.Transaction

	id := c.Param("id")

	rules := govalidator.MapData{
		"amount": []string{"required"},
		"notes":  []string{},
		"type":   []string{},
	}

	opts := govalidator.Options{
		Request: c.Request(),
		Data:    &transaction,
		Rules:   rules,
	}

	v := govalidator.New(opts)
	e := v.ValidateJSON()

	if len(e) > 0 {
		err := map[string]interface{}{"Validation error": e}
		return c.JSON(http.StatusBadRequest, err)
	}

	if transaction.Type != "income" {
		transaction.Type = "expense"
	}

	user := claimUser(c)

	if db.Where("id = ? and user_id = ?", id, user.ID).Find(&existTransaction).RecordNotFound() {
		return c.JSON(http.StatusNotFound, "transaction not found")
	}

	db.Model(&existTransaction).Updates(map[string]interface{}{
		"amount": transaction.Amount,
		"type":   transaction.Type,
		"notes":  transaction.Notes,
	})

	return c.JSON(http.StatusOK, existTransaction)
}

func DeleteTransaction(c echo.Context) error {
	db := config.App.DBConfig
	var transaction models.Transaction

	id := c.Param("id")

	user := claimUser(c)

	if db.Where("id = ? and user_id = ?", id, user.ID).Find(&transaction).RecordNotFound() {
		return c.JSON(http.StatusNotFound, "transaction not found")
	}

	db.Delete(&transaction)

	return c.JSON(http.StatusOK, "transaction has deleted successfully")
}

func claimUser(c echo.Context) *config.JwtCustomClaims {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*config.JwtCustomClaims)
	return claims
}
