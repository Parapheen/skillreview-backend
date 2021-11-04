CREATE TABLE IF NOT EXISTS "users" (
  "id" bigserial PRIMARY KEY,
  "uuid" char(36) UNIQUE,
  "created_at" timestamptz DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamptz DEFAULT CURRENT_TIMESTAMP,
  "deleted_at" timestamptz,
  "nickname" varchar(255),
  "email" varchar(255),
  "steam64_id" varchar(255),
  "steam32_id" varchar(255),
  "avatar" varchar(255),
  "rank" varchar(255),
  "plan" varchar(100) DEFAULT 'basic'
);