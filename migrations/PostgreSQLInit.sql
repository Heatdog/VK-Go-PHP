CREATE TABLE IF NOT EXISTS users(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    login VARCHAR(50) UNIQUE NOT NULL,
    password VARCHAR(250) NOT NULL,
    CHECK (LENGTH(login) > 3 and LENGTH(password) > 3)
);