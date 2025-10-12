#!/usr/bin/env python3
import sqlite3
import sys
import xlrd

wb = xlrd.open_workbook('testdata/עותק של רשימת יישובים עדכנית 16.4.24.xlt')
ws = wb.sheet_by_index(0)

print(f"Строк: {ws.nrows}, Колонок: {ws.ncols}")

headers = [ws.cell_value(0, col) for col in range(ws.ncols)]
print(f"Заголовки: {headers[:5]}")

city_col = 0
code_col = 1

print(f"Колонка города: {city_col}, Колонка кода: {code_col}")

conn = sqlite3.connect('configs/dictionaries/city.db')
cursor = conn.cursor()

cursor.execute("DELETE FROM city_aliases")
cursor.execute("DELETE FROM city_codes")

count = 0
for row_idx in range(1, ws.nrows):
    city = ws.cell_value(row_idx, city_col)
    code = ws.cell_value(row_idx, code_col)
    
    if city and code:
        city = str(city).strip()
        code = str(code).strip()
        
        try:
            cursor.execute(
                "INSERT INTO city_codes (city_heb, city_code) VALUES (?, ?)",
                (city, code)
            )
            count += 1
        except sqlite3.IntegrityError:
            pass

conn.commit()
conn.close()

print(f"\nЗагружено {count} городов")
