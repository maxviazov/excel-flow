-- Добавляем алиасы для отсутствующих городов из файла משקל.xlsx

-- קרית ים -> קריית ים (D1189)
INSERT OR IGNORE INTO city_aliases (alias_heb, target_heb) 
VALUES ('קרית ים', 'קריית ים');

-- הרצליה -> הרצלייה (I372)  
INSERT OR IGNORE INTO city_aliases (alias_heb, target_heb)
VALUES ('הרצליה', 'הרצלייה');