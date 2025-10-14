function showTab(tab) {
    document.querySelectorAll('.tab').forEach(t => t.classList.remove('active'));
    document.querySelectorAll('.section').forEach(s => s.classList.remove('active'));
    event.target.classList.add('active');
    document.getElementById(tab).classList.add('active');
    if (tab === 'cities') loadCities();
    if (tab === 'drivers') loadDrivers();
}

async function loadCities() {
    const res = await fetch(`${API_BASE_URL}/api/admin/cities`);
    const cities = await res.json();
    const tbody = document.querySelector('#citiesTable tbody');
    tbody.innerHTML = cities.map(c => `
        <tr class="${c.is_alias ? 'alias' : ''}">
            <td>${c.code}</td>
            <td>${c.name_heb}</td>
            <td>${c.name_eng || '-'}</td>
            <td>${c.is_alias ? 'כינוי → ' + c.canon_code : 'עיר'}</td>
            <td>
                <button class="delete-btn" onclick="deleteCity('${c.code}', ${c.is_alias})">מחק</button>
            </td>
        </tr>
    `).join('');
}

async function addCity() {
    const code = document.getElementById('cityCode').value;
    const nameHeb = document.getElementById('cityNameHeb').value;
    const nameEng = document.getElementById('cityNameEng').value;
    await fetch(`${API_BASE_URL}/api/admin/cities`, {
        method: 'POST',
        headers: {'Content-Type': 'application/json'},
        body: JSON.stringify({code, name_heb: nameHeb, name_eng: nameEng})
    });
    document.getElementById('cityCode').value = '';
    document.getElementById('cityNameHeb').value = '';
    document.getElementById('cityNameEng').value = '';
    loadCities();
}

async function addAlias() {
    const aliasHeb = document.getElementById('aliasHeb').value;
    const cityCode = document.getElementById('aliasCityCode').value;
    await fetch(`${API_BASE_URL}/api/admin/cities/alias`, {
        method: 'POST',
        headers: {'Content-Type': 'application/json'},
        body: JSON.stringify({alias_heb: aliasHeb, city_code: cityCode})
    });
    document.getElementById('aliasHeb').value = '';
    document.getElementById('aliasCityCode').value = '';
    loadCities();
}

async function deleteCity(code, isAlias) {
    if (!confirm('האם למחוק?')) return;
    const url = isAlias ? `${API_BASE_URL}/api/admin/cities/alias?alias=${code}` : `${API_BASE_URL}/api/admin/cities?code=${code}`;
    await fetch(url, {method: 'DELETE'});
    loadCities();
}

async function loadDrivers() {
    const res = await fetch(`${API_BASE_URL}/api/admin/drivers`);
    const drivers = await res.json();
    const tbody = document.querySelector('#driversTable tbody');
    tbody.innerHTML = drivers.map(d => `
        <tr>
            <td>${d.name}</td>
            <td>${d.car_number || '-'}</td>
            <td>${d.phone || '-'}</td>
            <td>${d.cities || '-'}</td>
            <td>
                <button class="delete-btn" onclick="deleteDriver(${d.id})">מחק</button>
            </td>
        </tr>
    `).join('');
}

async function importCities() {
    const file = document.getElementById('citiesFile').files[0];
    if (!file) return alert('בחר קובץ');
    const formData = new FormData();
    formData.append('file', file);
    const res = await fetch(`${API_BASE_URL}/api/admin/cities/import`, {method: 'POST', body: formData});
    const result = await res.json();
    alert(`נוספו: ${result.added}, עודכנו: ${result.updated}, דולגו: ${result.skipped}`);
    loadCities();
}

async function importDrivers() {
    const file = document.getElementById('driversFile').files[0];
    if (!file) return alert('בחר קובץ');
    const formData = new FormData();
    formData.append('file', file);
    const res = await fetch(`${API_BASE_URL}/api/admin/drivers/import`, {method: 'POST', body: formData});
    const result = await res.json();
    alert(`נוספו: ${result.added}, עודכנו: ${result.updated}, דולגו: ${result.skipped}`);
    loadDrivers();
}

function downloadTemplate() {
    window.location.href = `${API_BASE_URL}/api/admin/drivers/template`;
}

async function addDriver() {
    const name = document.getElementById('driverName').value;
    const carNumber = document.getElementById('driverCarNumber').value;
    const phone = document.getElementById('driverPhone').value;
    const cities = document.getElementById('driverCities').value;
    await fetch(`${API_BASE_URL}/api/admin/drivers`, {
        method: 'POST',
        headers: {'Content-Type': 'application/json'},
        body: JSON.stringify({name, phone, car_number: carNumber, cities})
    });
    document.getElementById('driverName').value = '';
    document.getElementById('driverCarNumber').value = '';
    document.getElementById('driverPhone').value = '';
    document.getElementById('driverCities').value = '';
    loadDrivers();
}

async function deleteDriver(id) {
    if (!confirm('האם למחוק?')) return;
    await fetch(`${API_BASE_URL}/api/admin/drivers?id=${id}`, {method: 'DELETE'});
    loadDrivers();
}

loadCities();
