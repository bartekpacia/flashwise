CREATE DATABASE IF NOT EXISTS flashwise;

USE flashwise;

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_admin BOOLEAN DEFAULT FALSE,
    password_hash BINARY(60) NOT NULL,
    token VARCHAR(40) NOT NULL
);

CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL
);

CREATE TABLE flashcard_sets (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    is_public BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    modified_at TIMESTAMP DEFAULT NULL,
    author_id BIGINT UNSIGNED NOT NULL,
    category_id BIGINT UNSIGNED NOT NULL,
    FOREIGN KEY (author_id) REFERENCES users (id),
    FOREIGN KEY (category_id) REFERENCES categories (id)
);

CREATE TABLE flashcards (
    id SERIAL PRIMARY KEY,
    front VARCHAR(512) NOT NULL,
    back VARCHAR(512) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    modified_at TIMESTAMP DEFAULT NULL,
    author_id BIGINT UNSIGNED NOT NULL,
    set_id BIGINT UNSIGNED NOT NULL,
    FOREIGN KEY (author_id) REFERENCES users (id),
    FOREIGN KEY (set_id) REFERENCES flashcard_sets (id)
);
