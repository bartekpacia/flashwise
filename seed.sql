CREATE DATABASE flashwise;

USE flashwise;

CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  username VARCHAR(255) UNIQUE NOT NULL,
  email VARCHAR(255) UNIQUE NOT NULL,
  created_at TIMESTAMP DEFAULT NOW()
);

INSERT INTO users (username, email)
VALUES
    ('admin', 'admin@flashwise.com'),
    ('bartek', 'barpac02@gmail.com');


CREATE TABLE flashcard_sets (
    id SERIAL PRIMARY KEY,
    author BIGINT UNSIGNED NOT NULL,
    status ENUM('public', 'private') DEFAULT 'private',
    description VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    modified_at TIMESTAMP DEFAULT NULL,
    FOREIGN KEY (author) REFERENCES users(id)
);

INSERT INTO flashcard_sets (author, description)
VALUES
    (2, 'Ciekawostki o Polsce');


CREATE TABLE flashcards (
    id SERIAL PRIMARY KEY,
    front VARCHAR(512) NOT NULL,
    back VARCHAR(512) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    modified_at TIMESTAMP DEFAULT NULL,
    author BIGINT UNSIGNED NOT NULL,
    flashcard_set BIGINT UNSIGNED NOT NULL,
    FOREIGN KEY (author) REFERENCES users(id),
    FOREIGN KEY (flashcard_set) REFERENCES flashcard_sets(id)
);

INSERT INTO flashcards (front, back, author, flashcard_set)
VALUES
    ('Stolica Polski', 'Warszawa', 2, 1),
    ('Najpiękniejsze miasto Polski', 'Gliwice', 2, 1),
    ('Najpotężniejsza polska uczelnia', 'Politechnika Śląska', 2, 1);


-- SELECT flashcards.front, flashcards.back, users.username
--     -> FROM flashcards
--     -> JOIN users ON flashcards.author = users.id;
