package main

import (
	"fmt"
	"log"
	"os"
	"time"

	// "github.com/golang-jwt/jwt"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/gofiber/template/html/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
)

// ด้านหลังที่เป็น json คือกำหนดชื่อข้อมูลเวลาส่งออกไป เป็น json ตัวอย่างถ้าเป็น `่json:hello` เวลาส่งออกจะได้ "hello" : data
type Book struct {
	Id     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

var books []Book

func checkMiddleware(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	if claims["role"] != "admin" {
		return fiber.ErrUnauthorized
	}
	return c.Next()
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("load .env error")
	}
	engine := html.New("./views", ".html")
	app := fiber.New(fiber.Config{Views: engine})

	books = append(books, Book{Id: 1, Title: "Harry Potter", Author: "J.K"})
	books = append(books, Book{Id: 2, Title: "Jujutsu Kaisen ", Author: "Idk"})
	fmt.Println(books)
	app.Post("/login", loginUser)

	// JWT Middleware   method ที่อยู่หลังจาก line:57 จำเป็นต้องตรวจสอบผ่าน middleware ก่อนถึงจะใช้งานได้
	// โดยที่ตัว Secret-key เราจะกำหนดที่ SigningKey ถ้าจะ authorization ได้ก็ต้องมีการสร้าง Jwt โดยนำ Secret-key ตัวเดียวกันไปสร้าง
	app.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte(os.Getenv("SECRET")),
	})) // middleware ตรงนี้ตรวจว่า login ถูกไหม
	//ต้องมีการผ่าน middleware jwt ก่อนถึงจะไป middleware อันต่อไป
	app.Use(checkMiddleware) //middleware ตรงนี้จะตรวจว่า role ถูกไหม ถ้าถูกถึงจะอนุญาติส่งไป method ต่อไปได้

	app.Get("/books", getBooks) //Books โดยไม่ใส่ () จะเป็นการเอาทั้ง func ไปใส่
	app.Get("/books/:id", getBook)
	app.Post("/books", createBook)
	app.Put("/books/:id", updateBook)
	app.Delete("/books/:id", deleteBook)
	app.Post("/upload", uploadImage)
	app.Get("html-views", viewsHTML)
	app.Get("/config", getEnv)

	app.Listen(":8080")
}

func uploadImage(c *fiber.Ctx) error {
	file, err := c.FormFile("image")
	if err != nil {
		return c.Status(fiber.ErrBadRequest.Code).SendString(err.Error())
	}
	err = c.SaveFile(file, "./uploads/"+file.Filename)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.SendString("File upload complete!")
}

func viewsHTML(c *fiber.Ctx) error {
	return c.Render("index", fiber.Map{
		"Title": books[0].Title,
	})
}

func getEnv(c *fiber.Ctx) error {

	return c.JSON(fiber.Map{
		"SECRET": os.Getenv("SECRET"),
	})
}

var userTest = User{
	Email:    "user@test.com",
	Password: "123456",
}

// JWT Token check
func loginUser(c *fiber.Ctx) error {
	user := new(User)
	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	if user.Email != userTest.Email || user.Password != userTest.Password {
		return fiber.ErrUnauthorized // return กลับไปว่ายืนยันตัวตนไม่ถูกต้อง
	}

	// Create the Claims กำหนดว่าจะเอาอะไรไปสร้างเป็น Pattern บ้าง
	claims := jwt.MapClaims{
		"email": user.Email,
		"role":  "admin",
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	}

	// Create token สร้าง pattern
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{
		"message": "login success",
		"status":  "ok",
		"token":   t,
	})
}
