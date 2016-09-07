ALTER TABLE tokens ADD COLUMN next timestamp;
CREATE EXTENSION hstore;
ALTER TABLE tokens ADD COLUMN followers bigint[];
