// Admin i18n
const translations = {
    he: {
        title: 'ניהול מאגרי מידע',
        heading: 'ניהול מאגרי מידע',
        citiesTab: 'ערים',
        driversTab: 'נהגים',
        citiesManagement: 'ניהול ערים',
        driversManagement: 'ניהול נהגים',
        citySearch: '🔍 חיפוש לפי קוד או שם...',
        driverSearch: '🔍 חיפוש לפי שם או טלפון...',
        importFromExcel: 'ייבא מקובץ Excel',
        downloadTemplate: 'הורד תבנית Excel',
        addCity: 'הוסף עיר',
        addDriver: 'הוסף נהג',
        addAlias: 'הוסף כינוי',
        cityCode: 'קוד עיר',
        cityNameHeb: 'שם בעברית',
        cityNameEng: 'שם באנגלית',
        aliasHeb: 'כינוי בעברית',
        driverName: 'שם נהג',
        driverCarNumber: 'מספר רכב',
        driverPhone: 'טלפון',
        driverCityCodes: 'קודי ערים (מופרד בפסיק)',
        driverCityNames: 'שמות ערים (מופרד בפסיק)',
        code: 'קוד',
        hebrewName: 'שם עברי',
        englishName: 'שם אנגלי',
        type: 'סוג',
        actions: 'פעולות',
        name: 'שם',
        carNumber: 'מספר רכב',
        phone: 'טלפון',
        cityCodes: 'קודי ערים',
        cityNames: 'שמות ערים',
        delete: 'מחק'
    },
    ru: {
        title: 'Управление базами данных',
        heading: 'Управление базами данных',
        citiesTab: 'Города',
        driversTab: 'Водители',
        citiesManagement: 'Управление городами',
        driversManagement: 'Управление водителями',
        citySearch: '🔍 Поиск по коду или названию...',
        driverSearch: '🔍 Поиск по имени или телефону...',
        importFromExcel: 'Импорт из Excel',
        downloadTemplate: 'Скачать шаблон Excel',
        addCity: 'Добавить город',
        addDriver: 'Добавить водителя',
        addAlias: 'Добавить псевдоним',
        cityCode: 'Код города',
        cityNameHeb: 'Название на иврите',
        cityNameEng: 'Название на английском',
        aliasHeb: 'Псевдоним на иврите',
        driverName: 'Имя водителя',
        driverCarNumber: 'Номер машины',
        driverPhone: 'Телефон',
        driverCityCodes: 'Коды городов (через запятую)',
        driverCityNames: 'Названия городов (через запятую)',
        code: 'Код',
        hebrewName: 'Название на иврите',
        englishName: 'Название на английском',
        type: 'Тип',
        actions: 'Действия',
        name: 'Имя',
        carNumber: 'Номер машины',
        phone: 'Телефон',
        cityCodes: 'Коды городов',
        cityNames: 'Названия городов',
        delete: 'Удалить'
    },
    en: {
        title: 'Database Management',
        heading: 'Database Management',
        citiesTab: 'Cities',
        driversTab: 'Drivers',
        citiesManagement: 'Cities Management',
        driversManagement: 'Drivers Management',
        citySearch: '🔍 Search by code or name...',
        driverSearch: '🔍 Search by name or phone...',
        importFromExcel: 'Import from Excel',
        downloadTemplate: 'Download Excel Template',
        addCity: 'Add City',
        addDriver: 'Add Driver',
        addAlias: 'Add Alias',
        cityCode: 'City Code',
        cityNameHeb: 'Name in Hebrew',
        cityNameEng: 'Name in English',
        aliasHeb: 'Alias in Hebrew',
        driverName: 'Driver Name',
        driverCarNumber: 'Car Number',
        driverPhone: 'Phone',
        driverCityCodes: 'City Codes (comma separated)',
        driverCityNames: 'City Names (comma separated)',
        code: 'Code',
        hebrewName: 'Hebrew Name',
        englishName: 'English Name',
        type: 'Type',
        actions: 'Actions',
        name: 'Name',
        carNumber: 'Car Number',
        phone: 'Phone',
        cityCodes: 'City Codes',
        cityNames: 'City Names',
        delete: 'Delete'
    }
};

function detectBrowserLanguage() {
    const browserLang = navigator.language || navigator.userLanguage;
    if (browserLang.startsWith('he')) return 'he';
    if (browserLang.startsWith('ru')) return 'ru';
    if (browserLang.startsWith('en')) return 'en';
    return 'he';
}

let currentLang = localStorage.getItem('excelFlowLang') || detectBrowserLanguage();

function t(key) {
    return translations[currentLang][key] || key;
}

function setLanguage(lang) {
    if (!translations[lang]) return;
    currentLang = lang;
    localStorage.setItem('excelFlowLang', lang);
    document.documentElement.dir = lang === 'he' ? 'rtl' : 'ltr';
    document.documentElement.lang = lang;
    updateTranslations();
    updateLanguageButtons();
}

function updateTranslations() {
    document.title = t('title');
    document.querySelectorAll('[data-i18n]').forEach(el => {
        const key = el.getAttribute('data-i18n');
        if (el.tagName === 'INPUT' || el.tagName === 'TEXTAREA') {
            el.placeholder = t(key);
        } else if (el.tagName === 'BUTTON') {
            const icon = el.textContent.match(/^[^\w\s]+/)?.[0] || '';
            el.textContent = icon + (icon ? ' ' : '') + t(key);
        } else {
            el.textContent = t(key);
        }
    });
}

document.addEventListener('DOMContentLoaded', () => {
    setLanguage(currentLang);
    updateLanguageButtons();
});

function updateLanguageButtons() {
    document.querySelectorAll('[data-lang]').forEach(btn => {
        const lang = btn.getAttribute('data-lang');
        if (lang === currentLang) {
            btn.style.background = '#667eea';
            btn.style.color = 'white';
        } else {
            btn.style.background = 'white';
            btn.style.color = '#667eea';
        }
    });
}
