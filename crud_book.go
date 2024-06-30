package main

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// วิธีการส่งออกไปเป็น json
func getBooks(c *fiber.Ctx) error {
	return c.JSON(books)
}

func getBook(c *fiber.Ctx) error {
	bookId, err := strconv.Atoi(c.Params("id"))

	if err != nil {
		return c.Status(fiber.ErrBadRequest.Code).SendString(err.Error()) // (fiber.ErrBadRequest) ErrBadRequest คือ 400
	}
	for _, book := range books {
		if book.Id == bookId {
			return c.JSON(book)
		}
	}
	return c.Status(fiber.ErrNotFound.Code).SendString("Not found")
}

func createBook(c *fiber.Ctx) error {
	book := new(Book)                          //จองพื้นที่ไว้ ตัวนี้มีค่าเป้น *book
	if err := c.BodyParser(book); err != nil { // ต้องมีการเขียน c.BodyParser(data) อยู่ในโค้ด err:= c.BodyParser(data)แบบนี้ตัว BodyParser ก็ทำงานเหมือนกัน
		return c.Status(fiber.ErrBadRequest.Code).SendString(err.Error())
	}
	books = append(books, *book)
	return c.JSON(book)
}

func updateBook(c *fiber.Ctx) error {
	bookId, err := strconv.Atoi(c.Params("id"))

	// handle error เวลาที่ค่า param ที่รับมาจาก url ไม่เป็น int เช่น api/books/a "ตรง a"
	if err != nil {
		return c.Status(fiber.ErrBadRequest.Code).SendString(err.Error())
	}
	// ค่าที่ส่งมาที่ต้องการเปลี่ยน
	bookUpdate := new(Book)
	if err := c.BodyParser(bookUpdate); err != nil { // line นี้ต้องมีเวลาที่มีการส่งแบบ user ส่งข้อมูลมา
		return c.Status(fiber.ErrBadRequest.Code).SendString(err.Error())
	}
	for index, bookData := range books {
		if bookData.Id == bookId {
			// books[index] = *book ไม่ควรทำแบบนี้เพราะจะกลายเป็นว่าเปลี่ยน id ไปด้วย
			books[index].Title = bookUpdate.Title
			books[index].Author = bookUpdate.Author
			return c.JSON(books[index])
		}
	}
	return c.Status(fiber.ErrNotFound.Code).SendString("Not found")
}

func deleteBook(c *fiber.Ctx) error {
	bookId, err := strconv.Atoi(c.Params("id"))

	// handle error เวลาที่ค่า param ที่รับมาจาก url ไม่เป็น int เช่น api/books/a "ตรง a"
	if err != nil {
		return c.Status(fiber.ErrBadRequest.Code).SendString(err.Error())
	}

	for index, bookData := range books {
		if bookData.Id == bookId {
			// ทำให้เห็นภาพ append บรรทัดที่ 105
			// [1,2,3,4,5] ตัวอย่าง index = 3
			// books[:2] => [1,2]
			// books[2+1:] => [4,5]
			// [1,2] + [4,5]
			books = append(books[:index], books[index+1:]...)
			return c.SendStatus(fiber.StatusNoContent)
		}
	}
	return c.Status(fiber.ErrNotFound.Code).SendString("Not found")
}
