-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
                                     id UUID PRIMARY KEY,
                                     email TEXT NOT NULL UNIQUE,
                                     password_hash TEXT NOT NULL DEFAULT '',
                                     role TEXT NOT NULL CHECK (role IN ('admin', 'user')),
                                     created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS rooms (
                                     id UUID PRIMARY KEY,
                                     name TEXT NOT NULL,
                                     description TEXT NULL,
                                     capacity INTEGER NULL CHECK (capacity IS NULL OR capacity >= 0),
                                     created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS schedules (
                                         id UUID PRIMARY KEY,
                                         room_id UUID NOT NULL UNIQUE REFERENCES rooms(id) ON DELETE CASCADE,
                                         days_of_week SMALLINT[] NOT NULL,
                                         start_time TIME NOT NULL,
                                         end_time TIME NOT NULL,
                                         created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS slots (
                                     id UUID PRIMARY KEY,
                                     room_id UUID NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
                                     start_at TIMESTAMPTZ NOT NULL,
                                     end_at TIMESTAMPTZ NOT NULL,
                                     created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
                                     CONSTRAINT slots_non_empty CHECK (end_at > start_at),
                                     CONSTRAINT slots_unique_interval UNIQUE (room_id, start_at, end_at)
);

CREATE TABLE IF NOT EXISTS bookings (
                                        id UUID PRIMARY KEY,
                                        slot_id UUID NOT NULL REFERENCES slots(id) ON DELETE CASCADE,
                                        user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                                        status TEXT NOT NULL CHECK (status IN ('active', 'cancelled')),
                                        conference_link TEXT NULL,
                                        created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX IF NOT EXISTS bookings_active_slot_uniq
    ON bookings (slot_id)
    WHERE status = 'active';

CREATE INDEX IF NOT EXISTS slots_room_start_idx
    ON slots (room_id, start_at);

CREATE INDEX IF NOT EXISTS bookings_user_status_idx
    ON bookings (user_id, status);


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS bookings;
DROP TABLE IF EXISTS slots;
DROP TABLE IF EXISTS schedules;
DROP TABLE IF EXISTS rooms;
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
