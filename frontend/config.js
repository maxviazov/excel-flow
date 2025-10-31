// Configuration for fish.viazov.dev
const CONFIG = {
    API_BASE_URL: 'https://excel.viazov.dev/api',
    UPLOAD_ENDPOINT: '/upload',
    PROCESS_ENDPOINT: '/process',
    DOWNLOAD_ENDPOINT: '/download',
    VALIDATE_ENDPOINT: '/validate',
    ADMIN_CITIES_ENDPOINT: '/admin/cities',
    ADMIN_DRIVERS_ENDPOINT: '/admin/drivers',
    
    // UI Configuration
    MAX_FILE_SIZE: 50 * 1024 * 1024, // 50MB
    ALLOWED_FILE_TYPES: ['.xlsx', '.xls'],
    
    // Localization
    DEFAULT_LANGUAGE: 'he',
    RTL_SUPPORT: true
};