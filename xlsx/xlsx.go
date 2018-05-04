package xlsx

import (
	"math"
	"strconv"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/pulpfree/gales-fuelsale-report/model"
)

const (
	abc            = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	timeShortForm  = "20060102"
	timeRecordForm = "2006-01-02"
)

var sheetName = "Period"

// NewFile function
func NewFile(fs *model.FuelSales) (xlsx *excelize.File, err error) {

	sDte := fs.DateStart.Format(timeShortForm)
	eDte := fs.DateEnd.Format(timeShortForm)
	sheetName = sheetName + " " + sDte + " - " + eDte

	xlsx = excelize.NewFile()
	xlsx.SetSheetName("Sheet1", sheetName)
	xlsx.SetColWidth(sheetName, "A", "E", 12)

	style, err := xlsx.NewStyle(`{"font":{"color":"#666666", "size":12}}`)
	if err != nil {
		return xlsx, err
	}
	xlsx.SetCellStyle(sheetName, "A1", "E1", style)

	xlsx.SetCellValue(sheetName, "A1", "Station")
	xlsx.SetCellValue(sheetName, "B1", "Date")
	xlsx.SetCellValue(sheetName, "C1", "NL")
	xlsx.SetCellValue(sheetName, "D1", "SNL")
	xlsx.SetCellValue(sheetName, "E1", "DSL")
	xlsx.SetCellValue(sheetName, "F1", "CDSL")

	var cell string
	col := 1
	row := 2

	for _, sale := range fs.Sales {

		f := sale.Fuel
		// Format date to UTC so we're getting the right day/time
		dte := sale.Date.UTC().Format(timeRecordForm)

		cell = toChar(col) + strconv.Itoa(row)
		xlsx.SetCellValue(sheetName, cell, sale.StationName)
		col++

		cell = toChar(col) + strconv.Itoa(row)
		xlsx.SetCellValue(sheetName, cell, dte)
		col++

		cell = toChar(col) + strconv.Itoa(row)
		xlsx.SetCellValue(sheetName, cell, toFixed(f.NL, 2))
		col++

		cell = toChar(col) + strconv.Itoa(row)
		xlsx.SetCellValue(sheetName, cell, toFixed(f.SNL, 2))
		col++

		cell = toChar(col) + strconv.Itoa(row)
		xlsx.SetCellValue(sheetName, cell, toFixed(f.DSL, 2))
		col++

		cell = toChar(col) + strconv.Itoa(row)
		xlsx.SetCellValue(sheetName, cell, toFixed(f.CDSL, 2))
		col = 1
		row++
	}

	return xlsx, err
}

// see: https://stackoverflow.com/questions/36803999/golang-alphabetic-representation-of-a-number
// for a way to map int to letters
func toChar(i int) string {
	return abc[i-1 : i]
}

// Found these function at: https://stackoverflow.com/questions/18390266/how-can-we-truncate-float64-type-to-a-particular-precision-in-golang
// Looks like a good way to deal with precision
func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}
