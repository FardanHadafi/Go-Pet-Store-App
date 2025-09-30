CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE pets (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    species species_enum NOT NULL,
    price DECIMAL(10,2) NOT NULL CHECK (price >= 0),
    created_by INTEGER REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO pets (name, species, price) VALUES
('Fluffy', 'cat', 299.99),
('Buddy', 'dog', 450.00),
('Tweety', 'bird', 89.99),
('Goldie', 'fish', 15.50);

ALTER TABLE pets ADD COLUMN IF NOT EXISTS species VARCHAR(100);
ALTER TABLE pets ADD COLUMN IF NOT EXISTS price DECIMAL(10, 2);

SELECT * from pets;

SELECT * from users;

DELETE  from pets;

delete from users;