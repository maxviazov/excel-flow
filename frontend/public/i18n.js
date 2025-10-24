// i18n - Internationalization
const translations = {
    he: {
        title: 'Excel Flow - עיבוד קבצי Excel',
        heading: '🚀 Excel Flow',
        adminLink: '⚙️ ניהול מאגרי מידע',
        uploadTitle: 'בחר קבצי Excel לעיבוד',
        uploadDesc: 'גרור קבצים לכאן או לחץ לבחירה (עד 5 קבצים)',
        filesSelected: 'קבצים נבחרו:',
        startProcessing: '▶️ התחל עיבוד',
        processing: 'מעבד...',
        uploadingFile: '📤 מעלה קובץ...',
        fileUploaded: '✅ קובץ הועלה בהצלחה',
        uploadError: '❌ שגיאה בהעלאת קובץ',
        processingStarted: '⚙️ מתחיל עיבוד...',
        processingComplete: '✅ עיבוד הושלם בהצלחה!',
        processingError: '❌ שגיאה בעיבוד',
        downloadExcel: '⬇️ הורד Excel',
        downloadCSV: '📄 הורד CSV',
        history: '📜 היסטוריית עיבודים',
        clearHistory: '🗑️ נקה היסטוריה',
        confirmClear: 'האם למחוק את כל ההיסטוריה?',
        historyCleared: '🗑️ ההיסטוריה נמחקה',
        downloading: '⬇️ מוריד קובץ...',
        exportingCSV: '📄 מייצא ל-CSV...',
        preview: '👁️ תצוגה מקדימה',
        fileName: 'שם קובץ:',
        fileSize: 'גודל:',
        fileType: 'סוג:',
        fileReady: '✅ הקובץ מוכן לעיבוד',
        inputRows: 'שורות בקובץ המקור:',
        outputRows: 'שורות בקובץ הסופי:',
        processTime: 'זמן עיבוד:',
        batchQueue: '📦 תור עיבוד',
        files: 'קבצים',
        processFiles: 'עבד',
        downloadAll: '⬇️ הורד את כל הקבצים',
        maxFiles: 'מקסימום',
        filesAtOnce: 'קבצים בבת אחת',
        statusPending: '⏳ ממתין',
        statusProcessing: '⚙️ מעבד',
        statusSuccess: '✅ הושלם',
        statusError: '❌ שגיאה',
        readingFile: 'קורא קובץ...',
        processingData: 'מעבד נתונים...',
        completing: 'משלים עיבוד...',
        completed: 'הושלם!',
        language: 'שפה'
    },
    ru: {
        title: 'Excel Flow - Обработка Excel файлов',
        heading: '🚀 Excel Flow',
        adminLink: '⚙️ Управление базами данных',
        uploadTitle: 'Выберите Excel файлы для обработки',
        uploadDesc: 'Перетащите файлы сюда или нажмите для выбора (до 5 файлов)',
        filesSelected: 'Выбрано файлов:',
        startProcessing: '▶️ Начать обработку',
        processing: 'Обработка...',
        uploadingFile: '📤 Загрузка файла...',
        fileUploaded: '✅ Файл успешно загружен',
        uploadError: '❌ Ошибка загрузки файла',
        processingStarted: '⚙️ Начинаем обработку...',
        processingComplete: '✅ Обработка завершена успешно!',
        processingError: '❌ Ошибка обработки',
        downloadExcel: '⬇️ Скачать Excel',
        downloadCSV: '📄 Скачать CSV',
        history: '📜 История обработок',
        clearHistory: '🗑️ Очистить историю',
        confirmClear: 'Удалить всю историю?',
        historyCleared: '🗑️ История очищена',
        downloading: '⬇️ Скачивание файла...',
        exportingCSV: '📄 Экспорт в CSV...',
        preview: '👁️ Предпросмотр',
        fileName: 'Имя файла:',
        fileSize: 'Размер:',
        fileType: 'Тип:',
        fileReady: '✅ Файл готов к обработке',
        inputRows: 'Строк в исходном файле:',
        outputRows: 'Строк в итоговом файле:',
        processTime: 'Время обработки:',
        batchQueue: '📦 Очередь обработки',
        files: 'файлов',
        processFiles: 'Обработать',
        downloadAll: '⬇️ Скачать все файлы',
        maxFiles: 'Максимум',
        filesAtOnce: 'файлов за раз',
        statusPending: '⏳ Ожидание',
        statusProcessing: '⚙️ Обработка',
        statusSuccess: '✅ Завершено',
        statusError: '❌ Ошибка',
        readingFile: 'Чтение файла...',
        processingData: 'Обработка данных...',
        completing: 'Завершение обработки...',
        completed: 'Завершено!',
        language: 'Язык'
    },
    en: {
        title: 'Excel Flow - Excel File Processing',
        heading: '🚀 Excel Flow',
        adminLink: '⚙️ Database Management',
        uploadTitle: 'Select Excel files for processing',
        uploadDesc: 'Drag files here or click to select (up to 5 files)',
        filesSelected: 'Files selected:',
        startProcessing: '▶️ Start Processing',
        processing: 'Processing...',
        uploadingFile: '📤 Uploading file...',
        fileUploaded: '✅ File uploaded successfully',
        uploadError: '❌ File upload error',
        processingStarted: '⚙️ Starting processing...',
        processingComplete: '✅ Processing completed successfully!',
        processingError: '❌ Processing error',
        downloadExcel: '⬇️ Download Excel',
        downloadCSV: '📄 Download CSV',
        history: '📜 Processing History',
        clearHistory: '🗑️ Clear History',
        confirmClear: 'Delete all history?',
        historyCleared: '🗑️ History cleared',
        downloading: '⬇️ Downloading file...',
        exportingCSV: '📄 Exporting to CSV...',
        preview: '👁️ Preview',
        fileName: 'File name:',
        fileSize: 'Size:',
        fileType: 'Type:',
        fileReady: '✅ File ready for processing',
        inputRows: 'Rows in source file:',
        outputRows: 'Rows in output file:',
        processTime: 'Processing time:',
        batchQueue: '📦 Processing Queue',
        files: 'files',
        processFiles: 'Process',
        downloadAll: '⬇️ Download all files',
        maxFiles: 'Maximum',
        filesAtOnce: 'files at once',
        statusPending: '⏳ Pending',
        statusProcessing: '⚙️ Processing',
        statusSuccess: '✅ Completed',
        statusError: '❌ Error',
        readingFile: 'Reading file...',
        processingData: 'Processing data...',
        completing: 'Completing processing...',
        completed: 'Completed!',
        language: 'Language'
    }
};

// Current language
let currentLang = localStorage.getItem('excelFlowLang') || 'he';

// Get translation
function t(key) {
    return translations[currentLang][key] || key;
}

// Set language
function setLanguage(lang) {
    if (!translations[lang]) return;
    currentLang = lang;
    localStorage.setItem('excelFlowLang', lang);
    
    // Update direction
    document.documentElement.dir = lang === 'he' ? 'rtl' : 'ltr';
    document.documentElement.lang = lang;
    
    // Update all translatable elements
    updateTranslations();
}

// Update all translations on page
function updateTranslations() {
    document.title = t('title');
    
    // Update elements with data-i18n attribute
    document.querySelectorAll('[data-i18n]').forEach(el => {
        const key = el.getAttribute('data-i18n');
        if (el.tagName === 'INPUT' && el.type === 'button') {
            el.value = t(key);
        } else if (el.tagName === 'BUTTON') {
            // Preserve icons
            const icon = el.textContent.match(/^[^\w\s]+/)?.[0] || '';
            el.textContent = icon + (icon ? ' ' : '') + t(key);
        } else {
            el.textContent = t(key);
        }
    });
}

// Initialize language on load
document.addEventListener('DOMContentLoaded', () => {
    setLanguage(currentLang);
});
