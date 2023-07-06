CREATE TABLE
    "posts" (
        "id" uuid NOT NULL DEFAULT (uuid_generate_v4()),
        "user_id" UUID,
        "title" varchar NOT NULL,
        "photo" varchar NOT NULL,
        "content" TEXT NOT NULL,
        "created_at" timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
        "updated_at" timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
        CONSTRAINT "posts_pkey" PRIMARY KEY ("id"),
        CONSTRAINT "fk_user" FOREIGN KEY("user_id") REFERENCES "users"("id") ON DELETE CASCADE
    );