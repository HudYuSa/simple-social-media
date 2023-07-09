CREATE Table
    IF NOT EXISTS "comments"(
        "id" UUID NOT NULL DEFAULT(uuid_generate_v4()),
        "user_id" UUID NOT NULL,
        "post_id" UUID NOT NULL,
        "content" TEXT NOT NULL,
        "created_at" timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
        -- constraint
        CONSTRAINT "comments_pkey" PRIMARY KEY ("id"),
        CONSTRAINT "fk_user" FOREIGN KEY("user_id") REFERENCES "users"("id") ON DELETE CASCADE,
        CONSTRAINT "fk_post" FOREIGN KEY("post_id") REFERENCES "posts"("id") ON DELETE CASCADE
    );