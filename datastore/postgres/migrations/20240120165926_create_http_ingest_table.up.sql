CREATE TABLE ingests (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    url_id uuid  NOT NULL REFERENCES urls(id),

    request jsonb NOT NULL DEFAULT '{}'::jsonb,

    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);
