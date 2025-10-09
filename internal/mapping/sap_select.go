package mapping

// MapSAPtoInternal задаёт соответствия: internal → SAP header.
// v1: хардкод; v1.1 — поднимем из configs/mapping/sap_to_staging.yaml
func MapSAPtoInternal() map[string]string {
	return map[string]string{
		"client_id":             "קוד כרטיס",
		"client_license_number": "מספר עוסק מורשה",
		"client_name_he":        "שם לועזי", // имя клиента (первая колонка F)
		"client_name_ru":        "שם כרטיס", // русское имя
		"client_address":        "כתובת",
		"order_id":              "אסמכתת בסיס",
		"date_raw":              "תאריך אסמכתא",
		"total_weight_raw":      "סה'כ משקל",
		"total_packaging_raw":   "סה'כ אריזות",
		"district_ru":           "מחוז",
	}
}

// ApplySelect переименовывает: оставляет только нужные поля.
func ApplySelect(raw []map[string]string, sel map[string]string) []map[string]string {
	out := make([]map[string]string, 0, len(raw))
	for _, r := range raw {
		row := map[string]string{}
		for internal, sapName := range sel {
			row[internal] = r[sapName]
		}
		out = append(out, row)
	}
	return out
}
