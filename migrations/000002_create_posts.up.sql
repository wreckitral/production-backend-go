CREATE TABLE posts (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),

    -- If a user is deleted, their posts are deleted too.
    -- That keeps orphaned posts out of the database.
    author_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    title text NOT NULL,
    body text NOT NULL,

    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),

    CONSTRAINT posts_title_not_blank CHECK (length(trim(title)) > 0),
    CONSTRAINT posts_body_not_blank CHECK (length(trim(body)) > 0)
);

-- Speeds up queries like: list all posts by this author.
CREATE INDEX posts_author_id_idx ON posts(author_id);

-- Speeds up homepage/list queries ordered by newest first.
CREATE INDEX posts_created_at_idx ON posts(created_at DESC);
