-- +goose Up
CREATE TABLE tasks (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    status VARCHAR(20) NOT NULL,
    due_date TIMESTAMP WITH TIME ZONE NOT NULL
);

-- +goose Down
DROP TABLE tasks;
