package main

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type CreateUserDTO struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func main() {
	fmt.Println("Hello, Gorm")
	dsn := "host=localhost user=postgres password=postgrespw dbname=learn_gorm port=55000 sslmode=disable TimeZone=Asia/Bangkok"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		fmt.Println(err)
	}

	e := echo.New()
	e.Use(middleware.Recover())
	e.GET("/users", func(c echo.Context) error {
		ctx := c.Request().Context()
		users := make([]User, 0)
		db.WithContext(ctx).Find(&users)
		return c.JSON(200, &users)
	})

	e.POST("/migrate", func(c echo.Context) error {
		fmt.Println("migrating database...")
		ctx := c.Request().Context()
		err := db.WithContext(ctx).AutoMigrate(&User{}, &Book{})
		if err != nil {
			fmt.Println(err)
			return c.JSON(400, map[string]string{
				"error": err.Error(),
			})
		}
		return c.JSON(200, map[string]string{
			"error": "",
			"msg":   "database migrated",
		})

	})
	e.POST("/user", func(c echo.Context) error {
		fmt.Println("creating user...")
		ctx := c.Request().Context()
		createUserDto := new(CreateUserDTO)
		if err := c.Bind(createUserDto); err != nil {
			e.Logger.Error(err)
			return c.JSON(400, err)
		}

		hashed, _ := hashPassword(createUserDto.Password)
		user := &User{
			Name:     createUserDto.Name,
			Username: createUserDto.Username,
			Password: hashed,
		}
		db.WithContext(ctx).Create(user)
		return c.JSON(201, user)
	})

	e.Logger.Fatal(e.Start(":3300"))
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
