-- CREATE TABLE
CREATE TABLE IF NOT EXISTS foo (
    foo_id int primary key GENERATED ALWAYS AS IDENTITY,
    label varchar(32)
);

CREATE TABLE IF NOT EXISTS bar (
    bar_id int primary key GENERATED ALWAYS AS IDENTITY,
    label varchar(32),
    foo_id int references foo(foo_id)
);


-- INSERT DATA
INSERT INTO foo(label) VALUES
('foo1'),
('foo2'),
('foo3');

INSERT INTO bar(label, foo_id) VALUES
('bar1',1),
('bar2',1),
('bar3',1),
('bar4',1),
('bar5',1),
('bar6',2),
('bar7',2),
('bar8',2);
