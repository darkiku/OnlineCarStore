const API_URL = 'http://localhost:9000/api';

// State
let allCars = [];
let garage = JSON.parse(localStorage.getItem('garage')) || [];

// Navigation
document.getElementById('catalog-btn').addEventListener('click', () => showPage('catalog-page'));
document.getElementById('garage-btn').addEventListener('click', () => showPage('garage-page'));
document.getElementById('add-car-btn').addEventListener('click', () => showPage('add-car-page'));

function showPage(pageId) {
    document.querySelectorAll('.page').forEach(page => page.classList.remove('active'));
    document.querySelectorAll('.nav-btn').forEach(btn => btn.classList.remove('active'));

    document.getElementById(pageId).classList.add('active');

    if (pageId === 'catalog-page') {
        document.getElementById('catalog-btn').classList.add('active');
    } else if (pageId === 'garage-page') {
        document.getElementById('garage-btn').classList.add('active');
        renderGarage();
    } else if (pageId === 'add-car-page') {
        document.getElementById('add-car-btn').classList.add('active');
    }
}

// Fetch cars
async function fetchCars(filters = {}) {
    try {
        let url = `${API_URL}/cars`;
        const params = new URLSearchParams();

        if (filters.make) params.append('make', filters.make);
        if (filters.body_type) params.append('body_type', filters.body_type);
        if (filters.fuel_type) params.append('fuel_type', filters.fuel_type);
        if (filters.transmission) params.append('transmission', filters.transmission);

        if (params.toString()) {
            url += '?' + params.toString();
        }

        const response = await fetch(url);
        allCars = await response.json();

        // Client-side price filtering
        if (filters.min_price || filters.max_price) {
            allCars = allCars.filter(car => {
                if (filters.min_price && car.price < parseFloat(filters.min_price)) return false;
                if (filters.max_price && car.price > parseFloat(filters.max_price)) return false;
                return true;
            });
        }

        renderCars();
    } catch (error) {
        console.error('Error fetching cars:', error);
        document.getElementById('cars-list').innerHTML = '<p>Ошибка загрузки автомобилей</p>';
    }
}

// Render cars
function renderCars() {
    const carsList = document.getElementById('cars-list');

    if (allCars.length === 0) {
        carsList.innerHTML = '<p>Автомобили не найдены</p>';
        return;
    }

    carsList.innerHTML = allCars.map(car => `
        <div class="car-card">
            <img src="${car.image_url || 'https://via.placeholder.com/400x300?text=' + car.make + '+' + car.model}" 
                 alt="${car.make} ${car.model}" class="car-image">
            <div class="car-info">
                <div class="car-title">${car.make} ${car.model}</div>
                <div class="car-specs">
                    <div class="spec-item">📅 ${car.year} год</div>
                    <div class="spec-item">⚡ ${car.horsepower || 'N/A'} л.с.</div>
                    <div class="spec-item">🔧 ${car.transmission}</div>
                    <div class="spec-item">⛽ ${car.fuel_type}</div>
                    <div class="spec-item">🎨 ${car.color || 'N/A'}</div>
                    <div class="spec-item">📏 ${car.engine_size || 'N/A'} л</div>
                </div>
                <div class="car-price">${formatPrice(car.price)} ₸</div>
                <div class="car-actions">
                    ${isInGarage(car.id)
        ? `<button class="btn btn-secondary" disabled>В гараже ✓</button>`
        : `<button class="btn btn-success" onclick="addToGarage('${car.id}')">В гараж</button>`
    }
                    <button class="btn btn-danger" onclick="deleteCar('${car.id}')">Удалить</button>
                </div>
            </div>
        </div>
    `).join('');
}

// Format price
function formatPrice(price) {
    return new Intl.NumberFormat('ru-RU').format(price);
}

// Garage functions
function isInGarage(carId) {
    return garage.some(car => car.id === carId);
}

function addToGarage(carId) {
    const car = allCars.find(c => c.id === carId);
    if (car && !isInGarage(carId)) {
        garage.push(car);
        localStorage.setItem('garage', JSON.stringify(garage));
        updateGarageCount();
        renderCars();
    }
}

function removeFromGarage(carId) {
    garage = garage.filter(car => car.id !== carId);
    localStorage.setItem('garage', JSON.stringify(garage));
    updateGarageCount();
    renderGarage();
    renderCars();
}

function updateGarageCount() {
    document.getElementById('garage-count').textContent = garage.length;
}

