-- +goose Up
CREATE TABLE IF NOT EXISTS search_queries (
    query         VARCHAR(255) NOT NULL,
    language      VARCHAR(10)  NOT NULL DEFAULT 'en',
    count         INT          NOT NULL DEFAULT 1,
    last_searched TIMESTAMP    DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (query, language)
);

-- +goose Down
DROP TABLE IF EXISTS search_queries;
