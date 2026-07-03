CREATE TABLE IF NOT EXISTS payment_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id TEXT NOT NULL UNIQUE,
    type TEXT NOT NULL,
    payment_id TEXT NOT NULL,
    payload JSONB NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS payment_risks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    payment_id TEXT NOT NULL UNIQUE,
    event_id TEXT NOT NULL,
    score INT NOT NULL,
    level TEXT NOT NULL,
    reasons JSONB NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);