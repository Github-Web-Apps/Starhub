CREATE EXTENSION hstore;

CREATE TABLE tokens(
  user_id bigint primary key,
  token text not null,
  next timestamp,
  followers text[] default array[]::varchar[],
  created_at timestamp default current_timestamp,
  updated_at timestamp default current_timestamp
);

CREATE INDEX idx_tokens_user_id ON tokens(user_id);
CREATE INDEX idx_tokens_next ON tokens(next);
