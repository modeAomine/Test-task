CREATE TABLE tokens (
                        id SERIAL PRIMARY KEY,
                        user_id INT NOT NULL,
                        token VARCHAR(255) NOT NULL,
                        created_at TIMESTAMP NOT NULL DEFAULT NOW(),
                        expires_at TIMESTAMP NOT NULL,
                        FOREIGN KEY (user_id) REFERENCES users(id)
);