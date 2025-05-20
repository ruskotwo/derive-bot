-- +goose Up
ALTER TABLE quests
ADD COLUMN file TEXT NULL DEFAULT NULL AFTER description,
LOCK=NONE;
-- +goose Down
ALTER TABLE quests
DROP COLUMN file,
LOCK=NONE;
