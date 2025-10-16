const uploadArea = document.getElementById('uploadArea');
const fileInput = document.getElementById('fileInput');
const fileInfo = document.getElementById('fileInfo');
const fileName = document.getElementById('fileName');
const processBtn = document.getElementById('processBtn');
const spinner = document.getElementById('spinner');
const log = document.getElementById('log');
const result = document.getElementById('result');
const preview = document.getElementById('preview');
const progressContainer = document.getElementById('progressContainer');
const progressFill = document.getElementById('progressFill');
const progressText = document.getElementById('progressText');
const history = document.getElementById('history');
const historyList = document.getElementById('historyList');
const toastContainer = document.getElementById('toastContainer');

let uploadedFile = null;
let currentFileData = null;

// Initialize
loadHistory();

// Upload Area Events
uploadArea.onclick = () => fileInput.click();

uploadArea.ondragover = (e) => {
    e.preventDefault();
    uploadArea.classList.add('dragover');
};

uploadArea.ondragleave = () => uploadArea.classList.remove('dragover');

uploadArea.ondrop = (e) => {
    e.preventDefault();
    uploadArea.classList.remove('dragover');
    if (e.dataTransfer.files.length) {
        fileInput.files = e.dataTransfer.files;
        handleFileSelect();
    }
};

fileInput.onchange = handleFileSelect;

function handleFileSelect() {
    const file = fileInput.files[0];
    if (file) {
        uploadedFile = null;
        processBtn.disabled = true;
        fileName.textContent = file.name;
        fileInfo.style.display = 'block';
        preview.style.display = 'none';
        uploadFile(file);
    }
}

async function uploadFile(file) {
    showToast('📤 מעלה קובץ...', 'info');
    addLog('📤 מעלה קובץ...');
    
    const formData = new FormData();
    formData.append('file', file);

    try {
        const res = await fetch(`${API_BASE_URL}/api/upload`, { method: 'POST', body: formData });
        const data = await res.json();
        
        if (res.ok) {
            uploadedFile = data.fullPath || data.path;
            currentFileData = { name: file.name, size: file.size };
            showToast('✅ קובץ הועלה בהצלחה', 'success');
            addLog('✅ קובץ הועלה בהצלחה');
            processBtn.disabled = false;
            
            // Show preview (simulated - in real app would parse Excel)
            showPreview(file);
        } else {
            showToast('❌ שגיאה בהעלאת קובץ', 'error');
            addLog('❌ שגיאה בהעלאת קובץ: ' + (data.error || 'Unknown error'), true);
        }
    } catch (err) {
        showToast('❌ שגיאה: ' + err.message, 'error');
        addLog('❌ שגיאה: ' + err.message, true);
    }
}

async function showPreview(file) {
    preview.innerHTML = `
        <h3>👁️ תצוגה מקדימה</h3>
        <div class="preview-stats">
            <div class="preview-stat"><strong>שם קובץ:</strong> ${file.name}</div>
            <div class="preview-stat"><strong>גודל:</strong> ${formatFileSize(file.size)}</div>
            <div class="preview-stat"><strong>סוג:</strong> ${file.type || 'Excel'}</div>
        </div>
        <p style="color: #6c757d; font-size: 14px;">✅ הקובץ מוכן לעיבוד</p>
    `;
    preview.style.display = 'block';
}

