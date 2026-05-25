-- +goose Up
CREATE TABLE IF NOT EXISTS pages (
    title        TEXT         NOT NULL,
    url          VARCHAR(768) NOT NULL UNIQUE,
    language     VARCHAR(10)  NOT NULL DEFAULT 'en',
    last_updated TIMESTAMP    NULL,
    content      LONGTEXT     NOT NULL,
    PRIMARY KEY (url),
    CONSTRAINT chk_language CHECK (language IN ('en', 'da'))
);

-- +goose Down
DROP TABLE IF EXISTS pages;
