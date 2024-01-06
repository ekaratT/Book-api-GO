package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	e := echo.New()

	b := BookHandler{}
	b.Initialize()

	e.GET("/books", b.GetAllBook)
	e.POST("/books", b.SaveBook)
	e.GET("/books/:id", b.GetBookById)
	e.PUT("/books/:id", b.UpdateBook)
	e.DELETE("/books/:id", b.DeleteBook)
	e.Logger.Fatal(e.Start(":8080"))
}

type Book struct {
	Id     uint64  `gorm:"primary_key; autoIncrement:true" json:"id"`
	Title  string  `json:"title"`
	Author string  `json:"author"`
	Price  float64 `json:"price"`
	Pages  int     `json:"pages"`
}
type BookHandler struct {
	DB *gorm.DB
}

func (b *BookHandler) Initialize() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(".env file couldn't be loaded.")
	}
	dsn := os.Getenv("DSN")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	db.AutoMigrate(&Book{})
	b.DB = db
}

func (b *BookHandler) GetAllBook(c echo.Context) error {
	books := []Book{}

	b.DB.Find(&books)

	return c.JSON(http.StatusOK, books)
}

func (b *BookHandler) GetBookById(c echo.Context) error {
	id := c.Param("id")
	book := Book{}

	if err := b.DB.Find(&book, id).Error; err != nil {
		return c.NoContent(http.StatusNotFound)
	}

	return c.JSON(http.StatusOK, book)
}

func (b *BookHandler) SaveBook(c echo.Context) error {
	book := Book{}

	if err := c.Bind(&book); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	if err := b.DB.Save(&book).Error; err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, book)
}

func (b *BookHandler) UpdateBook(c echo.Context) error {
	id := c.Param("id")
	book := Book{}

	if err := b.DB.Find(&book, id).Error; err != nil {
		return c.NoContent(http.StatusNotFound)
	}

	if err := c.Bind(&book); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	if err := b.DB.Save(&book).Error; err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, book)
}

func (b *BookHandler) DeleteBook(c echo.Context) error {
	id := c.Param("id")
	book := Book{}

	if err := b.DB.Find(&book, id).Error; err != nil {
		return c.NoContent(http.StatusNotFound)
	}

	if err := b.DB.Delete(&book, id).Error; err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusNoContent)
}