processBtn.onclick = async () => {
    if (!uploadedFile) return;

    processBtn.disabled = true;
    result.style.display = 'none';
    
    // Show progress
    progressContainer.style.display = 'block';
    updateProgress(0, 'מתחיל עיבוד...');
    
    showToast('⚙️ מתחיל עיבוד...', 'info');
    addLog('⚙️ מתחיל עיבוד...');

    const outputFile = `moh_${Date.now()}.xlsx`;

    try {
        // Simulate progress
        updateProgress(20, 'קורא קובץ...');
        await sleep(300);
        
        updateProgress(40, 'מעבד נתונים...');
        
        const res = await fetch(`${API_BASE_URL}/api/process`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ inputFile: uploadedFile, outputFile })
        });

        updateProgress(80, 'משלים עיבוד...');
        const data = await res.json();
        
        updateProgress(100, 'הושלם!');
        await sleep(500);
        progressContainer.style.display = 'none';

        if (data.success) {
            showToast('✅ עיבוד הושלם בהצלחה!', 'success');
            addLog('✅ עיבוד הושלם בהצלחה');
            addLog(`📊 שורות קלט: ${data.inputRows}`);
            addLog(`📊 שורות פלט: ${data.outputRows}`);
            addLog(`⏱️ זמן עיבוד: ${data.processTime}`);
            
            result.className = 'result';
            result.innerHTML = `
                <h3>✅ הקובץ עובד בהצלחה!</h3>
                <p><strong>שורות בקובץ המקור:</strong> ${data.inputRows}</p>
                <p><strong>שורות בקובץ הסופי:</strong> ${data.outputRows}</p>
                <p><strong>זמן עיבוד:</strong> ${data.processTime}</p>
                <button class="btn" onclick="downloadFile('${data.outputFile}')">⬇️ הורד קובץ</button>
            `;
            result.style.display = 'block';
            
            // Add to history
            addToHistory({
                fileName: currentFileData.name,
                outputFile: data.outputFile,
                inputRows: data.inputRows,
                outputRows: data.outputRows,
                timestamp: Date.now()
            });
        } else {
            showToast('❌ שגיאה בעיבוד', 'error');
            addLog('❌ שגיאה: ' + data.message, true);
            result.className = 'result error';
            result.innerHTML = `<h3>❌ שגיאה בעיבוד</h3><p>${data.message}</p>`;
            result.style.display = 'block';
        }
    } catch (err) {
        progressContainer.style.display = 'none';
        showToast('❌ שגיאה: ' + err.message, 'error');
        addLog('❌ שגיאה: ' + err.message, true);
    }

    processBtn.disabled = false;
};

// Toast Notifications
function showToast(message, type = 'info') {
    const toast = document.createElement('div');
    toast.className = `toast ${type}`;
    
    const icon = type === 'success' ? '✅' : type === 'error' ? '❌' : 'ℹ️';
    toast.innerHTML = `<span style="font-size: 20px;">${icon}</span><span>${message}</span>`;
    
    toastContainer.appendChild(toast);
    
    setTimeout(() => {
        toast.style.animation = 'slideDown 0.3s reverse';
        setTimeout(() => toast.remove(), 300);
    }, 3000);
}

// Progress Bar
function updateProgress(percent, text) {
    progressFill.style.width = percent + '%';
    progressFill.textContent = percent + '%';
    progressText.textContent = text;
}

// History Management
function addToHistory(item) {
    let historyData = JSON.parse(localStorage.getItem('excelFlowHistory') || '[]');
    historyData.unshift(item);
    historyData = historyData.slice(0, 10); // Keep last 10
    localStorage.setItem('excelFlowHistory', JSON.stringify(historyData));
    loadHistory();
}

function loadHistory() {
    const historyData = JSON.parse(localStorage.getItem('excelFlowHistory') || '[]');
    
    if (historyData.length === 0) {
        history.style.display = 'none';
        return;
    }
    
    history.style.display = 'block';
    historyList.innerHTML = historyData.map(item => `
        <div class="history-item">
            <div class="info">
                <div><strong>${item.fileName}</strong></div>
                <div class="time">${new Date(item.timestamp).toLocaleString('he-IL')}</div>
                <div style="font-size: 12px; color: #6c757d;">
                    קלט: ${item.inputRows} | פלט: ${item.outputRows}
                </div>
            </div>
            <button class="download-btn" onclick="downloadFile('${item.outputFile}')">⬇️ הורד</button>
        </div>
    `).join('');
}

function clearHistory() {
    if (confirm('האם למחוק את כל ההיסטוריה?')) {
        localStorage.removeItem('excelFlowHistory');
        loadHistory();
        showToast('🗑️ ההיסטוריה נמחקה', 'info');
    }
}

// Logging
function addLog(message, isError = false) {
    log.style.display = 'block';
    const line = document.createElement('div');
    line.className = 'log-line';
    line.style.color = isError ? '#ff4444' : '#00ff00';
    line.textContent = `[${new Date().toLocaleTimeString('he-IL')}] ${message}`;
    log.appendChild(line);
    log.scrollTop = log.scrollHeight;
}

// Download
function downloadFile(filename) {
    window.location.href = `${API_BASE_URL}/api/download/${filename}`;
    showToast('⬇️ מוריד קובץ...', 'info');
}

// Utilities
function formatFileSize(bytes) {
    if (bytes < 1024) return bytes + ' B';
    if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB';
    return (bytes / (1024 * 1024)).toFixed(1) + ' MB';
}

function sleep(ms) {
    return new Promise(resolve => setTimeout(resolve, ms));
}
