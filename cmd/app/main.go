package main

import (
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

func main() {
	// Загружаем страницу
	urlConfluence := "https://confluence.hflabs.ru/pages/viewpage.action?pageId=1181220999"
	doc, err := goquery.NewDocument(urlConfluence)
	if err != nil {
		fmt.Printf("Unable to retrieve page: %v\n", err)
		fmt.Scanln()
		return
	}

	// Ищем таблицу на странице
	table := doc.Find("table")
	if table.Length() == 0 {
		fmt.Println("Unable to find table on page")
		fmt.Scanln()
		return
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
	pathKey := "./cmd/keys/polar-office-377708-b66a7701095e.json"
	srv, err := sheets.NewService(context.Background(), option.WithCredentialsFile(pathKey))
	if err != nil {
		fmt.Printf("Unable to retrieve Sheets client: %v\n", err)
		fmt.Scanln()
		return
	}

	// Определяем ID таблицы
	spreadsheetID := "1N34uWDj33UlkCazW2v-SFtJy5NE1HiS5TxWk8UwDang"
	urlGoogleShets := "https://docs.google.com/spreadsheets/d/1N34uWDj33UlkCazW2v-SFtJy5NE1HiS5TxWk8UwDang/edit#gid=0"
	// Записываем данные в таблицу
	writeOption := "RAW"
	_, err = srv.Spreadsheets.Values.Update(spreadsheetID, "A1", &sheets.ValueRange{Values: data}).ValueInputOption(writeOption).Do()
	if err != nil {
		fmt.Printf("Unable to write data to sheet: %v\n", err)
		fmt.Scanln()
		return
	}

	fmt.Printf("Перенос c cайта %s\n в таблицу %s\n выполнен успешно!\nНажмите Enter для выхода из программы.", urlConfluence, urlGoogleShets)
	fmt.Scanln()

}
