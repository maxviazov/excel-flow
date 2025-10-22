-- Добавляем алиасы для орфографических вариантов городов

-- קירית ביאליק -> קריית ביאליק (D1186)
INSERT OR IGNORE INTO city_aliases (alias_heb, target_heb)
VALUES ('קירית ביאליק', 'קריית ביאליק');

-- קרית ביאליק -> קריית ביאליק (D1186) - еще один вариант
INSERT OR IGNORE INTO city_aliases (alias_heb, target_heb)
VALUES ('קרית ביאליק', 'קריית ביאליק');

-- נהריה -> נהרייה (B917)
INSERT OR IGNORE INTO city_aliases (alias_heb, target_heb)
VALUES ('נהריה', 'נהרייה');

-- מעלות -> מעלות-תרשיחא (B865)
INSERT OR IGNORE INTO city_aliases (alias_heb, target_heb)
VALUES ('מעלות', 'מעלות-תרשיחא');