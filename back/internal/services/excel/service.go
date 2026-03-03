package excel

import (
	"fmt"
	"keden-service/back/internal/models"

	"github.com/xuri/excelize/v2"
)

type ExcelService struct{}

func NewExcelService() *ExcelService {
	return &ExcelService{}
}

// GenerateFromData builds an Excel file from typed DocumentFields and DocumentItem rows.
// Column layout matches the example.xlsx template exactly.
func (s *ExcelService) GenerateFromData(fields *models.DocumentFields, items []models.DocumentItem) ([]byte, error) {
	f := excelize.NewFile()
	defer f.Close()

	sheet := "Таможенная декларация"
	index, err := f.NewSheet(sheet)
	if err != nil {
		return nil, err
	}
	f.SetActiveSheet(index)
	f.DeleteSheet("Sheet1")

	// Styles
	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 12, Color: "FFFFFF"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"2F5496"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
		Border:    allBorders(),
	})
	cellStyle, _ := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Vertical: "center", WrapText: true},
		Border:    allBorders(),
	})
	titleStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 14},
		Alignment: &excelize.Alignment{Horizontal: "center"},
	})
	labelStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
	})

	// Title
	f.MergeCell(sheet, "A1", "J1")
	f.SetCellValue(sheet, "A1", "Результаты обработки таможенной декларации")
	f.SetCellStyle(sheet, "A1", "J1", titleStyle)

	// Header fields — two columns: label | value
	type row struct {
		label string
		value interface{}
	}
	headerRows := []row{
		{"Номер декларации", fields.DeclarationNumber},
		{"Дата", fields.Date},
		{"Отправитель", fields.Sender},
		{"Получатель", fields.Receiver},
		{"Страна происхождения", fields.CountryOrigin},
		{"Страна назначения", fields.CountryDest},
		{"Валюта", fields.Currency},
		{"Общая стоимость", fields.TotalValue},
		{"Таможенная стоимость", fields.CustomsValue},
	}

	r := 3
	for _, hr := range headerRows {
		f.SetCellValue(sheet, fmt.Sprintf("A%d", r), hr.label)
		f.SetCellStyle(sheet, fmt.Sprintf("A%d", r), fmt.Sprintf("A%d", r), labelStyle)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", r), hr.value)
		r++
	}

	// Items table
	if len(items) > 0 {
		r += 2
		f.MergeCell(sheet, fmt.Sprintf("A%d", r), fmt.Sprintf("J%d", r))
		f.SetCellValue(sheet, fmt.Sprintf("A%d", r), "Товары")
		f.SetCellStyle(sheet, fmt.Sprintf("A%d", r), fmt.Sprintf("J%d", r), titleStyle)
		r++

		headers := []string{"№", "Код ТН ВЭД", "Описание", "Кол-во", "Ед.изм.", "Вес нетто (кг)", "Вес брутто (кг)", "Стоимость", "Пошлина", "НДС"}
		for i, h := range headers {
			cell := fmt.Sprintf("%s%d", col(i), r)
			f.SetCellValue(sheet, cell, h)
			f.SetCellStyle(sheet, cell, cell, headerStyle)
		}
		r++

		for _, item := range items {
			vals := []interface{}{
				item.Number,
				item.HSCode,
				item.Description,
				item.Quantity,
				item.Unit,
				item.WeightNet,
				item.WeightGross,
				item.Value,
				item.DutyRate,
				item.VATRate,
			}
			for i, v := range vals {
				cell := fmt.Sprintf("%s%d", col(i), r)
				f.SetCellValue(sheet, cell, v)
				f.SetCellStyle(sheet, cell, cell, cellStyle)
			}
			r++
		}
	}

	f.SetColWidth(sheet, "A", "A", 25)
	f.SetColWidth(sheet, "B", "B", 20)
	f.SetColWidth(sheet, "C", "C", 35)
	f.SetColWidth(sheet, "D", "J", 15)

	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func col(i int) string {
	return string(rune('A' + i))
}

func allBorders() []excelize.Border {
	return []excelize.Border{
		{Type: "left", Color: "000000", Style: 1},
		{Type: "right", Color: "000000", Style: 1},
		{Type: "top", Color: "000000", Style: 1},
		{Type: "bottom", Color: "000000", Style: 1},
	}
}
