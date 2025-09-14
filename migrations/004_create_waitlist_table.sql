-- +goose Up
CREATE TABLE IF NOT EXISTS waitlist (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,
    event_id VARCHAR(36) NOT NULL,
    quantity INTEGER NOT NULL CHECK (quantity > 0),
    priority INTEGER NOT NULL DEFAULT 0,
    status VARCHAR(20) NOT NULL CHECK (status IN ('active', 'notified', 'expired', 'converted')) DEFAULT 'active',
    joined_at TIMESTAMP NOT NULL DEFAULT NOW(),
    notified_at TIMESTAMP,
    expires_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    
    FOREIGN KEY (event_id) REFERENCES events(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    
    -- Ensure a user can only be on waitlist once per event
    UNIQUE(user_id, event_id)
);

-- Indexes for efficient querying
CREATE INDEX idx_waitlist_event_id ON waitlist(event_id);
CREATE INDEX idx_waitlist_user_id ON waitlist(user_id);
CREATE INDEX idx_waitlist_status ON waitlist(status);
CREATE INDEX idx_waitlist_priority_joined ON waitlist(event_id, priority DESC, joined_at ASC) WHERE status = 'active';

-- +goose Down
DROP TABLE IF EXISTS waitlist;
