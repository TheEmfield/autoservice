let currentUser = null;

async function login() {
    const username = document.getElementById('username').value;
    const password = document.getElementById('password').value;
    try {
        const res = await fetch('/api/login', {
            method: 'POST',
            headers: {'Content-Type': 'application/json'},
            body: JSON.stringify({username, password})
        });
        if (res.ok) {
            currentUser = await res.json();
            document.getElementById('login-view').classList.add('hidden');
            document.getElementById('app-view').classList.remove('hidden');
            document.getElementById('user-role').textContent = currentUser.role === 'admin' ? 'Администратор' : 'Оператор';
            
            const isAdmin = currentUser.role === 'admin';
            document.getElementById('admin-cars').classList.toggle('hidden', !isAdmin);
            document.getElementById('admin-masters').classList.toggle('hidden', !isAdmin);
            document.getElementById('admin-services').classList.toggle('hidden', !isAdmin);
            
            loadAll();
        } else {
            alert('Неверный логин или пароль');
        }
    } catch (e) { alert('Ошибка соединения'); }
}

function logout() {
    currentUser = null;
    document.getElementById('app-view').classList.add('hidden');
    document.getElementById('login-view').classList.remove('hidden');
}

function showTab(tabName) {
    document.querySelectorAll('.tab-content').forEach(el => el.classList.add('hidden'));
    document.getElementById('tab-' + tabName).classList.remove('hidden');
    document.querySelectorAll('.tab').forEach(el => el.classList.remove('active'));
    event.target.classList.add('active');
    if(tabName === 'works') loadSelects();
}

async function loadAll() {
    loadTable('/api/cars', 'cars-table', renderCarRow);
    loadTable('/api/masters', 'masters-table', renderMasterRow);
    loadTable('/api/services', 'services-table', renderServiceRow);
    loadTable('/api/works', 'works-table', renderWorkRow);
}

async function loadTable(url, tableId, renderFn) {
    const res = await fetch(url);
    const data = await res.json();
    const tbody = document.querySelector(`#${tableId} tbody`);
    tbody.innerHTML = data.map(renderFn).join('');
}

function renderCarRow(c) {
    const delBtn = currentUser.role === 'admin' ? `<button class="danger" onclick="deleteItem('/api/cars', ${c.id}, loadAll)">Удалить</button>` : '-';
    return `<tr><td>${c.num}</td><td>${c.mark}</td><td>${c.color}</td><td>${c.is_foreign ? 'Иномарка' : 'Отеч.'}</td><td>${delBtn}</td></tr>`;
}
async function addCar() {
    const body = {
        num: document.getElementById('car-num').value,
        mark: document.getElementById('car-mark').value,
        color: document.getElementById('car-color').value,
        is_foreign: parseInt(document.getElementById('car-foreign').value)
    };
    await postAndReload('/api/cars', body, loadAll);
}

function renderMasterRow(m) {
    const delBtn = currentUser.role === 'admin' ? `<button class="danger" onclick="deleteItem('/api/masters', ${m.id}, loadAll)">Удалить</button>` : '-';
    return `<tr><td>${m.id}</td><td>${m.name}</td><td>${delBtn}</td></tr>`;
}
async function addMaster() {
    await postAndReload('/api/masters', {name: document.getElementById('master-name').value}, loadAll);
}

function renderServiceRow(s) {
    const delBtn = currentUser.role === 'admin' ? `<button class="danger" onclick="deleteItem('/api/services', ${s.id}, loadAll)">Удалить</button>` : '-';
    return `<tr><td>${s.name}</td><td>${s.cost_our}</td><td>${s.cost_foreign}</td><td>${delBtn}</td></tr>`;
}
async function addService() {
    const body = {
        name: document.getElementById('serv-name').value,
        cost_our: parseFloat(document.getElementById('serv-cost-our').value),
        cost_foreign: parseFloat(document.getElementById('serv-cost-for').value)
    };
    await postAndReload('/api/services', body, loadAll);
}

function renderWorkRow(w) {
    const delBtn = currentUser.role === 'admin' ? `<button class="danger" onclick="deleteItem('/api/works', ${w.id}, loadAll)">Удалить</button>` : '-';
    return `<tr><td>${w.date_work}</td><td>${w.master_id}</td><td>${w.car_id}</td><td>${w.service_id}</td><td>${delBtn}</td></tr>`;
}
async function loadSelects() {
    const masters = await (await fetch('/api/masters')).json();
    const cars = await (await fetch('/api/cars')).json();
    const services = await (await fetch('/api/services')).json();
    
    fillSelect('work-master', masters, 'id', 'name');
    fillSelect('work-car', cars, 'id', 'num');
    fillSelect('work-service', services, 'id', 'name');
}
function fillSelect(id, data, valKey, textKey) {
    const sel = document.getElementById(id);
    sel.innerHTML = '<option value="">Выберите...</option>' + data.map(d => `<option value="${d[valKey]}">${d[textKey]}</option>`).join('');
}
async function addWork() {
    const body = {
        date_work: document.getElementById('work-date').value,
        master_id: parseInt(document.getElementById('work-master').value),
        car_id: parseInt(document.getElementById('work-car').value),
        service_id: parseInt(document.getElementById('work-service').value)
    };
    await postAndReload('/api/works', body, loadAll);
}

function downloadReport(type) {
    if (type === 'cost') {
        const start = document.getElementById('rep-cost-start').value;
        const end = document.getElementById('rep-cost-end').value;
        window.open(`/api/reports/cost?start=${start}&end=${end}`, '_blank');
    } else {
        const year = document.getElementById('rep-top-year').value;
        const month = document.getElementById('rep-top-month').value;
        window.open(`/api/reports/top-masters?year=${year}&month=${month}`, '_blank');
    }
}

async function postAndReload(url, body, reloadFn) {
    const res = await fetch(url, {method: 'POST', headers: {'Content-Type': 'application/json'}, body: JSON.stringify(body)});
    if (res.ok) { reloadFn(); } else { alert('Ошибка: ' + await res.text()); }
}
async function deleteItem(url, id, reloadFn) {
    if (!confirm('Удалить запись?')) return;
    const res = await fetch(`${url}?id=${id}`, {method: 'DELETE'});
    if (res.ok) { reloadFn(); } else { alert('Ошибка удаления (возможно, сработал триггер БД): ' + await res.text()); }
}