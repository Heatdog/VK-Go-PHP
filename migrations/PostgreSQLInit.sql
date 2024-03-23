CREATE TABLE IF NOT EXISTS users(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    login VARCHAR(50) UNIQUE NOT NULL,
    password VARCHAR(250) NOT NULL,
    CHECK (LENGTH(login) > 3 and LENGTH(password) > 3)
);

CREATE TABLE IF NOT EXISTS advert(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    body VARCHAR(1200) NOT NULL,
    price INTEGER NOT NULL,
    image_adr VARCHAR(500) NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id),
    CHECK (LENGTH(title) > 3 and LENGTH(body) > 3 and price > 0) 
);