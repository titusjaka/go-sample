-- +migrate Up
CREATE TABLE IF NOT EXISTS test ( id int );

INSERT INTO test (id) VALUES (1);

-- +migrate Down
DROP TABLE test;
