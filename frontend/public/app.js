const uploadArea = document.getElementById('uploadArea');
const fileInput = document.getElementById('fileInput');
const fileInfo = document.getElementById('fileInfo');
const fileName = document.getElementById('fileName');
const processBtn = document.getElementById('processBtn');
const spinner = document.getElementById('spinner');
const log = document.getElementById('log');
const result = document.getElementById('result');

let uploadedFile = null;

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
        uploadFile(file);
    }
}

async function uploadFile(file) {
    addLog('📤 מעלה קובץ...');
    const formData = new FormData();
    formData.append('file', file);

    try {
        const res = await fetch(`${API_BASE_URL}/api/upload`, { method: 'POST', body: formData });
        const data = await res.json();
        
        if (res.ok) {
            uploadedFile = data.path;
            addLog('✅ קובץ הועלה בהצלחה');
            processBtn.disabled = false;
        } else {
            addLog('❌ שגיאה בהעלאת קובץ: ' + data.error, true);
        }
    } catch (err) {
        addLog('❌ שגיאה: ' + err.message, true);
    }
}

processBtn.onclick = async () => {
    if (!uploadedFile) return;

    processBtn.disabled = true;
    spinner.style.display = 'block';
    result.style.display = 'none';
    addLog('⚙️ מתחיל עיבוד...');

    const outputFile = `moh_${Date.now()}.xlsx`;

    try {
        const res = await fetch(`${API_BASE_URL}/api/process`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ inputFile: uploadedFile, outputFile })
        });

        const data = await res.json();
        spinner.style.display = 'none';

        if (data.success) {
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
        } else {
            addLog('❌ שגיאה: ' + data.message, true);
            result.className = 'result error';
            result.innerHTML = `<h3>❌ שגיאה בעיבוד</h3><p>${data.message}</p>`;
            result.style.display = 'block';
        }
    } catch (err) {
        spinner.style.display = 'none';
        addLog('❌ שגיאה: ' + err.message, true);
    }

    processBtn.disabled = false;
};

function addLog(message, isError = false) {
    log.style.display = 'block';
    const line = document.createElement('div');
    line.className = 'log-line';
    line.style.color = isError ? '#ff4444' : '#00ff00';
    line.textContent = `[${new Date().toLocaleTimeString('he-IL')}] ${message}`;
    log.appendChild(line);
    log.scrollTop = log.scrollHeight;
}

function downloadFile(filename) {
    window.location.href = `${API_BASE_URL}/api/download/${filename}`;
}
