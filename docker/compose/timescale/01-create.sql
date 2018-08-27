CREATE EXTENSION IF NOT EXISTS timescaledb CASCADE;

CREATE TABLE speedy (
  time        TIMESTAMPTZ       NOT NULL,
  mac         MACADDR           NOT NULL,
  download    BIGINT            NOT NULL,
  upload      BIGINT            NOT NULL
);

CREATE TABLE speedy_metadata (
  mac         MACADDR           PRIMARY KEY,
  ipv4        INET              NULL,
  ipv6        INET              NULL
);

SELECT create_hypertable('speedy', 'time');

CREATE INDEX ON speedy (mac, time DESC);