-- DROP TABLE
DROP TABLE IF EXISTS bar;
DROP TABLE IF EXISTS foo;

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


-- INSERT DATA
INSERT INTO foo(foo_id, label, secret, value ,weight)
VALUES ('20000000-0000-0000-0000-000000000001', 'foo1', 'secret1', 1, 1.0),
       ('20000000-0000-0000-0000-000000000002', 'foo2', 'secret2', 2, 2.0),
       ('20000000-0000-0000-0000-000000000003', 'foo3', 'secret3', 3, 3.0);

INSERT INTO bar(bar_id, label, secret, value, foo_id)
VALUES ('20000000-0000-0000-0001-000000000001','bar1', 'secret1', 1, '20000000-0000-0000-0000-000000000001'),
       ('20000000-0000-0000-0001-000000000002','bar2', 'secret2', 2, '20000000-0000-0000-0000-000000000001'),
       ('20000000-0000-0000-0001-000000000003','bar3', 'secret3', 3, '20000000-0000-0000-0000-000000000001'),
       ('20000000-0000-0000-0001-000000000004','bar4', 'secret4', 4, '20000000-0000-0000-0000-000000000001'),
       ('20000000-0000-0000-0001-000000000005','bar5', 'secret5', 5, '20000000-0000-0000-0000-000000000001'),
       ('20000000-0000-0000-0001-000000000006','bar6', 'secret6', 6, '20000000-0000-0000-0000-000000000002');

