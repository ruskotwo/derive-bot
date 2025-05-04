-- +goose Up
CREATE TABLE IF NOT EXISTS users
(
    id          INT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
    telegram_id INT(11) UNSIGNED NOT NULL,
    lang        VARCHAR(5)       NOT NULL DEFAULT 'ru',
    created_at  TIMESTAMP        NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    UNIQUE KEY users_ux_telegram_id (telegram_id)
);
-- +goose Down
DROP TABLE IF EXISTS users;
