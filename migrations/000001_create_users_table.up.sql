CREATE TYPE user_role AS ENUM ('user', 'admin');

CREATE TABLE users (
                       id SERIAL PRIMARY KEY,
                       username VARCHAR(255) UNIQUE NOT NULL,
                       password VARCHAR(255) NOT NULL,
                       hashed_password VARCHAR(255),
                       role user_role DEFAULT 'user'
);