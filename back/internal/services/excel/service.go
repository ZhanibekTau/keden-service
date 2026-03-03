package excel

import (
	"fmt"
	"keden-service/back/internal/services/ai"

	"github.com/xuri/excelize/v2"
)

type ExcelService struct{}

func NewExcelService() *ExcelService {
	return &ExcelService{}
}

func (s *ExcelService) GenerateFromAIResponse(resp *ai.AIResponse) ([]byte, error) {
	f := excelize.NewFile()
	defer f.Close()

	sheetName := "Таможенная декларация"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return nil, err
	}
	f.SetActiveSheet(index)
	f.DeleteSheet("Sheet1")

	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 12, Color: "FFFFFF"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"2F5496"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})

	cellStyle, _ := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Vertical: "center", WrapText: true},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})

	titleStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 14},
		Alignment: &excelize.Alignment{Horizontal: "center"},
	})

	f.MergeCell(sheetName, "A1", "J1")
	f.SetCellValue(sheetName, "A1", "Результаты обработки таможенной декларации")
	f.SetCellStyle(sheetName, "A1", "J1", titleStyle)

	row := 3
	fieldLabelStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
	})

	fieldNames := map[string]string{
		"declaration_number": "Номер декларации",
		"date":               "Дата",
		"sender":             "Отправитель",
		"receiver":           "Получатель",
		"country_origin":     "Страна происхождения",
		"country_dest":       "Страна назначения",
		"currency":           "Валюта",
		"total_value":        "Общая стоимость",
		"customs_value":      "Таможенная стоимость",
	}

	fieldOrder := []string{"declaration_number", "date", "sender", "receiver", "country_origin", "country_dest", "currency", "total_value", "customs_value"}

	for _, key := range fieldOrder {
		label, ok := fieldNames[key]
		if !ok {
			label = key
		}
		val, exists := resp.Fields[key]
		if !exists {
			continue
		}
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), label)
		f.SetCellStyle(sheetName, fmt.Sprintf("A%d", row), fmt.Sprintf("A%d", row), fieldLabelStyle)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), fmt.Sprintf("%v", val))
		row++
	}

	for key, val := range resp.Fields {
		if _, handled := fieldNames[key]; handled {
			continue
		}
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), key)
		f.SetCellStyle(sheetName, fmt.Sprintf("A%d", row), fmt.Sprintf("A%d", row), fieldLabelStyle)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), fmt.Sprintf("%v", val))
		row++
	}

	row += 2
	if len(resp.Items) > 0 {
		f.MergeCell(sheetName, fmt.Sprintf("A%d", row), fmt.Sprintf("J%d", row))
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), "Товары")
		f.SetCellStyle(sheetName, fmt.Sprintf("A%d", row), fmt.Sprintf("J%d", row), titleStyle)
		row++

		headers := []string{"№", "Код ТН ВЭД", "Описание", "Кол-во", "Ед.изм.", "Вес нетто (кг)", "Вес брутто (кг)", "Стоимость", "Пошлина", "НДС"}
		for i, h := range headers {
			cell := fmt.Sprintf("%s%d", string(rune('A'+i)), row)
			f.SetCellValue(sheetName, cell, h)
			f.SetCellStyle(sheetName, cell, cell, headerStyle)
		}
		row++

		for _, item := range resp.Items {
			cols := []string{"number", "hs_code", "description", "quantity", "unit", "weight_net", "weight_gross", "value", "duty_rate", "vat_rate"}
			for i, col := range cols {
				cell := fmt.Sprintf("%s%d", string(rune('A'+i)), row)
				val, ok := item[col]
				if ok {
					f.SetCellValue(sheetName, cell, val)
				}
				f.SetCellStyle(sheetName, cell, cell, cellStyle)
			}
			row++
		}
	}

	f.SetColWidth(sheetName, "A", "A", 25)
	f.SetColWidth(sheetName, "B", "B", 20)
	f.SetColWidth(sheetName, "C", "C", 35)
	f.SetColWidth(sheetName, "D", "J", 15)

	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
