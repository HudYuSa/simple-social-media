CREATE Table
    IF NOT EXISTS "comments"(
        "id" UUID NOT NULL DEFAULT(uuid_generate_v4()),
        "user_id" UUID REFERENCES "users"("id"),
        "parent_comment_id" UUID REFERENCES "comments"("id"),
        "post_id" UUID REFERENCES "posts"("id"),
        "comment" TEXT NOT NULL,
        "created_at" timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
        -- constraint
        CONSTRAINT "comments_pkey" PRIMARY KEY ("id"),
        CONSTRAINT "fk_user" FOREIGN KEY("user_id") REFERENCES "users"("id") ON DELETE CASCADE,
        CONSTRAINT "fk_post" FOREIGN KEY("post_id") REFERENCES "posts"("id") ON DELETE CASCADE,
        CONSTRAINT "fk_parent_comment" FOREIGN KEY("parent_comment_id") REFERENCES "comments"("id") ON DELETE CASCADE
    );