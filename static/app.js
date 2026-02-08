const API_URL = 'http://localhost:9000/api';

// State
let allCars = [];
let garage = JSON.parse(localStorage.getItem('garage')) || [];

// Auth Helper Functions
function getAuthHeaders() {
    const token = localStorage.getItem('token');
    const headers = {
        'Content-Type': 'application/json'
    };

    if (token) {
        headers['Authorization'] = `Bearer ${token}`;
    }

    return headers;
}

function checkAuth() {
    const token = localStorage.getItem('token');
    if (!token) {
        alert('You must be logged in to perform this action');
        window.location.href = 'login.html';
        return false;
    }
    return true;
}

function handleAuthError(response) {
    if (response.status === 401) {
        alert('Session expired. Please login again');
        localStorage.removeItem('token');
        localStorage.removeItem('user');
        window.location.href = 'login.html';
        return true;
    }
    return false;
}

// Update UI based on auth status
function updateAuthUI() {
    const token = localStorage.getItem('token');
    const user = JSON.parse(localStorage.getItem('user') || 'null');
    const authContainer = document.getElementById('auth-container');

    if (token && user) {
        authContainer.innerHTML = `
            <span class="user-name">Hello, ${user.username}!</span>
            <button class="btn btn-secondary" onclick="logout()">Logout</button>
        `;
    } else {
        authContainer.innerHTML = `
            <a href="login.html" class="btn btn-primary">Login</a>
            <a href="register.html" class="btn btn-success">Register</a>
        `;
    }
}

function logout() {
    localStorage.removeItem('token');
    localStorage.removeItem('user');
    window.location.reload();
}

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
        document.getElementById('cars-list').innerHTML = '<p>Error loading cars</p>';
    }
}

// Render cars
function renderCars() {
    const carsList = document.getElementById('cars-list');
    const isAuthenticated = !!localStorage.getItem('token');

    if (allCars.length === 0) {
        carsList.innerHTML = '<p>No cars found</p>';
        return;
    }

    carsList.innerHTML = allCars.map(car => `
        <div class="car-card">
            <img src="${car.image_url || 'https://via.placeholder.com/400x300?text=' + car.make + '+' + car.model}" 
                 alt="${car.make} ${car.model}" class="car-image">
            <div class="car-info">
                <div class="car-title">${car.make} ${car.model}</div>
                <div class="car-specs">
                    <div class="spec-item">Year: ${car.year}</div>
                    <div class="spec-item">HP: ${car.horsepower || 'N/A'}</div>
                    <div class="spec-item">Trans: ${car.transmission}</div>
                    <div class="spec-item">Fuel: ${car.fuel_type}</div>
                    <div class="spec-item">Color: ${car.color || 'N/A'}</div>
                    <div class="spec-item">Engine: ${car.engine_size || 'N/A'} L</div>
                </div>
                <div class="car-price">$${formatPrice(car.price)}</div>
                <div class="car-actions">
                    ${isInGarage(car.id)
        ? `<button class="btn btn-secondary" disabled>In Garage ✓</button>`
        : `<button class="btn btn-success" onclick="addToGarage('${car.id}')">Add to Garage</button>`
    }
                    ${isAuthenticated
        ? `
                        <button class="btn btn-primary" onclick="openEditModal('${car.id}')">Edit</button>
                        <button class="btn btn-danger" onclick="deleteCar('${car.id}')">Delete</button>
                        `
        : ''
    }
                </div>
            </div>
        </div>
    `).join('');
}

// Format price
function formatPrice(price) {
    return new Intl.NumberFormat('en-US').format(price);
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
                    <span class="spec-label">Year:</span>
                    <span class="spec-value">${car.year}</span>
                </div>
                <div class="compare-spec-row">
                    <span class="spec-label">Price:</span>
                    <span class="spec-value">$${formatPrice(car.price)}</span>
                </div>
                <div class="compare-spec-row">
                    <span class="spec-label">Mileage:</span>
                    <span class="spec-value">${car.mileage ? car.mileage.toLocaleString() + ' km' : 'N/A'}</span>
                </div>
                <div class="compare-spec-row">
                    <span class="spec-label">Horsepower:</span>
                    <span class="spec-value">${car.horsepower || 'N/A'} hp</span>
                </div>
                <div class="compare-spec-row">
                    <span class="spec-label">Engine:</span>
                    <span class="spec-value">${car.engine_size || 'N/A'} L</span>
                </div>
                <div class="compare-spec-row">
                    <span class="spec-label">Body Type:</span>
                    <span class="spec-value">${car.body_type}</span>
                </div>
                <div class="compare-spec-row">
                    <span class="spec-label">Fuel:</span>
                    <span class="spec-value">${car.fuel_type}</span>
                </div>
                <div class="compare-spec-row">
                    <span class="spec-label">Transmission:</span>
                    <span class="spec-value">${car.transmission}</span>
                </div>
                <div class="compare-spec-row">
                    <span class="spec-label">Color:</span>
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

    if (!checkAuth()) {
        return;
    }

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
            headers: getAuthHeaders(),
            body: JSON.stringify(carData),
        });

        if (handleAuthError(response)) {
            return;
        }

        if (response.ok) {
            alert('Car added successfully!');
            e.target.reset();
            fetchCars();
            showPage('catalog-page');
        } else {
            const error = await response.text();
            alert('Error adding car: ' + error);
        }
    } catch (error) {
        console.error('Error adding car:', error);
        alert('Error adding car');
    }
});

