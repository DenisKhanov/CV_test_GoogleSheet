package main

import (
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
	"log"
)

func main() {
	// Загружаем страницу
	doc, err := goquery.NewDocument("https://confluence.hflabs.ru/pages/viewpage.action?pageId=1181220999")
	if err != nil {
		log.Fatalf("Unable to retrieve page: %v", err)
	}

	// Ищем таблицу на странице
	table := doc.Find("table")
	if table.Length() == 0 {
		log.Fatalf("Unable to find table on page")
	}

	// Создаем массив для хранения данных
	var data [][]interface{}

	// Выбрать заголовки таблицы
	var headers []interface{}
	table.Find("thead th").Each(func(i int, th *goquery.Selection) {
		headers = append(headers, th.Text())
	})
	data = append(data, headers)

	// Итерируемся по строкам таблицы
	table.Find("tr").Each(func(i int, tr *goquery.Selection) {
		// Итерируемся по ячейкам строки
		var row []interface{}
		tr.Find("td").Each(func(j int, td *goquery.Selection) {
			// Добавляем текст ячейки в массив строки
			row = append(row, td.Text())
		})
		// Добавляем строку в массив данных
		data = append(data, row)
	})

	// Инициализируем клиента Google Sheets API
	srv, err := sheets.NewService(context.Background(), option.WithCredentialsFile("polar-office-377708-b66a7701095e.json"))
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	// Определяем ID таблицы
	spreadsheetID := "11MxexsZ-8UZu7hJEUMZYU0itFp7oCRp9n2ZeNax4Eo8"

	// Записываем данные в таблицу
	writeOption := "RAW"
	resp, err := srv.Spreadsheets.Values.Update(spreadsheetID, "A1", &sheets.ValueRange{Values: data}).ValueInputOption(writeOption).Do()
	if err != nil {
		log.Fatalf("Unable to write data to sheet: %v", err)
	}

	fmt.Printf("%#v\n", resp)

}
