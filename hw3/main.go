package main

import (
	"fmt"
)

// Интерфейс Formatter с методом Format
type Formatter interface {
	Format(text string) string
}

// Структуры, удовлетворяющие интерфейсу Formatter

// Обычный текст (как есть)
type PlainText struct{}

func (p PlainText) Format(text string) string {
	return text
}

// Жирный шрифт
type Bold struct{}

func (b Bold) Format(text string) string {
	return "**" + text + "**"
}

// Код
type Code struct{}

func (c Code) Format(text string) string {
	return "`" + text + "`"
}

// Курсив
type Italic struct{}

func (i Italic) Format(text string) string {
	return "_" + text + "_"
}

// ChainFormatter для цепочки модификаторов
type ChainFormatter struct {
	formatters []Formatter
}

func (cf *ChainFormatter) AddFormatter(formatter Formatter) {
	cf.formatters = append(cf.formatters, formatter)
}

func (cf ChainFormatter) Format(text string) string {
	for _, formatter := range cf.formatters {
		text = formatter.Format(text)
	}
	return text
}

func main() {
	// Примеры использования
	plainText := PlainText{}
	bold := Bold{}
	code := Code{}
	italic := Italic{}

	// Обычный текст
	fmt.Println("Обычный текст:", plainText.Format("Hello, World!"))

	// Жирный шрифт
	fmt.Println("Жирный шрифт:", bold.Format("Hello, World!"))

	// Код
	fmt.Println("Код:", code.Format("Hello, World!"))

	// Курсив
	fmt.Println("Курсив:", italic.Format("Hello, World!"))

	// Цепочка модификаторов
	chainFormatter := &ChainFormatter{}
	chainFormatter.AddFormatter(code)
	chainFormatter.AddFormatter(bold)
	chainFormatter.AddFormatter(italic)

	// Форматируем строку с использованием цепочки модификаторов
	fmt.Println("Цепочка модификаторов:", chainFormatter.Format("Hello, World!"))
}
