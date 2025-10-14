package writer

import (
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/xuri/excelize/v2"

	"github.com/maxviazov/excel-flow/internal/drivers"
	"github.com/maxviazov/excel-flow/internal/pipelines"
	"github.com/maxviazov/excel-flow/internal/textutil"
)

func WriteMOH(path string, groups map[pipelines.GroupKey]*pipelines.GroupVal, driverRegistry *drivers.Registry) error {
	f := excelize.NewFile()
	sh := "Sheet1"
	
	// Create header row from official MOH template
	headers := []string{
		"שם הספק",
		"ח\"פ ספק ",
		"מספר משרד הבריאות",
		"תאריך",
		"מס.רכב",
		"שם הנהג",
		"טלפון נהג",
		"לקוח",
		"סוג לקוח (קמעונאי,מפעל/מחסן)",
		"קוד עיר",
		"כתובת",
		"ח\"פ לקוח \nאו מספר אישור משרד הבריאות במקרים בהם המשלוח הוא למפעל מאושר",
		"מספר סניף הרשת",
		"מספר תעודת משלוח",
		"בשר בהמות גולמי",
		"בשר בהמות מיבוא קפוא",
		"בשר בהמות מעובד",
		"עוף גולמי (עוף שחוט)",
		"עוף מעובד",
		"דגים גולמי (מקומי)",
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
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellStr(sh, cell, h)
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
		
		// Set values according to MOH template
		f.SetCellStr(sh, fmt.Sprintf("A%d", row), "דולינה גרופ בע\"מ") // שם הספק
		f.SetCellValue(sh, fmt.Sprintf("B%d", row), 511777856) // ח"פ ספק
		f.SetCellStr(sh, fmt.Sprintf("C%d", row), "P1908") // מספר משרד הבריאות
		f.SetCellValue(sh, fmt.Sprintf("D%d", row), time.Now()) // תאריך
		
		// 5-7: vehicle and driver info
		if driverRegistry != nil {
			driver := driverRegistry.GetRandomDriverForCity(v.CityCode)
			if driver != nil {
				f.SetCellStr(sh, fmt.Sprintf("E%d", row), driver.LicenseNumber) // מס.רכב
				f.SetCellStr(sh, fmt.Sprintf("F%d", row), driver.Name) // שם הנהג
				f.SetCellStr(sh, fmt.Sprintf("G%d", row), driver.Phone) // טלפון נהג
			}
		}
		
		if v.ClientName != "" {
			f.SetCellStr(sh, fmt.Sprintf("H%d", row), textutil.SanitizeForMOH(v.ClientName)) // לקוח
		}
		f.SetCellStr(sh, fmt.Sprintf("I%d", row), "קמעונאי") // סוג לקוח
		if v.CityCode != "" {
			f.SetCellStr(sh, fmt.Sprintf("J%d", row), v.CityCode) // קוד עיר
		}
		if v.Address != "" {
			f.SetCellStr(sh, fmt.Sprintf("K%d", row), textutil.SanitizeForMOH(v.Address)) // כתובת
		}
		f.SetCellValue(sh, fmt.Sprintf("L%d", row), parseNumber(k.ClientLicense)) // ח"פ לקוח
		// M is empty - מספר סניף הרשת
		f.SetCellValue(sh, fmt.Sprintf("N%d", row), parseNumber(k.OrderID)) // מספר תעודת משלוח
		// O-Y are empty (meat/fish categories)
		f.SetCellValue(sh, fmt.Sprintf("Z%d", row), v.TotalPackages) // סה"כ קרטונים
		f.SetCellValue(sh, fmt.Sprintf("AA%d", row), v.TotalWeight) // סה"כ משקל
		f.SetCellValue(sh, fmt.Sprintf("AB%d", row), 1) // סבב יומי
		// AC, AD are empty
		row++
	}

	return f.SaveAs(path)
}

func boolPtr(b bool) *bool {
	return &b
}

func parseNumber(s string) interface{} {
	if num, err := strconv.ParseFloat(s, 64); err == nil {
		return num
	}
	return s
}
