package handlers

import (
	"credibooktest/config"
	"credibooktest/models"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/spf13/viper"
	"github.com/thedevsaddam/govalidator"
	"golang.org/x/crypto/bcrypt"
)

func Login(c echo.Context) error {
	type LoginData struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var user models.User
	var count int
	var tokenData string

	db := config.App.DBConfig

	login := new(LoginData)
	c.Bind(login)

	db.Where("username = ?", login.Username).First(&user).Count(&count)
	if count == 1 {
		//validate password
		err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(login.Password))

		if err == nil {
			claims := &config.JwtCustomClaims{
				user.ID,
				user.Username,
				jwt.StandardClaims{
					ExpiresAt: time.Now().Add(time.Hour * 6).Unix(),
				},
				user.IsAdmin,
			}
			// Create token with claims
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

			// Generate token as response
			tokenData, err = token.SignedString([]byte(viper.GetString("jwtSign")))
			if err != nil {
				return err
			}
			c.Set("x-credential-identifier", user.Username)
			return c.JSON(http.StatusOK, tokenData)
		}
	}

	return c.JSON(http.StatusUnauthorized, "Invalid username or password")
}

func AddUser(c echo.Context) error {
	var user, existingUser models.User
	var count int
	db := config.App.DBConfig

	rules := govalidator.MapData{
		"username": []string{"required"},
		"password": []string{"required"},
	}

	opts := govalidator.Options{
		Request: c.Request(),
		Data:    &user,
		Rules:   rules,
	}

	v := govalidator.New(opts)
	e := v.ValidateJSON()

	if len(e) > 0 {
		err := map[string]interface{}{"Validation error": e}
		return c.JSON(http.StatusBadRequest, err)
	}

	db.Where("username = ?", user.Username).Find(&existingUser).Count(&count)
	if count > 0 {
		return c.JSON(http.StatusCreated, "Username has used, please try another username")
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 8)
	user.Password = string(password)

	db.Create(&user)

	return c.JSON(http.StatusCreated, "New user created")
}

func GetAllUsers(c echo.Context) error {
	db := config.App.DBConfig
	var usersResponses []models.UserResponse
	var pages, perpages int

	if page := c.FormValue("page"); page != "" {
		pages, _ = strconv.Atoi(page)
	}

	if perpage := c.FormValue("perpage"); perpage != "" {
		perpages, _ = strconv.Atoi(perpage)
	}

	db = db.Table("users").Find(&usersResponses)

	p := Paginator{
		DB:      db,
		Page:    pages,
		PerPage: perpages,
	}

	result := p.paginate(&usersResponses)

	return c.JSON(http.StatusCreated, result)
}
