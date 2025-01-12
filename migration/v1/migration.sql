SET timezone = 'UTC';

CREATE EXTENSION postgis;

-- Simple domains
CREATE DOMAIN non_empty_text AS VARCHAR(255) CHECK (VALUE <> ''); -- max length 255
CREATE DOMAIN non_empty_large_text AS VARCHAR(1000) CHECK (VALUE <> ''); -- max length 1000
CREATE DOMAIN email_domain AS TEXT CHECK (VALUE ~ '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$'); -- RFC 5322
CREATE DOMAIN phone_domain AS TEXT CHECK (VALUE ~ '^\+\d{5,15}$'); -- E.164

-- Address domain
CREATE TYPE address_type AS (
    street non_empty_text,
    street_number VARCHAR(10),
    extra non_empty_text,
    zip_code VARCHAR(10),
    city non_empty_text
);
CREATE DOMAIN address_domain AS address_type CHECK (VALUE IS NULL OR (
    (VALUE).street IS NOT NULL AND
    (VALUE).street_number ~ '^[0-9]+[A-Za-z]?(-[0-9]+[A-Za-z]?)?$' AND -- valid street number (e.g. 123, 123a, 123-125, 123a-125)
    (VALUE).zip_code ~ '^\d{5}$' AND -- 5 digit zip code
    (VALUE).city IS NOT NULL
));

CREATE TABLE users (
    id UUID PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
);
CREATE INDEX idx_users_deleted_at ON users(deleted_at);
CREATE UNIQUE INDEX idx_users_kinde_id ON users(kinde_id);

CREATE TABLE user_versions (
    id UUID PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    user_id UUID NOT NULL REFERENCES users(id) ON UPDATE RESTRICT ON DELETE RESTRICT,
    first_name non_empty_text NOT NULL,
    last_name non_empty_text NOT NULL,
    email email_domain NOT NULL,
    phone phone_domain NOT NULL,
    is_admin BOOLEAN NOT NULL DEFAULT FALSE,
);
CREATE INDEX idx_user_versions_created_at ON user_versions(created_at);
CREATE INDEX idx_user_versions_id ON user_versions(user_id);

-- Prevent updates and deletes on immutable tables
CREATE OR REPLACE FUNCTION func_no_update_or_delete() RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'UPDATE' THEN
        IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = TG_TABLE_NAME AND column_name = 'deleted_at') THEN
            IF OLD.deleted_at IS NULL AND NEW.deleted_at IS NOT NULL THEN
                RETURN NEW; -- allow soft delete
            END IF;
        END IF;
        RAISE EXCEPTION 'cannot update rows in this table';
    END IF;
    IF TG_OP = 'DELETE' THEN
        IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = TG_TABLE_NAME AND column_name = 'deleted_at') THEN
            IF OLD.deleted_at IS NOT NULL THEN
                RAISE EXCEPTION 'row is already deleted';
            ELSE
                EXECUTE format('UPDATE %I SET deleted_at = now() WHERE ctid = $1', TG_TABLE_NAME) USING OLD.ctid;
            END IF;
        ELSE
            RAISE EXCEPTION 'cannot delete rows in this table';
        END IF;
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;
DO $$
DECLARE
    table_name TEXT;
    tables TEXT[] := ARRAY[
        'users',
    ];
BEGIN
    FOREACH table_name IN ARRAY tables LOOP
        EXECUTE format('
            CREATE TRIGGER trig_no_update_or_delete_%s
            BEFORE UPDATE OR DELETE ON %s
            FOR EACH ROW
            EXECUTE FUNCTION func_no_update_or_delete()', table_name, table_name);
    END LOOP;
END;
$$ LANGUAGE plpgsql;
