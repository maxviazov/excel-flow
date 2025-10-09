package normalize

import (
	"strconv"
	"strings"
	"time"
)

type NormRow = map[string]any

// Normalize принимает строковые поля и создаёт типизированную запись.
func Normalize(in map[string]string) (NormRow, []string) {
	var issues []string

	trim := func(s string) string { return strings.TrimSpace(s) }
	num := func(s string) (float64, bool) {
		s = strings.ReplaceAll(trim(s), ",", ".")
		if s == "" {
			return 0, false
		}
		f, err := strconv.ParseFloat(s, 64)
		return f, err == nil
	}
	parseDate := func(s string) (time.Time, bool) {
		s = trim(s)
		// v1: пробуем ISO и dd/mm/yy; (excel-serial добавим в v1.1)
		if t, err := time.Parse("2006-01-02", s); err == nil {
			return t, true
		}
		if t, err := time.Parse("02/01/06", s); err == nil {
			return t, true
		}
		if t, err := time.Parse("02/01/2006", s); err == nil {
			return t, true
		}
		return time.Time{}, false
	}

	out := NormRow{
		"client_id":             trim(in["client_id"]),
		"client_license_number": trim(in["client_license_number"]),
		"client_name_he":        trim(in["client_name_he"]),
		"client_name_ru":        trim(in["client_name_ru"]),
		"client_address":        trim(in["client_address"]),
		"order_id":              trim(in["order_id"]),
		"district_ru":           trim(in["district_ru"]),
	}

	if w, ok := num(in["total_weight_raw"]); ok && w >= 0 {
		out["total_weight"] = w // данные уже в кг
	} else {
		out["total_weight"] = 0.0
		issues = append(issues, "BAD_WEIGHT")
	}
	if p, ok := num(in["total_packaging_raw"]); ok && p >= 0 {
		out["total_packaging"] = p
	} else {
		out["total_packaging"] = 0.0
		issues = append(issues, "BAD_PACKAGING")
	}
	if d, ok := parseDate(in["date_raw"]); ok {
		out["date"] = d
	} else {
		issues = append(issues, "BAD_DATE")
	}

	return out, issues
}
