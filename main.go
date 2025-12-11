package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
)

// Структура для объявления
type Announcement struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Author  string `json:"author"`
	Date    string `json:"date"`
}

// Глобальные переменные
var announcements []Announcement
var nextID = 1

// Загружаем данные
func loadData() error {
	fmt.Println("📂 Загружаю данные...")

	// Читаем файл
	data, err := os.ReadFile("bulletin.json")
	if err != nil {
		return err
	}

	// Ваш формат
	var dataStruct struct {
		Items []Announcement `json:"announcements"`
	}

	if err := json.Unmarshal(data, &dataStruct); err != nil {
		return err
	}

	announcements = dataStruct.Items

	// Добавляем ID если их нет
	for i := range announcements {
		if announcements[i].ID == 0 {
			announcements[i].ID = nextID
			nextID++
		}
	}

	fmt.Printf("✅ Загружено %d объявлений\n", len(announcements))
	return nil
}

// Сохраняем данные
func saveData() error {
	data := struct {
		Items []Announcement `json:"announcements"`
	}{
		Items: announcements,
	}

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile("bulletin.json", jsonData, 0644)
}

func main() {
	// Загружаем данные
	if err := loadData(); err != nil {
		fmt.Printf("❌ Ошибка загрузки: %v\n", err)
		fmt.Println("Создаю новый файл...")

		// Создаем пустой файл
		emptyData := struct {
			Items []Announcement `json:"announcements"`
		}{
			Items: []Announcement{},
		}

		jsonData, _ := json.MarshalIndent(emptyData, "", "  ")
		os.WriteFile("bulletin.json", jsonData, 0644)
		announcements = []Announcement{}
	}

	// Создаем приложение
	app := fiber.New()

	// Все объявления
	app.Get("/api/announcements", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"success": true,
			"data":    announcements,
			"count":   len(announcements),
		})
	})

	// Одно объявление
	app.Get("/api/announcements/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		for _, ann := range announcements {
			if fmt.Sprint(ann.ID) == id {
				return c.JSON(fiber.Map{
					"success": true,
					"data":    ann,
				})
			}
		}
		return c.Status(404).JSON(fiber.Map{
			"error":   true,
			"message": "Не найдено",
		})
	})

	// Создать
	app.Post("/api/announcements", func(c *fiber.Ctx) error {
		type Request struct {
			Title   string `json:"title"`
			Content string `json:"content"`
			Author  string `json:"author"`
		}

		var req Request
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": true, "message": "Ошибка JSON"})
		}

		newAnn := Announcement{
			ID:      nextID,
			Title:   req.Title,
			Content: req.Content,
			Author:  req.Author,
			Date:    time.Now().Format(time.RFC3339),
		}

		nextID++
		announcements = append(announcements, newAnn)
		saveData()

		return c.JSON(fiber.Map{
			"success": true,
			"data":    newAnn,
		})
	})

	// Обновить
	app.Put("/api/announcements/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")

		type Request struct {
			Title   string `json:"title"`
			Content string `json:"content"`
			Author  string `json:"author"`
		}

		var req Request
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": true, "message": "Ошибка JSON"})
		}

		for i, ann := range announcements {
			if fmt.Sprint(ann.ID) == id {
				if req.Title != "" {
					announcements[i].Title = req.Title
				}
				if req.Content != "" {
					announcements[i].Content = req.Content
				}
				if req.Author != "" {
					announcements[i].Author = req.Author
				}
				announcements[i].Date = time.Now().Format(time.RFC3339)

				saveData()
				return c.JSON(fiber.Map{
					"success": true,
					"data":    announcements[i],
				})
			}
		}

		return c.Status(404).JSON(fiber.Map{"error": true, "message": "Не найдено"})
	})

	// Удалить
	app.Delete("/api/announcements/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")

		for i, ann := range announcements {
			if fmt.Sprint(ann.ID) == id {
				announcements = append(announcements[:i], announcements[i+1:]...)
				saveData()
				return c.JSON(fiber.Map{
					"success": true,
					"message": "Удалено",
				})
			}
		}

		return c.Status(404).JSON(fiber.Map{"error": true, "message": "Не найдено"})
	})

	// Главная
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("🏫 API Доски объявлений работает!\n\nИспользуйте:\nGET /api/announcements - все\nPOST /api/announcements - создать\nGET /api/announcements/1 - одно\nPUT /api/announcements/1 - обновить\nDELETE /api/announcements/1 - удалить")
	})

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "OK",
			"service": "Class Bulletin API",
		})
	})

	// Запуск
	fmt.Println("🚀 Сервер запущен: http://localhost:8080")
	app.Listen(":8080")
}
