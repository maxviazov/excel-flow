package pipelines

// Порядок и имена колонок — как в шаблоне МОЗ (упрощённая часть)
var MOHHeaders = []string{
	"שם הספק",           // supplier_name
	"ח\"פ ספק",          // supplier_tax_id
	"מספר משרד הבריאות", // moh_registration
	"תאריך",             // shipment_date (YYYY-MM-DD)
	"מס.רכב",            // vehicle_number
	"שם הנהג",           // driver_name
	"טלפון נהג",         // driver_phone
	"לקוח",              // client_name
	"סוג לקוח (קמעונאי,מפעל/מחסן)", // client_type
	"קוד עיר", // city_code
	"כתובת",   // address
	"ח\"פ לקוח / מספר אישור משרד הבריאות", // client_license_number
	"מספר סניף הרשת",                      // branch_number
	"מספר תעודת משלוח",                    // delivery_note_number

	// категории O..X — заполняем нулями
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

	"סה\"כ קרטונים", // total_cartons
	"סה\"כ משקל",    // total_weight
	"סבב יומי",      // daily_round
	"קוד ביטול דיווח משלוח\n(למקרים בהם נדרש לבטל תעודת משלוח שדווחה ולא יצאה מהמפעל לשיווק",
	"משווק באמצעות",
}

// Константы предприятия (вынеси куда хочешь)
type SupplierPreset struct {
	Name       string // שם הספק
	TaxID      string // ח"פ ספק
	MOHReg     string // מספר משרד הבריאות
	ClientType string // סוג לקוח
	DailyRound int    // סבב יומי
}

type LogisticsPreset struct {
	VehicleNumber string
	DriverName    string
	DriverPhone   string
}

func ProjectGroupsToMOH(groups map[GroupKey]*GroupVal, sup SupplierPreset) []map[string]any {
	out := make([]map[string]any, 0, len(groups))
	for k, v := range groups {
		row := map[string]any{
			"שם הספק":           sup.Name,
			"ח\"פ ספק":          sup.TaxID,
			"מספר משרד הבריאות": sup.MOHReg,
			"תאריך":             k.Date.Format("2006-01-02"),

			"מס.רכב":    "", // пока пусто; позже подставим из logistics
			"שם הנהג":   "",
			"טלפון נהג": "",

			"לקוח": v.ClientName,
			"סוג לקוח (קמעונאי,מפעל/מחסן)": sup.ClientType,
			"קוד עיר": "", // пока пусто (добавим lookup позже)
			"כתובת":   v.Address,

			"ח\"פ לקוח / מספר אישור משרד הבריאות": k.ClientLicense,
			"מספר סניף הרשת":                      "",
			"מספר תעודת משלוח":                    k.OrderID,

			// категории — нули
			"בשר בהמות גולמי (משקל)":      0.0,
			"בשר בהמות מיבוא קפוא (משקל)": 0.0,
			"בשר בהמות מעובד (משקל)":      0.0,
			"עוף גולמי (עוף שחוט) (משקל)": 0.0,
			"עוף מעובד (משקל)":            0.0,
			"דגים גולמי (מקומי) (משקל)":   0.0,
			"דגים יבוא":                   0.0,
			"דגים מעובדים":                0.0,
			"מוצרים מוכנים לאכילה":        v.TotalWeight, // твой основной вес
			"נוסף א":                      0.0,
			"נוסף ב":                      0.0,

			"סה\"כ קרטונים": v.TotalPackages,
			"סה\"כ משקל":    v.TotalWeight,
			"סבב יומי":      sup.DailyRound,
			"קוד ביטול דיווח משלוח\n(למקרים בהם נדרש לבטל תעודת משלוח שדווחה ולא יצאה מהמפעל לשיווק": "",
			"משווק באמצעות": "",
		}
		out = append(out, row)
	}
	return out
}
