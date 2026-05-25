-- +goose Up
CREATE TABLE IF NOT EXISTS search_queries (
    query         VARCHAR(255) NOT NULL,
    language      VARCHAR(10)  NOT NULL DEFAULT 'en',
    count         INT          NOT NULL DEFAULT 1,
    last_searched TIMESTAMP    DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (query, language)
) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- +goose Down
DROP TABLE IF EXISTS search_queries;
