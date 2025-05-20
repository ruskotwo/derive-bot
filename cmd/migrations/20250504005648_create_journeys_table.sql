-- +goose Up
CREATE TABLE IF NOT EXISTS journeys
(
    id               INT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
    user_id          INT(11) UNSIGNED NOT NULL DEFAULT 0,
    quest_id         INT(11) UNSIGNED NOT NULL DEFAULT 0,
    progress         TINYINT UNSIGNED NOT NULL DEFAULT 0,
    complete_till_at TIMESTAMP        NULL     DEFAULT NULL,
    created_at       TIMESTAMP        NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at       TIMESTAMP        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    KEY journeys_ix_user_id_progress (user_id, progress)
);
-- +goose Down
DROP TABLE IF EXISTS journeys;
