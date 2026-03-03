CREATE TABLE users (
                       id BIGSERIAL PRIMARY KEY,
                       full_name VARCHAR(150) NOT NULL,
                       email VARCHAR(150) UNIQUE NOT NULL,
                       password_hash TEXT NOT NULL,
                       role VARCHAR(20) NOT NULL CHECK (role IN ('admin','teacher','student')),
                       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);