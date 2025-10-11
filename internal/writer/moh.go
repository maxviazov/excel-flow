package writer

import (
	"fmt"
	"sort"
	"time"

	"github.com/xuri/excelize/v2"

	"github.com/maxviazov/excel-flow/internal/pipelines"
)

func WriteMOH(path string, groups map[pipelines.GroupKey]*pipelines.GroupVal) error {
	f := excelize.NewFile()
	sh := "data"
	f.NewSheet(sh)
	f.DeleteSheet("Sheet1")

	// Set RTL (right-to-left) for Hebrew
	f.SetSheetView(sh, 0, &excelize.ViewOptions{RightToLeft: boolPtr(true)})

	// MOH headers in exact order
	headers := []string{
		"שם הספק",
		"ח\"פ ספק",
		"מספר משרד הבריאות",
		"תאריך",
		"מס.רכב",
		"שם הנהג",
		"טלפון נהג",
		"לקוח",
		"סוג לקוח (קמעונאי,מפעל/מחסן)",
		"קוד עיר",
		"כתובת",
		"ח\"פ לקוח / מספר אישור משרד הבריאות",
		"מספר סניף הרשת",
		"מספר תעודת משלוח",
		"בשר בהמות גולמי (משקל)",
		"בשר בהמות מיבוא קפוא (משקל)",
		"בשר בהמות מעובד (משקל)",
		"עוף גולמי (עוף שחוט) (משקל)",
		"עוף מעובד (משקל)",
		"דגים גולמי (מקומי) (משקל)",
		"דגים יבוא",
		"דגים מעובדים",
		"מוצרים מוכנים לאכילה",
		"נוסף א",
		"נוסף ב",
		"סה\"כ קרטונים",
		"סה\"כ משקל",
		"סבב יומי",
		"קוד ביטול דיווח משלוח\n(למקרים בהם נדרש לבטל תעודת משלוח שדווחה ולא יצאה מהמפעל לשיווק",
		"משווק באמצעות",
	}

	// Write headers
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sh, cell, h)
	}

	// Sort keys for stability
	keys := make([]pipelines.GroupKey, 0, len(groups))
	for k := range groups {
		keys = append(keys, k)
	}
	sort.Slice(
		keys, func(i, j int) bool {
			if keys[i].ClientLicense != keys[j].ClientLicense {
				return keys[i].ClientLicense < keys[j].ClientLicense
			}
			if keys[i].OrderID != keys[j].OrderID {
				return keys[i].OrderID < keys[j].OrderID
			}
			return keys[i].Date.Before(keys[j].Date)
		},
	)

	// Write data rows
	row := 2
	for _, k := range keys {
		v := groups[k]
		// Set non-zero values only
		setCellIfNotEmpty := func(col int, val any) {
			cell, _ := excelize.CoordinatesToCellName(col, row)
			if str, ok := val.(string); ok && str != "" && str != "0" {
				f.SetCellValue(sh, cell, val)
			} else if num, ok := val.(float64); ok && num != 0 {
				f.SetCellValue(sh, cell, fmt.Sprintf("%.3f", num))
			}
		}

		// Set required fields
		setCellIfNotEmpty(1, "דולינה גרופ בע\"מ")
		setCellIfNotEmpty(2, "511777856")
		setCellIfNotEmpty(3, "P1908")
		setCellIfNotEmpty(4, time.Now().Format("2006-01-02"))
		// 5-7: vehicle info - empty
		setCellIfNotEmpty(8, v.ClientName)
		setCellIfNotEmpty(9, "קמעונאי")
		setCellIfNotEmpty(10, v.CityCode) // код города
		setCellIfNotEmpty(11, v.Address)
		setCellIfNotEmpty(12, k.ClientLicense)
		// 13: branch number - empty
		setCellIfNotEmpty(14, k.OrderID) // один номер документа на строку
		// 15-22: meat/fish categories - all zero, skip
		setCellIfNotEmpty(23, v.TotalWeight) // מוצרים מוכנים לאכילה
		// 24-25: additional categories - zero, skip
		setCellIfNotEmpty(26, v.TotalPackages) // סה"כ קרטונים
		setCellIfNotEmpty(27, v.TotalWeight)   // סה"כ משקל
		setCellIfNotEmpty(28, "1")             // סבב יומי
		// 29-30: cancellation/marketing - empty
		row++
	}

	return f.SaveAs(path)
}

func boolPtr(b bool) *bool {
	return &b
}
