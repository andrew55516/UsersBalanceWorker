CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "username" varchar NOT NULL,
  "balance" decimal NOT NULL
);

COMMENT ON COLUMN "users"."balance" IS 'must be non-negative';

INSERT INTO users (username, balance) VALUES ('internal Wallet', 0);