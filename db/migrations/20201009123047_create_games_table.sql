-- migrate:up

CREATE TABLE IF NOT EXISTS games (
  id VARCHAR(255) PRIMARY KEY,
  turn BIGINT NOT NULL,
  result BIGINT NOT NULL,
  board_side1 JSONB NOT NULL DEFAULT '{}'::jsonb,
  board_side2 JSONB NOT NULL DEFAULT '{}'::jsonb
);

-- migrate:down

DROP TABLE IF EXISTS games;
