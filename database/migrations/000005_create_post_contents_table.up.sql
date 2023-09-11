CREATE TABLE
    "post_contents"(
        "id" UUID NOT NULL DEFAULT (uuid_generate_v4()),
        "post_id" UUID NOT NULL,
        "source" VARCHAR(200) NOT NULL,
        CONSTRAINT "post_contents_pkey" PRIMARY KEY ("id"),
        CONSTRAINT "fk_post" FOREIGN KEY("post_Id") REFERENCES "posts"("id") ON DELETE CASCADE
    )