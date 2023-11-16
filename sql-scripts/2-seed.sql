USE flashwise;

INSERT INTO users (username, email, is_admin, password_hash, token)
VALUES (
    'admin',
    'admin@flashwise.com',
    TRUE,
    '$2a$14$y1kbWJKhY5b2u8JwOmlbuOWLTY6tK9fsXYZ.btRImHl1o0Nh4Dfx2',
    'c8dd1617cf1f35965cb8d2f827f4c2d834f2958b'
), (
    'bartek',
    'barpac02@gmail.com',
    FALSE,
    '$2a$14$.w8xD7D9IGe3.ju4gSaviO/2BtX/137Rg3XhsLFj9roHr.pw718pi',
    '6d0c1a5ecb334e176c5d13e8d24c282a8b45684d'
);

INSERT INTO categories (title, slug)
VALUES ('Geografia', 'geografia'),
('Angielski', 'angielski'),
('Niemiecki', 'niemiecki'),
('Historia', 'historia');

INSERT INTO flashcard_sets (title, author_id, category_id)
VALUES ('Ciekawostki o Polsce', 2, 1);

INSERT INTO flashcards (front, back, author_id, set_id)
VALUES ('Stolica Polski', 'Warszawa', 2, 1),
('Najpiękniejsze miasto Polski', 'Gliwice', 2, 1),
('Najpotężniejsza polska uczelnia', 'Politechnika Śląska', 2, 1);
