package main

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Cheese struct {
	gorm.Model
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int    `json:"price"`
}

var db *gorm.DB

func cheese(c echo.Context) error {
	var cheese []Cheese

	if err := db.Find(&cheese).Error; err != nil {
		c.JSON(http.StatusInternalServerError, echo.Map{"error": "Internal Server Error(500)"})
	}

	return c.JSON(http.StatusOK, cheese)
}

func createCheese(c echo.Context) error {
	cheeseCreate := new(Cheese)

	if err := c.Bind(cheeseCreate); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Bad input"})
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
