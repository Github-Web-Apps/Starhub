CREATE TABLE tokens(
  user_id bigint primary key,
  token text not null
);

CREATE INDEX idx_token_user_id ON tokens(user_id);
