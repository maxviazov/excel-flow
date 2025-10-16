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
let batchFiles = [];
const MAX_BATCH_SIZE = 5;

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
    const files = Array.from(fileInput.files);
    if (files.length === 0) return;
    
    if (files.length === 1) {
        // Single file mode
        uploadedFile = null;
        processBtn.disabled = true;
        fileName.textContent = files[0].name;
        fileInfo.style.display = 'block';
        preview.style.display = 'none';
        batchQueue.style.display = 'none';
        uploadFile(files[0]);
    } else {
        // Batch mode
        if (files.length > MAX_BATCH_SIZE) {
            showToast(`מקסימום ${MAX_BATCH_SIZE} קבצים בבת אחת`, 'error');
            return;
        }
        batchFiles = files.map(f => ({ file: f, status: 'pending', uploadedPath: null, result: null }));
        fileName.textContent = `${files.length} קבצים נבחרו`;
        fileInfo.style.display = 'block';
        preview.style.display = 'none';
        renderBatchQueue();
        processBtn.disabled = false;
        processBtn.textContent = `▶️ עבד ${files.length} קבצים`;
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
    if (batchFiles.length > 0) {
        await processBatch();
        return;
    }
    
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
                <div style="display: flex; gap: 10px;">
                    <button class="btn" onclick="downloadFile('${data.outputFile}')" style="flex: 1;">⬇️ הורד Excel</button>
                    <button class="btn btn-secondary" onclick="downloadCSV('${data.outputFile}')" style="flex: 1;">📄 הורד CSV</button>
                </div>
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

function downloadCSV(filename) {
    window.location.href = `${API_BASE_URL}/api/export-csv/${filename}`;
    showToast('📄 מייצא ל-CSV...', 'info');
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

// Batch Processing
function renderBatchQueue() {
    batchQueue.style.display = 'block';
    batchQueue.innerHTML = `
        <h3>📦 תור עיבוד (${batchFiles.length} קבצים)</h3>
        ${batchFiles.map((item, idx) => `
            <div class="batch-item" id="batch-${idx}">
                <div class="name">${item.file.name}</div>
                <div class="status ${item.status}">${getStatusText(item.status)}</div>
                ${item.status === 'pending' ? `<button class="remove-btn" onclick="removeBatchItem(${idx})">✕</button>` : ''}
            </div>
        `).join('')}
    `;
}

function getStatusText(status) {
    const texts = {
        pending: '⏳ ממתין',
        processing: '⚙️ מעבד',
        success: '✅ הושלם',
        error: '❌ שגיאה'
    };
    return texts[status] || status;
}

function removeBatchItem(idx) {
    batchFiles.splice(idx, 1);
    if (batchFiles.length === 0) {
        batchQueue.style.display = 'none';
        processBtn.disabled = true;
    } else {
        renderBatchQueue();
        processBtn.textContent = `▶️ עבד ${batchFiles.length} קבצים`;
    }
}

async function processBatch() {
    processBtn.disabled = true;
    showToast(`מתחיל עיבוד ${batchFiles.length} קבצים...`, 'info');
    
    for (let i = 0; i < batchFiles.length; i++) {
        const item = batchFiles[i];
        item.status = 'processing';
        renderBatchQueue();
        
        try {
            // Upload
            const formData = new FormData();
            formData.append('file', item.file);
            const uploadRes = await fetch(`${API_BASE_URL}/api/upload`, { method: 'POST', body: formData });
            const uploadData = await uploadRes.json();
            
            if (!uploadRes.ok) throw new Error('Upload failed');
            
            item.uploadedPath = uploadData.fullPath || uploadData.path;
            
            // Validate
            const validateRes = await fetch(`${API_BASE_URL}/api/validate`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ filePath: item.uploadedPath })
            });
            const validateData = await validateRes.json();
            
            if (!validateData.valid) {
                throw new Error(validateData.error || 'Validation failed');
            }
            
            // Process
            const outputFile = `moh_${Date.now()}_${i}.xlsx`;
            const processRes = await fetch(`${API_BASE_URL}/api/process`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ inputFile: item.uploadedPath, outputFile })
            });
            
            const processData = await processRes.json();
            
            if (processData.success) {
                item.status = 'success';
                item.result = processData;
                addToHistory({
                    fileName: item.file.name,
                    outputFile: processData.outputFile,
                    inputRows: processData.inputRows,
                    outputRows: processData.outputRows,
                    timestamp: Date.now()
                });
            } else {
                item.status = 'error';
            }
        } catch (err) {
            item.status = 'error';
            showToast(`שגיאה בעיבוד ${item.file.name}`, 'error');
        }
        
        renderBatchQueue();
    }
    
    const successCount = batchFiles.filter(f => f.status === 'success').length;
    showToast(`הושלם! ${successCount}/${batchFiles.length} קבצים עובדו בהצלחה`, 'success');
    
    // Show download all button
    if (successCount > 0) {
        const downloadAllBtn = document.createElement('div');
        downloadAllBtn.className = 'batch-actions';
        downloadAllBtn.innerHTML = `
            <button class="btn download-all-btn" onclick="downloadAllBatch()">⬇️ הורד את כל הקבצים (${successCount})</button>
        `;
        batchQueue.appendChild(downloadAllBtn);
    }
    
    processBtn.disabled = false;
    processBtn.textContent = '▶️ התחל עיבוד';
}

function downloadAllBatch() {
    const successFiles = batchFiles.filter(f => f.status === 'success');
    successFiles.forEach((item, idx) => {
        setTimeout(() => {
            downloadFile(item.result.outputFile);
        }, idx * 500);
    });
    showToast(`מוריד ${successFiles.length} קבצים...`, 'info');
}
