-- +goose Up
ALTER TABLE pages ADD FULLTEXT INDEX ft_title (title);
ALTER TABLE pages ADD FULLTEXT INDEX ft_content (content);

-- +goose Down
ALTER TABLE pages DROP INDEX ft_title;
ALTER TABLE pages DROP INDEX ft_content;
