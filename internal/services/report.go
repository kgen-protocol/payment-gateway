package services

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/aakritigkmit/payment-gateway/internal/model"
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

func ExportProductsToExcel(products []model.Product) error {
	f := excelize.NewFile()
	sheet := "Products"
	f.NewSheet(sheet)

	headers := []string{
		"Product ID", "Name", "Type", "Description", "Availability Zones", "Tags",
		"Destination", "Source", "Operator", "Service", "SubService", "Prices", "Rates",
		"Validity", "Number of Benefits", "Benefit", "Required Credit Party Identifier Fields",
		"Required Sender Fields", "Required Beneficiary Fields", "Required Debit Party Identifier Fields",
		"Required Additional Identifier Fields", "Required Statement Identifier Fields", "Promotions", "Regions",
	}

	// Write headers
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
	}

	// Write data
	for i, p := range products {
		row := i + 2

		// Helper functions
		jsonStr := func(v interface{}) string {
			if v == nil {
				return ""
			}
			b, _ := json.Marshal(v)
			return string(b)
		}

		joinStrings := func(arr []string) string {
			return strings.Join(arr, ", ")
		}

		var flatCreditFields []string

		for _, arr := range p.RequiredCreditPartyIdentifierFields {
			flatCreditFields = append(flatCreditFields, strings.Join(arr, "|"))
		}

		values := []interface{}{
			p.UniqueId,
			p.Name,
			p.Type,
			p.Description,
			joinStrings(p.AvailabilityZones),
			jsonStr(p.Tags),
			jsonStr(p.Destination),
			jsonStr(p.Source),
			jsonStr(p.Operator),
			jsonStr(p.Service),
			jsonStr(p.Service.SubService),
			jsonStr(p.Prices),
			jsonStr(p.Rates),
			jsonStr(p.Validity),
			len(p.Benefits),
			jsonStr(p.Benefits),
			strings.Join(flatCreditFields, ", "),
			jsonStr(p.RequiredSenderFields),
			jsonStr(p.RequiredBeneficiaryFields),
			jsonStr(p.RequiredDebitPartyIdentifierFields),
			jsonStr(p.RequiredAdditionalIdentifierFields),
			jsonStr(p.RequiredStatementIdentifierFields),
			jsonStr(p.Promotions),
			jsonStr(p.Regions),
		}

		for j, val := range values {
			cell, _ := excelize.CoordinatesToCellName(j+1, row)
			f.SetCellValue(sheet, cell, val)
		}
	}

	f.DeleteSheet("Sheet1")
	if err := f.SaveAs("product_fetch_report.xlsx"); err != nil {
		return fmt.Errorf("failed to save Excel file: %w", err)
	}

	log.Println("Product fetch report saved as 'product_fetch_report.xlsx'")
	return nil
}
