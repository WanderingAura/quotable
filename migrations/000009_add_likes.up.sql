CREATE TABLE IF NOT EXISTS likes (
    user_id bigint NOT NULL REFERENCES users ON DELETE CASCADE,
    quote_id bigint NOT NULL REFERENCES quotes ON DELETE CASCADE,
    val smallint NOT NULL CHECK(val = 0 OR val = 1),
    PRIMARY KEY (user_id, quote_id)
);