package parser

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
	"time"

	"mongotest/internal/models"

	"github.com/xuri/excelize/v2"
)

// ParseTradeRow maps a row of strings to a Trade model.
func ParseTradeRow(row []string) models.Trade {
	t := models.Trade{
		OrderDate:     SafeGet(row, 0),
		DeliveryDate:  SafeGet(row, 1),
		OrderType:     SafeGet(row, 2),
		SellerOrderID: SafeGet(row, 3),
		OrderID:       SafeGet(row, 4), // Some CSVs have OrderID at index 3 or 4
		Inventory:     SafeGet(row, 4), // In some layouts
		Status:        SafeGet(row, 5),
		Tracking:      SafeGet(row, 6),
		Created:       time.Now(),
	}

	// Adjusting for jumbled indices if needed, but keeping the core logic consistent
	t.PackageCount = ParseInt(SafeGet(row, 7))
	t.Qty = ParseInt(SafeGet(row, 8))
	t.SalePrice1 = ParseFloat(SafeGet(row, 9))
	t.SalePrice2 = ParseFloat(SafeGet(row, 10))
	t.ShippingCharged = ParseFloat(SafeGet(row, 11))
	t.AMZFee = ParseFloat(SafeGet(row, 12))
	t.TotalRate = ParseFloat(SafeGet(row, 13))
	t.TaxFees = ParseFloat(SafeGet(row, 14))
	t.Shipping = ParseFloat(SafeGet(row, 15))
	t.UnitCost = ParseFloat(SafeGet(row, 16))
	t.TotalCost = ParseFloat(SafeGet(row, 17))
	t.Refund = ParseFloat(SafeGet(row, 18))
	t.Profit = ParseFloat(SafeGet(row, 19))
	t.Loss = ParseFloat(SafeGet(row, 20))
	t.ROI = ParseFloat(SafeGet(row, 21))
	t.ReturnRefund = ParseFloat(SafeGet(row, 22))
	t.VeeqoCredits = ParseFloat(SafeGet(row, 23))
	t.ThreePLCost = ParseFloat(SafeGet(row, 24))
	t.Net = ParseFloat(SafeGet(row, 25))

	// Re-check some mappings that vary between orders and sales
	if len(row) > 26 {
		t.Inventory = SafeGet(row, 26)
	}

	return t
}

// LoadTradesFromFile loads trade records from a CSV or Excel file.
func LoadTradesFromFile(path string) ([]models.Trade, error) {
	if strings.HasSuffix(path, ".csv") {
		return LoadCSV(path)
	} else if strings.HasSuffix(path, ".xlsx") {
		return LoadExcel(path)
	}
	return nil, fmt.Errorf("unsupported file format: %s", path)
}

// LoadCSV reads trade records from a CSV file.
func LoadCSV(path string) ([]models.Trade, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	reader := csv.NewReader(f)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var trades []models.Trade
	for i, record := range records {
		if i == 0 || len(record) < 9 {
			continue // skip header or invalid rows
		}
		trades = append(trades, ParseTradeRow(record))
	}
	return trades, nil
}

// LoadExcel reads trade records from an Excel file.
func LoadExcel(path string) ([]models.Trade, error) {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		return nil, fmt.Errorf("no sheets found in excel file")
	}

	rows, err := f.GetRows(sheets[0])
	if err != nil {
		return nil, err
	}

	var trades []models.Trade
	for i, row := range rows {
		if i == 0 || len(row) < 9 {
			continue
		}
		trades = append(trades, ParseTradeRow(row))
	}
	return trades, nil
}

// SafeGet safely retrieves a trimmed string from a row at the given index.
func SafeGet(row []string, index int) string {
	if index < len(row) {
		return strings.TrimSpace(row[index])
	}
	return ""
}

// ParseInt parses a string into an integer.
func ParseInt(s string) int {
	var val int
	fmt.Sscanf(s, "%d", &val)
	return val
}

// ParseFloat parses a string into a float64, handling currency symbols and commas.
func ParseFloat(s string) float64 {
	s = strings.ReplaceAll(s, "$", "")
	s = strings.ReplaceAll(s, ",", "")
	var val float64
	fmt.Sscanf(s, "%f", &val)
	return val
}
