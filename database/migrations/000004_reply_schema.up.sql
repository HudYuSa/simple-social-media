CREATE TABLE
    IF NOT EXISTS "replies"(
        "id" UUID NOT NULL DEFAULT(uuid_generate_v4()),
        "user_id" UUID NOT NULL,
        "comment_id" UUID NULL,
        "mention_id" UUID NULL,
        "content" TEXT NOT NULL,
        "created_at" timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
        -- constraint
        CONSTRAINT "replies_pkey" PRIMARY KEY ("id"),
        CONSTRAINT "fk_user" FOREIGN KEY("user_id") REFERENCES "users"("id") ON DELETE CASCADE,
        CONSTRAINT "fk_parent_comment" FOREIGN KEY("comment_id") REFERENCES "comments"("id") ON DELETE CASCADE,
        CONSTRAINT "fk_mention" FOREIGN KEY("mention_id") REFERENCES "users"("id") ON DELETE
        SET NULL ON UPDATE
        SET NULL
    )