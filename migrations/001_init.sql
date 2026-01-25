CREATE TABLE users (
                       id SERIAL PRIMARY KEY,
                       username VARCHAR(50) UNIQUE NOT NULL,
                       email VARCHAR(100) UNIQUE NOT NULL,
                       password_hash VARCHAR(255) NOT NULL,
                       role VARCHAR(20) DEFAULT 'user'
);

CREATE TABLE cars (
                      id SERIAL PRIMARY KEY,
                      make VARCHAR(50) NOT NULL,
                      model VARCHAR(50) NOT NULL,
                      year INT NOT NULL,
                      price DECIMAL(12,2) NOT NULL,
                      mileage INT,
                      body_type VARCHAR(30),
                      description TEXT,
                      image_url VARCHAR(255)
);

CREATE TABLE orders (
                        id SERIAL PRIMARY KEY,
                        user_id INT REFERENCES users(id),
                        car_id INT REFERENCES cars(id),
                        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                        status VARCHAR(20) DEFAULT 'pending'
);