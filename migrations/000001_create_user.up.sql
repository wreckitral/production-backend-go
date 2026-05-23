CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE users (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),

    email text NOT NULL UNIQUE,

    password_hash text NOT NULL,

    created_at timestamptz NOT NULL DEFAULT now(),

    CONSTRAINT users_email_lowercase CHECK (email = lower(email)),
    CONSTRAINT users_email_not_blank CHECK (length(trim(email)) > 0)
);
