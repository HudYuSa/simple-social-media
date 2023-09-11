CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE
    IF NOT EXISTS "users" (
        "id" uuid NOT NULL DEFAULT (uuid_generate_v4()),
        "username" varchar(50) NOT NULL,
        "email" varchar(50) NOT NULL,
        "password" varchar(250) NOT NULL,
        "profile_picture" varchar(250) NOT NULL,
        "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        "updated_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        CONSTRAINT "users_pkey" PRIMARY KEY ("id")
    );

CREATE UNIQUE INDEX IF NOT EXISTS "users_email_key" ON "users" ("email");