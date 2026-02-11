const API_URL = 'http://localhost:3000/api';

// State
let allCars = [];
let garage = JSON.parse(localStorage.getItem('garage')) || [];
let favorites = [];

// ============= AUTH HELPER FUNCTIONS =============

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
        console.error('Authentication error: 401 Unauthorized');
        alert('Session expired. Please login again');
        localStorage.removeItem('token');
        localStorage.removeItem('user');
        window.location.href = 'login.html';
        return true;
    }
    return false;
}

function updateAuthUI() {
    const token = localStorage.getItem('token');
    const user = JSON.parse(localStorage.getItem('user') || 'null');
    const authContainer = document.getElementById('auth-container');

    console.log('updateAuthUI called. Token:', token ? 'exists' : 'missing', 'User:', user);

    if (token && user) {
        authContainer.innerHTML = `
            <span class="user-name">Hello, ${user.username}!</span>
            <button class="btn btn-secondary" onclick="logout()">Logout</button>
        `;
        // Load favorites count ONLY if user is logged in
        // Добавляем небольшую задержку чтобы токен успел прогрузиться
        setTimeout(() => {
            loadFavoritesCount();
        }, 100);
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

// ============= NAVIGATION =============

document.getElementById('catalog-btn').addEventListener('click', () => showPage('catalog-page'));
document.getElementById('garage-btn').addEventListener('click', () => showPage('garage-page'));
document.getElementById('add-car-btn').addEventListener('click', () => showPage('add-car-page'));
document.getElementById('favorites-btn').addEventListener('click', () => {
    if (checkAuth()) {
        showPage('favorites-page');
        loadFavoritesPage();
    }
});

function showPage(pageId) {
    document.querySelectorAll('.page').forEach(page => page.classList.remove('active'));
    document.querySelectorAll('.nav-btn').forEach(btn => btn.classList.remove('active'));

    document.getElementById(pageId).classList.add('active');

    if (pageId === 'catalog-page') {
        document.getElementById('catalog-btn').classList.add('active');
    } else if (pageId === 'garage-page') {
        document.getElementById('garage-btn').classList.add('active');
        displayGarage();
    } else if (pageId === 'add-car-page') {
        document.getElementById('add-car-btn').classList.add('active');
    } else if (pageId === 'favorites-page') {
        document.getElementById('favorites-btn').classList.add('active');
    }
}

// ============= FAVORITES FUNCTIONS (SERIK) =============

async function loadFavoritesCount() {
    const token = localStorage.getItem('token');

    // Не делаем запрос если нет токена
    if (!token) {
        console.log('No token found, skipping favorites count load');
        return;
    }

    try {
        console.log('Loading favorites count...');
        const response = await fetch(`${API_URL}/favorites/count`, {
            headers: getAuthHeaders()
        });

        console.log('Favorites count response status:', response.status);

        if (response.status === 401) {
            // НЕ показываем alert при первой загрузке!
            // Просто логируем ошибку
            console.error('Token is invalid or expired');
            localStorage.removeItem('token');
            localStorage.removeItem('user');
            // Перезагружаем UI без редиректа
            updateAuthUI();
            return;
        }

        if (response.ok) {
            const data = await response.json();
            document.getElementById('favorites-count').textContent = data.count || 0;
            console.log('Favorites count loaded:', data.count);
        } else {
            console.error('Failed to load favorites count:', response.status);
        }
    } catch (error) {
        console.error('Error loading favorites count:', error);
        // Не показываем alert на ошибках сети при загрузке
    }
}

async function toggleFavorite(carId, button) {
    if (!checkAuth()) return;

    const isFavorite = button.classList.contains('active');

    try {
        let response;
        if (isFavorite) {
            // Remove from favorites
            response = await fetch(`${API_URL}/favorites/${carId}`, {
                method: 'DELETE',
                headers: getAuthHeaders()
            });
        } else {
            // Add to favorites
            response = await fetch(`${API_URL}/favorites`, {
                method: 'POST',
                headers: getAuthHeaders(),
                body: JSON.stringify({ car_id: carId })
            });
        }

        if (handleAuthError(response)) return;

        if (response.ok) {
            button.classList.toggle('active');
            button.innerHTML = button.classList.contains('active') ? '❤️' : '🤍';
            loadFavoritesCount();

            // Reload favorites page if we're on it
            if (document.getElementById('favorites-page').classList.contains('active')) {
                loadFavoritesPage();
            }
        } else {
            alert('Failed to update favorites');
        }
    } catch (error) {
        console.error('Error toggling favorite:', error);
        alert('Error updating favorites');
    }
}

async function loadFavoritesPage() {
    const container = document.getElementById('favorites-container');
    container.innerHTML = '<div class="loading">Loading favorites...</div>';

    try {
        const response = await fetch(`${API_URL}/favorites`, {
            headers: getAuthHeaders()
        });

        if (handleAuthError(response)) return;

        if (response.ok) {
            favorites = await response.json();
            displayFavorites();
        } else {
            container.innerHTML = '<div class="error">Failed to load favorites</div>';
        }
    } catch (error) {
        console.error('Error loading favorites:', error);
        container.innerHTML = '<div class="error">Error loading favorites</div>';
    }
}

function displayFavorites() {
    const container = document.getElementById('favorites-container');

    if (!favorites || favorites.length === 0) {
        container.innerHTML = '<div class="empty-state">No favorites yet. Add some cars to your favorites!</div>';
        return;
    }

    container.innerHTML = favorites.map(fav => createCarCard(fav.car, true)).join('');
}

// ============= REVIEWS FUNCTIONS (DIMASH) =============

async function showReviewsModal(carId) {
    const modal = document.getElementById('reviews-modal');
    const modalContent = document.getElementById('reviews-content');

    modal.style.display = 'block';
    modalContent.innerHTML = '<div class="loading">Loading reviews...</div>';

    try {
        const response = await fetch(`${API_URL}/cars/${carId}/reviews`);

        if (response.ok) {
            const data = await response.json();
            displayReviews(data, carId);
        } else {
            modalContent.innerHTML = '<div class="error">Failed to load reviews</div>';
        }
    } catch (error) {
        console.error('Error loading reviews:', error);
        modalContent.innerHTML = '<div class="error">Error loading reviews</div>';
    }
}

function displayReviews(data, carId) {
    const modalContent = document.getElementById('reviews-content');
    const token = localStorage.getItem('token');
    const currentUser = JSON.parse(localStorage.getItem('user') || 'null');

    let html = `
        <div class="reviews-header">
            <h3>Reviews</h3>
            <div class="average-rating">
                ${data.average_rating ? `
                    <span class="rating-stars">${'⭐'.repeat(Math.round(data.average_rating))}</span>
                    <span class="rating-value">${data.average_rating.toFixed(1)} / 5</span>
                    <span class="review-count">(${data.reviews.length} reviews)</span>
                ` : '<span>No reviews yet</span>'}
            </div>
        </div>
    `;

    if (token) {
        html += `
            <div class="add-review-form">
                <h4>Write a Review</h4>
                <form id="review-form" onsubmit="submitReview(event, '${carId}')">
                    <div class="rating-input">
                        <label>Rating:</label>
                        <select name="rating" required>
                            <option value="">Select rating</option>
                            <option value="5">⭐⭐⭐⭐⭐ (5)</option>
                            <option value="4">⭐⭐⭐⭐ (4)</option>
                            <option value="3">⭐⭐⭐ (3)</option>
                            <option value="2">⭐⭐ (2)</option>
                            <option value="1">⭐ (1)</option>
                        </select>
                    </div>
                    <div class="form-group">
                        <label>Your Review:</label>
                        <textarea name="comment" rows="4" required placeholder="Share your experience with this car..."></textarea>
                    </div>
                    <button type="submit" class="btn btn-primary">Submit Review</button>
                </form>
            </div>
        `;
    }

    html += '<div class="reviews-list">';

    if (data.reviews && data.reviews.length > 0) {
        data.reviews.forEach(review => {
            const isOwnReview = currentUser && review.user_id === currentUser.id;
            html += `
                <div class="review-item" data-review-id="${review.id}">
                    <div class="review-header">
                        <span class="review-author">${review.username || 'Anonymous'}</span>
                        <span class="review-rating">${'⭐'.repeat(review.rating)}</span>
                    </div>
                    <p class="review-comment">${review.comment}</p>
                    <div class="review-footer">
                        <span class="review-date">${new Date(review.created_at).toLocaleDateString()}</span>
                        ${isOwnReview ? `
                            <div class="review-actions">
                                <button class="btn btn-sm btn-secondary" onclick="editReview('${review.id}', '${review.comment}', ${review.rating})">Edit</button>
                                <button class="btn btn-sm btn-danger" onclick="deleteReview('${review.id}')">Delete</button>
                            </div>
                        ` : ''}
                    </div>
                </div>
            `;
        });
    } else {
        html += '<div class="empty-state">No reviews yet. Be the first to review this car!</div>';
    }

    html += '</div>';
    modalContent.innerHTML = html;
}

async function submitReview(event, carId) {
    event.preventDefault();

    if (!checkAuth()) return;

    const form = event.target;
    const formData = new FormData(form);
    const rating = parseInt(formData.get('rating'));
    const comment = formData.get('comment');

    try {
        const response = await fetch(`${API_URL}/cars/${carId}/reviews`, {
            method: 'POST',
            headers: getAuthHeaders(),
            body: JSON.stringify({ rating, comment })
        });

        if (handleAuthError(response)) return;

        if (response.ok) {
            showReviewsModal(carId);
        } else {
            const error = await response.text();
            alert(`Error: ${error}`);
        }
    } catch (error) {
        console.error('Error submitting review:', error);
        alert('Failed to submit review');
    }
}

async function deleteReview(reviewId) {
    if (!confirm('Are you sure you want to delete this review?')) return;

    try {
        const response = await fetch(`${API_URL}/reviews/${reviewId}`, {
            method: 'DELETE',
            headers: getAuthHeaders()
        });

        if (handleAuthError(response)) return;

        if (response.ok) {
            // Reload the modal - get carId from the current modal
            const modal = document.getElementById('reviews-modal');
            const carId = modal.dataset.carId;
            if (carId) {
                showReviewsModal(carId);
            } else {
                closeModal('reviews-modal');
            }
        } else {
            alert('Failed to delete review');
        }
    } catch (error) {
        console.error('Error deleting review:', error);
        alert('Error deleting review');
    }
}

function editReview(reviewId, currentComment, currentRating) {
    const reviewItem = document.querySelector(`[data-review-id="${reviewId}"]`);
    const commentElement = reviewItem.querySelector('.review-comment');
    const actionsElement = reviewItem.querySelector('.review-actions');

    const editForm = `
        <form class="edit-review-form" onsubmit="updateReview(event, '${reviewId}')">
            <div class="rating-input">
                <label>Rating:</label>
                <select name="rating" required>
                    <option value="5" ${currentRating === 5 ? 'selected' : ''}>⭐⭐⭐⭐⭐ (5)</option>
                    <option value="4" ${currentRating === 4 ? 'selected' : ''}>⭐⭐⭐⭐ (4)</option>
                    <option value="3" ${currentRating === 3 ? 'selected' : ''}>⭐⭐⭐ (3)</option>
                    <option value="2" ${currentRating === 2 ? 'selected' : ''}>⭐⭐ (2)</option>
                    <option value="1" ${currentRating === 1 ? 'selected' : ''}>⭐ (1)</option>
                </select>
            </div>
            <textarea name="comment" rows="3" required>${currentComment}</textarea>
            <div class="edit-actions">
                <button type="submit" class="btn btn-sm btn-primary">Save</button>
                <button type="button" class="btn btn-sm btn-secondary" onclick="showReviewsModal('${document.getElementById('reviews-modal').dataset.carId}')">Cancel</button>
            </div>
        </form>
    `;

    commentElement.innerHTML = editForm;
    actionsElement.style.display = 'none';
}

async function updateReview(event, reviewId) {
    event.preventDefault();

    const form = event.target;
    const formData = new FormData(form);
    const rating = parseInt(formData.get('rating'));
    const comment = formData.get('comment');

    try {
        const response = await fetch(`${API_URL}/reviews/${reviewId}`, {
            method: 'PUT',
            headers: getAuthHeaders(),
            body: JSON.stringify({ rating, comment })
        });

        if (handleAuthError(response)) return;

        if (response.ok) {
            // Reload the modal - get carId from the current modal
            const modal = document.getElementById('reviews-modal');
            const carId = modal.dataset.carId;
            if (carId) {
                showReviewsModal(carId);
            } else {
                closeModal('reviews-modal');
            }
        } else {
            alert('Failed to update review');
        }
    } catch (error) {
        console.error('Error updating review:', error);
        alert('Error updating review');
    }
}

function closeModal(modalId) {
    document.getElementById(modalId).style.display = 'none';
}

// ============= CAR DISPLAY FUNCTIONS =============

async function fetchCars() {
    try {
        const response = await fetch(`${API_URL}/cars`);
        if (response.ok) {
            allCars = await response.json();
            displayCars(allCars);
        }
    } catch (error) {
        console.error('Error fetching cars:', error);
    }
}

function createCarCard(car, isFavorite = false) {
    const inGarage = garage.some(g => g.id === car.id);
    const token = localStorage.getItem('token');

    return `
        <div class="car-card">
            ${token ? `<button class="favorite-btn ${isFavorite ? 'active' : ''}" onclick="toggleFavorite('${car.id}', this)">${isFavorite ? '❤️' : '🤍'}</button>` : ''}
            <img src="${car.image_url || 'https://via.placeholder.com/300x200?text=No+Image'}" alt="${car.make} ${car.model}">
            <div class="car-info">
                <h3>${car.make} ${car.model}</h3>
                <p class="car-year">${car.year}</p>
                <p class="car-price">$${car.price.toLocaleString()}</p>
                <div class="car-specs">
                    <span>${car.mileage.toLocaleString()} km</span>
                    <span>${car.fuel_type}</span>
                    <span>${car.transmission}</span>
                </div>
                <p class="car-description">${car.description || 'No description available'}</p>
                <div class="car-actions">
                    ${!inGarage ?
        `<button class="btn btn-primary" onclick="addToGarage('${car.id}')">Add to Garage</button>` :
        `<button class="btn btn-secondary" disabled>In Garage</button>`
    }
                    <button class="btn btn-info" onclick="showReviewsModal('${car.id}')">Reviews</button>
                </div>
            </div>
        </div>
    `;
}

function displayCars(cars) {
    const container = document.getElementById('cars-container');

    if (!cars || cars.length === 0) {
        container.innerHTML = '<div class="empty-state">No cars found</div>';
        return;
    }

    container.innerHTML = cars.map(car => createCarCard(car)).join('');
}

function addToGarage(carId) {
    const car = allCars.find(c => c.id === carId);
    if (car && !garage.some(g => g.id === carId)) {
        garage.push(car);
        localStorage.setItem('garage', JSON.stringify(garage));
        updateGarageCount();
        displayCars(allCars);
    }
}

function removeFromGarage(carId) {
    garage = garage.filter(car => car.id !== carId);
    localStorage.setItem('garage', JSON.stringify(garage));
    updateGarageCount();
    displayGarage();
}

function updateGarageCount() {
    document.getElementById('garage-count').textContent = garage.length;
}

function displayGarage() {
    const container = document.getElementById('garage-container');

    if (garage.length === 0) {
        container.innerHTML = '<div class="empty-state">Your garage is empty. Add some cars!</div>';
        return;
    }

    container.innerHTML = garage.map(car => `
        <div class="car-card">
            <img src="${car.image_url || 'https://via.placeholder.com/300x200?text=No+Image'}" alt="${car.make} ${car.model}">
            <div class="car-info">
                <h3>${car.make} ${car.model}</h3>
                <p class="car-year">${car.year}</p>
                <p class="car-price">$${car.price.toLocaleString()}</p>
                <div class="car-specs">
                    <span>${car.mileage.toLocaleString()} km</span>
                    <span>${car.fuel_type}</span>
                    <span>${car.transmission}</span>
                </div>
                <button class="btn btn-danger" onclick="removeFromGarage('${car.id}')">Remove</button>
            </div>
        </div>
    `).join('');
}

// ============= FILTER FUNCTIONS =============

document.getElementById('apply-filters').addEventListener('click', applyFilters);
document.getElementById('reset-filters').addEventListener('click', resetFilters);

function applyFilters() {
    const make = document.getElementById('filter-make').value;
    const bodyType = document.getElementById('filter-body').value;
    const fuelType = document.getElementById('filter-fuel').value;
    const transmission = document.getElementById('filter-transmission').value;
    const minPrice = parseFloat(document.getElementById('filter-price-min').value) || 0;
    const maxPrice = parseFloat(document.getElementById('filter-price-max').value) || Infinity;

    const filtered = allCars.filter(car => {
        return (!make || car.make === make) &&
            (!bodyType || car.body_type === bodyType) &&
            (!fuelType || car.fuel_type === fuelType) &&
            (!transmission || car.transmission === transmission) &&
            (car.price >= minPrice && car.price <= maxPrice);
    });

    displayCars(filtered);
}

function resetFilters() {
    document.getElementById('filter-make').value = '';
    document.getElementById('filter-body').value = '';
    document.getElementById('filter-fuel').value = '';
    document.getElementById('filter-transmission').value = '';
    document.getElementById('filter-price-min').value = '';
    document.getElementById('filter-price-max').value = '';
    displayCars(allCars);
}

// ============= ADD CAR FORM =============

document.getElementById('add-car-form').addEventListener('submit', async (e) => {
    e.preventDefault();

    if (!checkAuth()) return;

    const formData = {
        make: document.getElementById('car-make').value,
        model: document.getElementById('car-model').value,
        year: parseInt(document.getElementById('car-year').value),
        price: parseFloat(document.getElementById('car-price').value),
        mileage: parseInt(document.getElementById('car-mileage').value),
        body_type: document.getElementById('car-body').value,
        fuel_type: document.getElementById('car-fuel').value,
        transmission: document.getElementById('car-transmission').value,
        color: document.getElementById('car-color').value,
        horsepower: parseInt(document.getElementById('car-horsepower').value),
        engine_size: parseFloat(document.getElementById('car-engine').value),
        description: document.getElementById('car-description').value,
        image_url: document.getElementById('car-image').value
    };

    try {
        const response = await fetch(`${API_URL}/cars`, {
            method: 'POST',
            headers: getAuthHeaders(),
            body: JSON.stringify(formData)
        });

        if (handleAuthError(response)) return;

        if (response.ok) {
            alert('Car added successfully!');
            document.getElementById('add-car-form').reset();
            fetchCars();
            showPage('catalog-page');
        } else {
            const error = await response.text();
            alert(`Error: ${error}`);
        }
    } catch (error) {
        console.error('Error adding car:', error);
        alert('Failed to add car');
    }
});

// ============= MODAL HANDLERS =============

// Close modal when clicking outside
window.onclick = function(event) {
    const reviewsModal = document.getElementById('reviews-modal');
    if (event.target === reviewsModal) {
        closeModal('reviews-modal');
    }
}

// ============= INITIALIZATION =============

console.log('Initializing app...');
updateAuthUI();
updateGarageCount();
fetchCars();

// Store carId in modal when opening
const originalShowReviewsModal = showReviewsModal;
showReviewsModal = function(carId) {
    document.getElementById('reviews-modal').dataset.carId = carId;
    originalShowReviewsModal(carId);
}