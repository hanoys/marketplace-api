CREATE TABLE IF NOT EXISTS advertisement(
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users,
    title VARCHAR(255),
    body VARCHAR(2048),
    image_url VARCHAR(2048),
    price money
)