// Edit Modal Functions
function openEditModal(carId) {
    const car = allCars.find(c => c.id === carId);
    if (!car) return;

    document.getElementById('edit-car-id').value = car.id;
    document.getElementById('edit-make').value = car.make;
    document.getElementById('edit-model').value = car.model;
    document.getElementById('edit-year').value = car.year;
    document.getElementById('edit-price').value = car.price;
    document.getElementById('edit-mileage').value = car.mileage || '';
    document.getElementById('edit-horsepower').value = car.horsepower || '';
    document.getElementById('edit-engine-size').value = car.engine_size || '';
    document.getElementById('edit-color').value = car.color || '';
    document.getElementById('edit-body-type').value = car.body_type;
    document.getElementById('edit-fuel-type').value = car.fuel_type;
    document.getElementById('edit-transmission').value = car.transmission;
    document.getElementById('edit-image-url').value = car.image_url || '';
    document.getElementById('edit-description').value = car.description || '';

    document.getElementById('edit-modal').style.display = 'block';
}

function closeEditModal() {
    document.getElementById('edit-modal').style.display = 'none';
}

// Close modal when clicking X
document.querySelector('.modal-close').addEventListener('click', closeEditModal);

// Close modal when clicking outside
window.addEventListener('click', (e) => {
    const modal = document.getElementById('edit-modal');
    if (e.target === modal) {
        closeEditModal();
    }
});

// Edit car form
document.getElementById('edit-car-form').addEventListener('submit', async (e) => {
    e.preventDefault();

    if (!checkAuth()) {
        return;
    }

    const carId = document.getElementById('edit-car-id').value;
    const carData = {
        make: document.getElementById('edit-make').value,
        model: document.getElementById('edit-model').value,
        year: parseInt(document.getElementById('edit-year').value),
        price: parseFloat(document.getElementById('edit-price').value),
        mileage: parseInt(document.getElementById('edit-mileage').value) || 0,
        horsepower: parseInt(document.getElementById('edit-horsepower').value) || 0,
        engine_size: parseFloat(document.getElementById('edit-engine-size').value) || 0,
        color: document.getElementById('edit-color').value || '',
        body_type: document.getElementById('edit-body-type').value,
        fuel_type: document.getElementById('edit-fuel-type').value,
        transmission: document.getElementById('edit-transmission').value,
        description: document.getElementById('edit-description').value || '',
        image_url: document.getElementById('edit-image-url').value || '',
    };

    try {
        const response = await fetch(`${API_URL}/cars/${carId}`, {
            method: 'PUT',
            headers: getAuthHeaders(),
            body: JSON.stringify(carData),
        });

        if (handleAuthError(response)) {
            return;
        }

        if (response.ok) {
            alert('Car updated successfully!');
            closeEditModal();
            fetchCars();
        } else {
            const error = await response.text();
            alert('Error updating car: ' + error);
        }
    } catch (error) {
        console.error('Error updating car:', error);
        alert('Error updating car');
    }
});

// Delete car
async function deleteCar(carId) {
    if (!checkAuth()) {
        return;
    }

    if (!confirm('Are you sure you want to delete this car?')) {
        return;
    }

    try {
        const response = await fetch(`${API_URL}/cars/${carId}`, {
            method: 'DELETE',
            headers: getAuthHeaders(),
        });

        if (handleAuthError(response)) {
            return;
        }

        if (response.ok) {
            alert('Car deleted successfully');
            removeFromGarage(carId);
            fetchCars();
        } else {
            const error = await response.text();
            alert('Error deleting car: ' + error);
        }
    } catch (error) {
        console.error('Error deleting car:', error);
        alert('Error deleting car');
    }
}

// Initialize
updateAuthUI();
updateGarageCount();
fetchCars();