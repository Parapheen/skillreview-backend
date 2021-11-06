CREATE TABLE IF NOT EXISTS "bases" (
  "id" bigserial,
  "uuid" char(36) UNIQUE,
  "created_at" timestamptz DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamptz DEFAULT CURRENT_TIMESTAMP,
  "deleted_at" timestamptz,
  PRIMARY KEY ("id")
)