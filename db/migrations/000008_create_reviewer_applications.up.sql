CREATE TABLE IF NOT EXISTS "reviewer_applications" (
  "id" bigserial,
  "uuid" char(36) UNIQUE,
  "author_id" bigserial,
  "created_at" timestamptz DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamptz DEFAULT CURRENT_TIMESTAMP,
  "deleted_at" timestamptz,
  "description" varchar(255) NOT NULL,
  "rating" integer NOT NULL,
  "state" varchar(100) DEFAULT 'recieved',
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_reviewer_applications_user" FOREIGN KEY ("author_id") REFERENCES "users"("id")
);
