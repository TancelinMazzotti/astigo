-- CREATE TABLE
CREATE TABLE IF NOT EXISTS foo
(
    foo_id UUID PRIMARY KEY,
    label  varchar(32) NOT NULL,
    secret varchar(32) NOT NULL,
    value  int NOT NULL DEFAULT 0,
    weight float NOT NULL DEFAULT 0.0,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz
);

CREATE TABLE IF NOT EXISTS bar
(
    bar_id UUID PRIMARY KEY,
    label  varchar(32) NOT NULL,
    secret varchar(32) NOT NULL,
    value  int NOT NULL DEFAULT 0,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz,
    foo_id UUID references foo (foo_id)
);
