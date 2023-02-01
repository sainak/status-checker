CREATE TABLE websites (
    id SERIAL PRIMARY KEY,
    url TEXT UNIQUE NOT NULL,
    added_at TIMESTAMP NOT NULL

);

CREATE TABLE website_statuses (
    id SERIAL PRIMARY KEY,
    up BOOLEAN NOT NULL,
    time TIMESTAMP NOT NULL,
    website_id INTEGER NOT NULL,
    FOREIGN KEY (website_id) REFERENCES websites(id)
);


