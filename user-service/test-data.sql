-- Test Data for User Service

-- Insert Movies
INSERT INTO movies (id, title, genre, duration, language) VALUES
(1, 'Inception', 'Sci-Fi', 148, 'English'),
(2, 'The Dark Knight', 'Action', 152, 'English'),
(3, 'Interstellar', 'Sci-Fi', 169, 'English')
ON CONFLICT (id) DO NOTHING;

-- Insert Theaters
INSERT INTO theaters (id, name, city) VALUES
(1, 'PVR Cinemas', 'Mumbai'),
(2, 'INOX Megaplex', 'Delhi'),
(3, 'Cinepolis', 'Bangalore')
ON CONFLICT (id) DO NOTHING;

-- Insert Shows
INSERT INTO shows (id, movie_id, theater_id, start_time, price) VALUES
(1, 1, 1, '2024-01-20 18:00:00', 250.00),
(2, 1, 2, '2024-01-20 21:00:00', 300.00),
(3, 2, 1, '2024-01-21 15:00:00', 280.00),
(4, 2, 3, '2024-01-21 19:00:00', 320.00),
(5, 3, 2, '2024-01-22 17:00:00', 350.00)
ON CONFLICT (id) DO NOTHING;

-- Reset sequences
SELECT setval('movies_id_seq', (SELECT MAX(id) FROM movies));
SELECT setval('theaters_id_seq', (SELECT MAX(id) FROM theaters));
SELECT setval('shows_id_seq', (SELECT MAX(id) FROM shows));
