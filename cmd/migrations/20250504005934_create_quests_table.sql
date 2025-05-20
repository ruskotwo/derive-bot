-- +goose Up
CREATE TABLE IF NOT EXISTS quests
(
    id          INT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
    description TEXT             NOT NULL,
    lang        VARCHAR(2)       NOT NULL DEFAULT 'ru',
    category_id INT(11) UNSIGNED NOT NULL DEFAULT 0,
    created_at  TIMESTAMP        NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    KEY quests_ix_lang_category_id (lang, category_id)
);
-- +goose Down
DROP TABLE IF EXISTS quests;
