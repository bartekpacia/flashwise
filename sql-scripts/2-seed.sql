USE flashwise;

INSERT INTO users (username, email, is_admin, password_hash, token)
VALUES
    ('admin', 'admin@flashwise.com', TRUE, 'pwhash', 't0k3n'),
    ('bartek', 'barpac02@gmail.com', FALSE, 'pwhash', 't0k3n');

INSERT INTO flashcard_sets (author_id, description)
VALUES
    (2, 'Ciekawostki o Polsce');

INSERT INTO flashcards (front, back, author_id, set_id)
VALUES
    ('Stolica Polski', 'Warszawa', 2, 1),
    ('Najpiękniejsze miasto Polski', 'Gliwice', 2, 1),
    ('Najpotężniejsza polska uczelnia', 'Politechnika Śląska', 2, 1);


-- SELECT flashcards.front, flashcards.back, users.username
--     -> FROM flashcards
--     -> JOIN users ON flashcards.author_id = users.id;
