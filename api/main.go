package main

import (
	"github.com/JohnKucharsky/golang-api/models"
	"github.com/JohnKucharsky/golang-api/storage"
	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
	"gorm.io/gorm"
	"log"
	"net/http"
)

type Repository struct {
	DB *gorm.DB
}

func (r *Repository) SetupRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Post("/books", r.CreateBook)
	api.Delete("/books/:id", r.DeleteBook)
	api.Get("/books/:id", r.GetBookById)
	api.Get("/books", r.GetBooks)
}

func (r *Repository) CreateBook(ctx *fiber.Ctx) error {
	book := models.Book{}

	err := ctx.BodyParser(&book)

	if err != nil {
		err := ctx.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": "request failed"},
		)
		if err != nil {
			return err
		}
	}

	err = r.DB.Create(&book).Error

	if err != nil {
		err := ctx.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not create the book"},
		)
		if err != nil {
			return err
		}
	}

	err = ctx.Status(http.StatusOK).JSON(
		&fiber.Map{"message": "book created"},
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetBooks(ctx *fiber.Ctx) error {
	bookModels := &[]models.Book{}

	err := r.DB.Find(bookModels).Error

	if err != nil {
		err := ctx.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not get books"},
		)
		if err != nil {
			return err
		}
	}

	err = ctx.Status(http.StatusOK).JSON(
		&fiber.Map{"data": bookModels},
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) DeleteBook(ctx *fiber.Ctx) error {
	bookModel := models.Book{}
	id := ctx.Params("id")

	if id == "" {
		err := ctx.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{
				"message": "id cannot be empty",
			},
		)
		if err != nil {
			return err
		}
		return nil
	}

	err := r.DB.Delete(bookModel, id)

	if err.Error != nil {
		ctx.Status(http.StatusBadRequest).JSON(
			&fiber.Map{
				"message": "could not delete book",
			},
		)
		return err.Error
	}

	ctx.Status(http.StatusOK).JSON(
		&fiber.Map{
			"message": "books deleted successfully",
		},
	)

	return nil
}

func (r *Repository) GetBookById(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	bookModel := models.Book{}

	if id == "" {
		ctx.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{
				"message": "id cannot be empty",
			},
		)
	}

	err := r.DB.Where("id = ?", id).First(bookModel).Error

	if err != nil {
		ctx.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not get the book"},
		)
		return err
	}

	ctx.Status(http.StatusOK).JSON(
		&fiber.Map{
			"data": bookModel,
		},
	)

	return nil
}

func main() {
	//start db
	log.Print("Prepare db...")
	db, err := storage.NewConnection()
	if err != nil {
		log.Fatal("Could not load the database")
	}
	//end start db

	//migrate
	err = storage.MigrateBooks(db)
	if err != nil {
		log.Fatal("Could not migrate db")
	}
	//end migrate

	//start app
	r := Repository{
		DB: db,
	}
	app := fiber.New()
	r.SetupRoutes(app)
	err = app.Listen(":8000")
	if err != nil {
		log.Fatal("Could not start the app")
	}
	//end start the app
}
