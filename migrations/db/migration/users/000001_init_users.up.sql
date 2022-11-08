CREATE TABLE "users"
(
    "id"       bigint PRIMARY KEY,
    "username" varchar NOT NULL,
    "balance"  decimal NOT NULL
);
COMMENT
ON COLUMN "users"."balance" IS 'must be non-negative';
INSERT INTO users (id, username, balance)
VALUES (1, 'internal Wallet', 0);