function renderGarage() {
    const emptyState = document.getElementById('garage-empty');
    const compareContainer = document.getElementById('garage-compare');
    const garageCars = document.getElementById('garage-cars');

    if (garage.length === 0) {
        emptyState.style.display = 'block';
        compareContainer.style.display = 'none';
        return;
    }

    emptyState.style.display = 'none';
    compareContainer.style.display = 'block';

    garageCars.innerHTML = garage.map(car => `
        <div class="compare-card">
            <button class="remove-from-garage" onclick="removeFromGarage('${car.id}')">×</button>
            <img src="${car.image_url || 'https://via.placeholder.com/400x300?text=' + car.make + '+' + car.model}" 
                 alt="${car.make} ${car.model}" class="car-image">
            <div class="car-title">${car.make} ${car.model}</div>
            <div class="compare-specs">
                <div class="compare-spec-row">
                    <span class="spec-label">Год:</span>
                    <span class="spec-value">${car.year}</span>
                </div>
                <div class="compare-spec-row">
                    <span class="spec-label">Цена:</span>
                    <span class="spec-value">${formatPrice(car.price)} ₸</span>
                </div>
                <div class="compare-spec-row">
                    <span class="spec-label">Пробег:</span>
                    <span class="spec-value">${car.mileage ? car.mileage.toLocaleString() + ' км' : 'N/A'}</span>
                </div>
                <div class="compare-spec-row">
                    <span class="spec-label">Мощность:</span>
                    <span class="spec-value">${car.horsepower || 'N/A'} л.с.</span>
                </div>
                <div class="compare-spec-row">
                    <span class="spec-label">Двигатель:</span>
                    <span class="spec-value">${car.engine_size || 'N/A'} л</span>
                </div>
                <div class="compare-spec-row">
                    <span class="spec-label">Тип кузова:</span>
                    <span class="spec-value">${car.body_type}</span>
                </div>
                <div class="compare-spec-row">
                    <span class="spec-label">Топливо:</span>
                    <span class="spec-value">${car.fuel_type}</span>
                </div>
                <div class="compare-spec-row">
                    <span class="spec-label">КПП:</span>
                    <span class="spec-value">${car.transmission}</span>
                </div>
                <div class="compare-spec-row">
                    <span class="spec-label">Цвет:</span>
                    <span class="spec-value">${car.color || 'N/A'}</span>
                </div>
            </div>
        </div>
    `).join('');
}

// Filters
document.getElementById('apply-filters').addEventListener('click', () => {
    const filters = {
        make: document.getElementById('filter-make').value,
        body_type: document.getElementById('filter-body').value,
        fuel_type: document.getElementById('filter-fuel').value,
        transmission: document.getElementById('filter-transmission').value,
        min_price: document.getElementById('filter-price-min').value,
        max_price: document.getElementById('filter-price-max').value,
    };

    fetchCars(filters);
});

document.getElementById('reset-filters').addEventListener('click', () => {
    document.getElementById('filter-make').value = '';
    document.getElementById('filter-body').value = '';
    document.getElementById('filter-fuel').value = '';
    document.getElementById('filter-transmission').value = '';
    document.getElementById('filter-price-min').value = '';
    document.getElementById('filter-price-max').value = '';
    fetchCars();
});

// Add car form
document.getElementById('add-car-form').addEventListener('submit', async (e) => {
    e.preventDefault();

    const formData = new FormData(e.target);
    const carData = {
        make: formData.get('make'),
        model: formData.get('model'),
        year: parseInt(formData.get('year')),
        price: parseFloat(formData.get('price')),
        mileage: parseInt(formData.get('mileage')) || 0,
        horsepower: parseInt(formData.get('horsepower')) || 0,
        engine_size: parseFloat(formData.get('engine_size')) || 0,
        color: formData.get('color') || '',
        body_type: formData.get('body_type'),
        fuel_type: formData.get('fuel_type'),
        transmission: formData.get('transmission'),
        description: formData.get('description') || '',
        image_url: formData.get('image_url') || '',
    };

    try {
        const response = await fetch(`${API_URL}/cars`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(carData),
        });

        if (response.ok) {
            alert('Автомобиль успешно добавлен!');
            e.target.reset();
            fetchCars();
            showPage('catalog-page');
        } else {
            alert('Ошибка при добавлении автомобиля');
        }
    } catch (error) {
        console.error('Error adding car:', error);
        alert('Ошибка при добавлении автомобиля');
    }
});

// Delete car
async function deleteCar(carId) {
    if (!confirm('Вы уверены, что хотите удалить этот автомобиль?')) {
        return;
    }

    try {
        const response = await fetch(`${API_URL}/cars/${carId}`, {
            method: 'DELETE',
        });

        if (response.ok) {
            alert('Автомобиль удален');
            removeFromGarage(carId);
            fetchCars();
        } else {
            alert('Ошибка при удалении автомобиля');
        }
    } catch (error) {
        console.error('Error deleting car:', error);
        alert('Ошибка при удалении автомобиля');
    }
}

// Initialize
updateGarageCount();
fetchCars();