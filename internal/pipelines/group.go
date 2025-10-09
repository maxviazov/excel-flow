package pipelines

import (
	"fmt"
	"time"
)

type GroupKey struct {
	ClientLicense string
	OrderID       string
	Date          time.Time
}

type GroupVal struct {
	Rows          []Row
	TotalWeight   float64
	TotalPackages float64
	Count         int
	ClientName    string
	Address       string
	OrderIDs      []string // список всех номеров документов
}

// BuildGroups: кидаем строки с пустыми ключами; остальное агрегируем.
func BuildGroups(rows []Row) (map[GroupKey]*GroupVal, []Issue) {
	g := map[GroupKey]*GroupVal{}
	issues := []Issue{}

	for i, r := range rows {
		lic, _ := r["client_license_number"].(string)
		ord, _ := r["order_id"].(string)
		dt, _ := r["date"].(time.Time)

		if lic == "" || ord == "" || dt.IsZero() {
			issues = append(
				issues, Issue{
					Level:  "warn",
					Code:   "DROP_NO_KEY",
					RowNum: i + 2,
					Detail: fmt.Sprintf("missing key parts: lic=%q ord=%q date=%v", lic, ord, dt),
				},
			)
			continue
		}

		key := GroupKey{ClientLicense: lic, OrderID: ord, Date: dt}
		if g[key] == nil {
			g[key] = &GroupVal{}
		}
		val := g[key]
		val.Rows = append(val.Rows, r)
		if val.ClientName == "" {
			if s, ok := r["client_name_he"].(string); ok && s != "" {
				val.ClientName = s
			}
		}
		if val.Address == "" {
			if s, ok := r["client_address"].(string); ok && s != "" {
				val.Address = s
			}
		}

		if w, ok := r["total_weight"].(float64); ok {
			val.TotalWeight += w
		}
		if p, ok := r["total_packaging"].(float64); ok {
			val.TotalPackages += p
		}
		val.Count++
	}

	return g, issues
}
