-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE example
(
    user_id INT         NOT NULL,
    counter INT         NOT NULL,
    ctime   TIMESTAMP   NOT NULL DEFAULT NOW(),
    mtime   TIMESTAMP   NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id)
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE example;
