# Cheat-Sheet: Configuration `משקל.yaml`

## 🧩 General Pipeline Flow

| Stage               | When Triggered    | What It Does                                                                               | Where It Logs         |
|---------------------|-------------------|--------------------------------------------------------------------------------------------|-----------------------|
| **source**          | Start of program  | Validates file, sheet, header row.                                                         | `info` / `error`      |
| **columns / types** | Reading file      | Maps Excel headers to internal keys, sets data types.                                      | `warn` / `error`      |
| **transforms**      | After read        | Trims spaces, normalizes Unicode, applies decimal separator and weight multiplier.         | `debug`               |
| **filters**         | Pre-processing    | Removes rows with negative or missing values.                                              | `info` (summary)      |
| **grouping**        | After filtering   | Aggregates by `client_license_number`, `order_id`. Sums weights, packages, sets min(date). | `debug`               |
| **city_extraction** | After aggregation | Extracts city from address or district, translates RU→HEB, looks up city code.             | `warn` (missing city) |
| **output**          | Final step        | Writes `output.xlsx` with mapped headers.                                                  | `info` (success)      |
| **error_handling**  | Any error         | Defines behavior on missing/invalid data.                                                  | `warn` / `error`      |
| **performance**     | During processing | Batch size, concurrency settings.                                                          | `debug`               |
| **logging**         | Continuous        | Controls verbosity, color, timestamp.                                                      | All levels            |

---

## 🌆 City Extraction Logic

1. Extract city from `client_address` (text before first comma).
2. If empty → take from `district_ru` (column מחוז).
3. Translate Russian → Hebrew using `data/ru_to_heb.csv`.
4. Find city code by Hebrew name in `data/city_aliases.csv`.
5. If not found → assign fallback code `9999`.

**Outputs:**

* `city_raw` — raw value from address or district.
* `city_heb` — normalized name in Hebrew.
* `city_code` — city numeric code.

---

## 📃 Supporting Data Files

| File                    | Purpose                      | Format     |
|-------------------------|------------------------------|------------|
| `משקל.xlsx`             | Source SAP report            | Excel      |
| `data/ru_to_heb.csv`    | Russian → Hebrew translation | `ru,heb`   |
| `data/city_aliases.csv` | Hebrew name → city code      | `heb,code` |
| `output.xlsx`           | Generated report             | Excel      |
| `log.txt`               | Log of all operations        | text       |

---

## 🔐 Error Handling

| Condition         | Behavior                    | Example Log                                  |
|-------------------|-----------------------------|----------------------------------------------|
| Missing column    | `fail`                      | `❌ Missing required column: תאריך אסמכתא`    |
| Invalid type      | `warn`                      | `⚠ Value 'abc' in weight column not numeric` |
| Missing city code | `use_fallback` → `9999`     | `⚠ City 'Назарет' not found, fallback=9999`  |
| Too many errors   | Stops after N (default 100) | `⛔ Max error threshold reached (100)`        |

---

## 📊 Sample Log Events

| Level   | Example                                   |
|---------|-------------------------------------------|
| `INFO`  | `✅ Loaded 2350 rows from משקל.xlsx`       |
| `INFO`  | `⛔ 14 rows skipped (negative weight)`     |
| `WARN`  | `⚠ City 'Хайфа' not found, fallback=9999` |
| `DEBUG` | `↪ Grouped 2300 → 214 records`            |
| `ERROR` | `❌ File not found: משקל.xlsx`             |

---

## 🛠 Performance Notes

* `batch_size`: processes rows in groups (1000 default).
* `parallel: false`: sequential for safety; can switch to `true` for large datasets.
* Recommended first-run mode: **debug off, parallel false** to ensure full trace.

---

## 🔗 Quick Reference of Keys

| Section           | Key                         | Description                               |
|-------------------|-----------------------------|-------------------------------------------|
| `filters`         | `skip_negative_rows.fields` | Columns checked for `< 0` values          |
| `grouping`        | `aggregates`                | Defines how numeric fields are summarized |
| `city_extraction` | `ru_to_heb_map`             | Path to RU→HEB translation table          |
| `city_extraction` | `city_code_map`             | Path to Hebrew→Code map                   |
| `error_handling`  | `on_missing_city_code`      | Behavior if city not found                |
| `output`          | `map`                       | Final field mapping to Excel headers      |

---

### ✅ Summary

This configuration defines a complete ETL pipeline for SAP shipment reports:

* Reads Excel → Cleans → Aggregates → Extracts/Translates City → Writes Output.
* Provides robust logging, filtering, and fallback mechanisms.
* 100% compatible with localized (Hebrew/Russian) data environments.
