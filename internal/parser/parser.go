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

// ParseTradeRow maps a row of strings to a Trade model, using baseDate to resolve relative dates.
func ParseTradeRow(row []string, baseDate string) models.Trade {
	t := models.Trade{
		OrderDate:     ParseFlexibleDate(SafeGet(row, 0), baseDate),
		DeliveryDate:  ParseFlexibleDate(SafeGet(row, 1), baseDate),
		OrderType:     SafeGet(row, 2),
		SellerOrderID: SafeGet(row, 3),
		OrderID:       SafeGet(row, 4),
		Inventory:     SafeGet(row, 4),
		Status:        SafeGet(row, 5),
		Tracking:      SafeGet(row, 6),
		Created:       time.Now(),
	}

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

	if len(row) > 26 {
		t.Inventory = SafeGet(row, 26)
	}

	return t
}

// ParseFlexibleDate tries to parse 'Sun, Feb 1' using a base year/month (YYYY-MM).
func ParseFlexibleDate(dateStr, baseDate string) time.Time {
	dateStr = strings.TrimSpace(dateStr)
	if dateStr == "" {
		return time.Time{}
	}

	// Case 1: Standard YYYY-MM-DD
	if t, err := time.Parse("2006-01-02", dateStr); err == nil {
		return t
	}

	// Case 2: "Sun, Feb 1" or "Feb 1"
	// Clean string: remove commas
	cleanStr := strings.ReplaceAll(dateStr, ",", "")
	parts := strings.Fields(cleanStr)
	if len(parts) > 0 {
		dayStr := parts[len(parts)-1]
		// Convert day to 2 digits
		dayInt := 0
		fmt.Sscanf(dayStr, "%d", &dayInt)
		if dayInt > 0 {
			fullDateStr := fmt.Sprintf("%s-%02d", baseDate, dayInt)
			if t, err := time.Parse("2006-01-02", fullDateStr); err == nil {
				return t
			}
		}
	}

	return time.Time{}
}

// LoadTradesFromFile loads trade records from a CSV or Excel file.
func LoadTradesFromFile(path, baseDate string) ([]models.Trade, error) {
	if strings.HasSuffix(path, ".csv") {
		return LoadCSV(path, baseDate)
	} else if strings.HasSuffix(path, ".xlsx") {
		return LoadExcel(path, baseDate)
	}
	return nil, fmt.Errorf("unsupported file format: %s", path)
}

// LoadCSV reads trade records from a CSV file.
func LoadCSV(path, baseDate string) ([]models.Trade, error) {
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
		trades = append(trades, ParseTradeRow(record, baseDate))
	}
	return trades, nil
}

// LoadExcel reads trade records from an Excel file.
func LoadExcel(path, baseDate string) ([]models.Trade, error) {
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
		trades = append(trades, ParseTradeRow(row, baseDate))
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
