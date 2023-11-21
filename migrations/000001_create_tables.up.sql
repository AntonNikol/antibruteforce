CREATE TABLE IF NOT EXISTS whitelist
(
    id         SERIAL PRIMARY KEY,
    prefix     VARCHAR(150) not null,
    mask       VARCHAR(150) not null,
    created_at TIMESTAMP DEFAULT now()
);

CREATE TABLE IF NOT EXISTS blacklist
(
    id         SERIAL PRIMARY KEY,
    prefix     VARCHAR(150) not null,
    mask       VARCHAR(150) not null,
    created_at TIMESTAMP DEFAULT now()
);
