CREATE TABLE IF NOT EXISTS Users(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    login VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    CHECK (LENGTH(login) > 3 and LENGTH(password) > 3)
);