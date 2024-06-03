-- +goose Up
CREATE TABLE people (
                        id SERIAL PRIMARY KEY,
                        name VARCHAR(50),
                        surname VARCHAR(50),
                        patronymic VARCHAR(50),
                        age INT,
                        gender VARCHAR(10),
                        nationality VARCHAR(50)
);

-- +goose Down
DROP TABLE people;