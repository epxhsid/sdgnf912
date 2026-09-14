-- migrate:up

CREATE TABLE outbox_events (
    id uuid PRIMARY KEY,
    event_type text NOT NULL,
    aggregate_id uuid NOT NULL,
    payload jsonb NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    published_at timestamptz,
    attempts integer NOT NULL DEFAULT 0,
    last_error text,
    next_retry_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX idx_outbox_events_processing ON outbox_events (next_retry_at) WHERE published_at IS NULL;

-- migrate:down

DROP TABLE outbox_events
