package main

import (
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Cheese struct {
	gorm.Model
	ID          int    `json:"id"`
	Name        string `json:"name" validate:"required,min=3,max=30"`
	Description string `json:"description" validate:"required,min=3,max=350"`
	Price       int    `json:"price" validate:"required,gt=0"`
}

type cheeseGet struct {
	Name        string `json:"name" validate:"required,min=3,max=30"`
	Description string `json:"description" validate:"required,min=3,max=350"`
	Price       int    `json:"price" validate:"required,gt=0"`
}

func toCheeseGet(c Cheese) cheeseGet {
	return cheeseGet{
		Name:        c.Name,
		Description: c.Description,
		Price:       c.Price,
	}
}

var validate = validator.New()

var db *gorm.DB

func cheese(c echo.Context) error {
	var cheese []Cheese

	if err := db.Find(&cheese).Error; err != nil {
		c.JSON(http.StatusInternalServerError, echo.Map{"error": "Internal Server Error(500)"})
	}

	var resp []cheeseGet

	for _, ch := range cheese {
		resp = append(resp, toCheeseGet(ch))
	}

	return c.JSON(http.StatusOK, resp)
}

func createCheese(c echo.Context) error {
	cheeseCreate := new(Cheese)

	if err := c.Bind(cheeseCreate); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Bad input"})
	}

	if err := validate.Struct(cheeseCreate); err != nil {
		return c.JSON(422, echo.Map{"validation errror": err.Error()})
	}

	if err := db.Create(cheeseCreate).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Error making cheese"})
	}

	return c.JSON(http.StatusOK, cheeseCreate)
}

func deleteCheese(c echo.Context) error {
	var cheese Cheese
	id := c.Param("id")

	if err := db.First(&cheese, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "Cheese not found"})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Database error"})
	}

	// Удаление
	if err := db.Delete(&cheese).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to delete cheese"})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "Cheese deleted", "id": id})
}

func updateCheese(c echo.Context) error {
	id := c.Param("id")
	var cheese Cheese

	if err := db.First(&cheese, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "Cheese not found"})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Database error"})
	}

	var updateData map[string]interface{}
	if err := c.Bind(&updateData); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid input"})
	}

	if err := db.Model(&cheese).Updates(updateData).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to update cheese"})
	}

	return c.JSON(http.StatusOK, cheese)
}

func main() {
	var err error

	db, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("DB init fail^ %v", err)
	}

	db.AutoMigrate(&Cheese{})

	app := echo.New()
	app.GET("/cheese", cheese)
	app.POST("/cheeseCreate", createCheese)
	app.DELETE("/cheeseDel", deleteCheese)
	app.PATCH("cheesePatch", updateCheese)

	app.Logger.Fatal(app.Start(":1488"))
}
