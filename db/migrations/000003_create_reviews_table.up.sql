CREATE TABLE IF NOT EXISTS "reviews" (
  "id" bigserial,
  "uuid" char(36) UNIQUE,
  "author_uuid" char(36) UNIQUE,
  "author_id" bigserial,
  "review_request_uuid" char(36) UNIQUE,
  "review_request_id" bigserial,
  "created_at" timestamptz DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamptz DEFAULT CURRENT_TIMESTAMP,
  "deleted_at" timestamptz,
  "description" varchar(255) NOT NULL,
  "rate_laning" integer NOT NULL,
  "rate_teamfights" integer NOT NULL,
  "rate_overall" integer NOT NULL,
  "state" varchar(100) DEFAULT 'submitted',
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_reviews_user" FOREIGN KEY ("author_id") REFERENCES "users"("id"),
  CONSTRAINT "fk_reviews_review_request" FOREIGN KEY ("review_request_id") REFERENCES "review_requests"("id") ON DELETE CASCADE
);