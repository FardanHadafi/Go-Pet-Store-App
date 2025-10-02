-- USERS table
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- PETS table
CREATE TABLE pets (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    species VARCHAR(50) NOT NULL,
    price DECIMAL(10, 2) NOT NULL CHECK (price >= 0),
    created_by INT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

drop TABLE pets;

drop table users;

INSERT INTO
    pets (name, species, price)
VALUES ('Fluffy', 'cat', 299.99),
    ('Buddy', 'dog', 450.00),
    ('Tweety', 'bird', 89.99),
    ('Goldie', 'fish', 15.50);

ALTER TABLE pets ADD COLUMN IF NOT EXISTS species VARCHAR(100);

ALTER TABLE pets ADD COLUMN IF NOT EXISTS price DECIMAL(10, 2);

SELECT * from pets;

SELECT * from users;

DELETE from pets;
-- Add the missing updated_at column
ALTER TABLE users
ADD COLUMN updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP;

-- Change email length
ALTER TABLE users ALTER COLUMN email TYPE VARCHAR(255);

-- Change created_at type
ALTER TABLE users ALTER COLUMN created_at TYPE TIMESTAMP;

-- Add indexes
CREATE INDEX IF NOT EXISTS idx_users_email ON users (email);

CREATE INDEX IF NOT EXISTS idx_users_username ON users (username);

ALTER TABLE users ADD CONSTRAINT unique_email UNIQUE (email);