CREATE TABLE IF NOT EXISTS "review_requests" (
  "id" bigserial,
  "uuid" char(36) UNIQUE,
  "author_id" bigserial,
  "created_at" timestamptz DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamptz DEFAULT CURRENT_TIMESTAMP,
  "deleted_at" timestamptz,
  "description" varchar(255) NOT NULL,
  "match_id" varchar(255) NOT NULL,
  "hero_played" integer NOT NULL,
  "author_rank" varchar(255),
  "self_rate_laning" integer NOT NULL,
  "self_rate_teamfights" integer NOT NULL,
  "self_rate_overall" integer NOT NULL,
  "state" varchar(100) DEFAULT 'open',
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_review_requests_user" FOREIGN KEY ("author_id") REFERENCES "users"("id") ON DELETE CASCADE
);