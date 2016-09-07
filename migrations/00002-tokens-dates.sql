ALTER TABLE tokens ADD COLUMN created_at timestamp default current_timestamp;
ALTER TABLE tokens ADD COLUMN updated_at timestamp default current_timestamp;
