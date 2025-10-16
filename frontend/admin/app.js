// Tab switching
function showTab(tab) {
    document.querySelectorAll('.section').forEach(s => s.classList.remove('active'));
    document.querySelectorAll('.tab').forEach(t => t.classList.remove('active'));
    document.getElementById(tab).classList.add('active');
    event.target.classList.add('active');
    
    if (tab === 'cities') loadCities();
    if (tab === 'drivers') loadDrivers();
}

// Cities Management
async function loadCities() {
    try {
        const res = await fetch(`${API_BASE_URL}/api/admin/cities`);
        const cities = await res.json();
        
        const tbody = document.querySelector('#citiesTable tbody');
        tbody.innerHTML = cities.map(city => `
            <tr>
                <td>${city.code}</td>
                <td>${city.name_heb}</td>
                <td>${city.name_eng || '-'}</td>
                <td>${city.is_alias ? 'כינוי' : 'עיר'}</td>
                <td>
                    <button onclick="deleteCity('${city.code}', ${city.is_alias})">מחק</button>
                </td>
            </tr>
        `).join('');
    } catch (err) {
        alert('שגיאה בטעינת ערים: ' + err.message);
    }
}

async function addCity() {
    const code = document.getElementById('cityCode').value;
    const nameHeb = document.getElementById('cityNameHeb').value;
    const nameEng = document.getElementById('cityNameEng').value;
    
    if (!code || !nameHeb) {
        alert('נא למלא קוד ושם בעברית');
        return;
    }
    
    try {
        const res = await fetch(`${API_BASE_URL}/api/admin/cities`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ code, name_heb: nameHeb, name_eng: nameEng })
        });
        
        if (res.ok) {
            alert('עיר נוספה בהצלחה');
            document.getElementById('cityCode').value = '';
            document.getElementById('cityNameHeb').value = '';
            document.getElementById('cityNameEng').value = '';
            loadCities();
        } else {
            const data = await res.json();
            alert('שגיאה: ' + data.error);
        }
    } catch (err) {
        alert('שגיאה: ' + err.message);
    }
}

async function deleteCity(code, isAlias) {
    if (!confirm(`האם למחוק ${isAlias ? 'כינוי' : 'עיר'} ${code}?`)) return;
    
    try {
        const res = await fetch(`${API_BASE_URL}/api/admin/cities?code=${code}`, {
            method: 'DELETE'
        });
        
        if (res.ok) {
            alert('נמחק בהצלחה');
            loadCities();
        } else {
            const data = await res.json();
            alert('שגיאה: ' + data.error);
        }
    } catch (err) {
        alert('שגיאה: ' + err.message);
    }
}

// Drivers Management
async function loadDrivers() {
    try {
        const res = await fetch(`${API_BASE_URL}/api/admin/drivers`);
        const drivers = await res.json();
        
        const tbody = document.querySelector('#driversTable tbody');
        tbody.innerHTML = drivers.map(driver => `
            <tr>
                <td>${driver.name}</td>
                <td>${driver.car_number || '-'}</td>
                <td>${driver.phone}</td>
                <td>${driver.cities || '-'}</td>
                <td>
                    <button onclick="editDriver(${driver.id}, '${driver.name}', '${driver.phone}', '${driver.car_number}', '${driver.cities}')">ערוך</button>
                    <button onclick="deleteDriver(${driver.id})">מחק</button>
                </td>
            </tr>
        `).join('');
    } catch (err) {
        alert('שגיאה בטעינת נהגים: ' + err.message);
    }
}

async function addDriver() {
    const name = document.getElementById('driverName').value;
    const carNumber = document.getElementById('driverCarNumber').value;
    const phone = document.getElementById('driverPhone').value;
    const cities = document.getElementById('driverCities').value;
    
    if (!name || !phone) {
        alert('נא למלא שם וטלפון');
        return;
    }
    
    try {
        const res = await fetch(`${API_BASE_URL}/api/admin/drivers`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ name, phone, car_number: carNumber, cities })
        });
        
        if (res.ok) {
            alert('נהג נוסף בהצלחה');
            document.getElementById('driverName').value = '';
            document.getElementById('driverCarNumber').value = '';
            document.getElementById('driverPhone').value = '';
            document.getElementById('driverCities').value = '';
            loadDrivers();
        } else {
            const data = await res.json();
            alert('שגיאה: ' + data.error);
        }
    } catch (err) {
        alert('שגיאה: ' + err.message);
    }
}

function editDriver(id, name, phone, carNumber, cities) {
    document.getElementById('driverName').value = name;
    document.getElementById('driverPhone').value = phone;
    document.getElementById('driverCarNumber').value = carNumber || '';
    document.getElementById('driverCities').value = cities || '';
    
    const btn = document.querySelector('#drivers .form button:last-child');
    btn.textContent = 'עדכן נהג';
    btn.onclick = () => updateDriver(id);
}

async function updateDriver(id) {
    const name = document.getElementById('driverName').value;
    const carNumber = document.getElementById('driverCarNumber').value;
    const phone = document.getElementById('driverPhone').value;
    const cities = document.getElementById('driverCities').value;
    
    try {
        const res = await fetch(`${API_BASE_URL}/api/admin/drivers`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ id, name, phone, car_number: carNumber, cities })
        });
        
        if (res.ok) {
            alert('נהג עודכן בהצלחה');
            document.getElementById('driverName').value = '';
            document.getElementById('driverCarNumber').value = '';
            document.getElementById('driverPhone').value = '';
            document.getElementById('driverCities').value = '';
            
            const btn = document.querySelector('#drivers .form button:last-child');
            btn.textContent = 'הוסף נהג';
            btn.onclick = addDriver;
            
            loadDrivers();
        } else {
            const data = await res.json();
            alert('שגיאה: ' + data.error);
        }
    } catch (err) {
        alert('שגיאה: ' + err.message);
    }
}

async function deleteDriver(id) {
    if (!confirm('האם למחוק נהג זה?')) return;
    
    try {
        const res = await fetch(`${API_BASE_URL}/api/admin/drivers?id=${id}`, {
            method: 'DELETE'
        });
        
        if (res.ok) {
            alert('נהג נמחק בהצלחה');
            loadDrivers();
        } else {
            const data = await res.json();
            alert('שגיאה: ' + data.error);
        }
    } catch (err) {
        alert('שגיאה: ' + err.message);
    }
}

// Initialize
loadCities();
