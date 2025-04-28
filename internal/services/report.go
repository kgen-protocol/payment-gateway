package services

import (
	"fmt"
	"log"

	"github.com/xuri/excelize/v2"
)

type ProductSaveError struct {
	ProductID int
	ErrorMsg  string
}

type FetchPageError struct {
	Page     int
	ErrorMsg string
}

func GenerateProductSyncReport(successCount, saveErrorCount, fetchErrorCount int, saveErrors []ProductSaveError, fetchErrors []FetchPageError) error {
	f := excelize.NewFile()

	// Create Summary Sheet
	summary := "Summary"
	f.NewSheet(summary)
	f.SetCellValue(summary, "A1", "Total Products Synced")
	f.SetCellValue(summary, "B1", successCount)

	f.SetCellValue(summary, "A2", "Product Save Errors")
	f.SetCellValue(summary, "B2", saveErrorCount)

	f.SetCellValue(summary, "A3", "Page Fetch Errors")
	f.SetCellValue(summary, "B3", fetchErrorCount)

	// Create Save Errors Sheet
	saveErrorSheet := "Product Save Errors"
	f.NewSheet(saveErrorSheet)
	f.SetCellValue(saveErrorSheet, "A1", "ProductID")
	f.SetCellValue(saveErrorSheet, "B1", "Error Message")

	for idx, errLog := range saveErrors {
		rowNum := idx + 2
		f.SetCellValue(saveErrorSheet, fmt.Sprintf("A%d", rowNum), errLog.ProductID)
		f.SetCellValue(saveErrorSheet, fmt.Sprintf("B%d", rowNum), errLog.ErrorMsg)
	}

	// Create Fetch Errors Sheet
	fetchErrorSheet := "Fetch Errors"
	f.NewSheet(fetchErrorSheet)
	f.SetCellValue(fetchErrorSheet, "A1", "Page")
	f.SetCellValue(fetchErrorSheet, "B1", "Error Message")

	for idx, errLog := range fetchErrors {
		rowNum := idx + 2
		f.SetCellValue(fetchErrorSheet, fmt.Sprintf("A%d", rowNum), errLog.Page)
		f.SetCellValue(fetchErrorSheet, fmt.Sprintf("B%d", rowNum), errLog.ErrorMsg)
	}

	// Delete default sheet
	f.DeleteSheet("Sheet1")

	// Save file
	err := f.SaveAs("product_sync_report.xlsx")
	if err != nil {
		return fmt.Errorf("failed to save xlsx file: %w", err)
	}

	log.Println("Product sync report saved as 'product_sync_report.xlsx'")
	return nil
}
