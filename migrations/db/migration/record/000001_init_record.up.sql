CREATE TABLE "service_record"
(
    "id"         bigint PRIMARY KEY,
    "service_id" bigint      NOT NULL,
    "user_id"    bigint      NOT NULL,
    "value"      decimal     NOT NULL,
    "status"     varchar     NOT NULL,
    "time"       timestamptz NOT NULL DEFAULT (now())
);
CREATE TABLE "credit_record"
(
    "id"      bigserial PRIMARY KEY,
    "user_id" bigint      NOT NULL,
    "value"   decimal     NOT NULL,
    "status"  varchar     NOT NULL,
    "time"    timestamptz NOT NULL DEFAULT (now())
);
CREATE TABLE "transfer_record"
(
    "id"           bigserial PRIMARY KEY,
    "user_from_id" bigint      NOT NULL,
    "user_to_id"   bigint      NOT NULL,
    "value"        decimal     NOT NULL,
    "comment"      varchar     NOT NULL,
    "status"       varchar     NOT NULL,
    "time"         timestamptz NOT NULL DEFAULT (now())
);
