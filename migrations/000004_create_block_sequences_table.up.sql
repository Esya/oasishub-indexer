CREATE TABLE IF NOT EXISTS block_sequences
(
    id                  BIGSERIAL                NOT NULL,
    created_at          TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at          TIMESTAMP WITH TIME ZONE NOT NULL,

    height              DECIMAL(65, 0)           NOT NULL,
    time                TIMESTAMP WITH TIME ZONE NOT NULL,

    hash                TEXT                     NOT NULL,
    proposer_entity_uid TEXT                     NOT NULL,
    transactions_count  INT,
    app_version         BIGINT                   NOT NULL,
    block_version       BIGINT                   NOT NULL,

    PRIMARY KEY (time, id)
);

-- Hypertable
SELECT create_hypertable('block_sequences', 'time', if_not_exists => TRUE);

-- Indexes
CREATE index idx_block_sequences_height on block_sequences (height, time DESC);
CREATE index idx_block_sequences_app_version on block_sequences (app_version, time DESC);
CREATE index idx_block_sequences_block_version on block_sequences (block_version, time DESC);
CREATE index idx_block_sequences_proposer_hash on block_sequences (hash, time DESC);
CREATE index idx_block_sequences_proposer_entity_uid on block_sequences (proposer_entity_uid, time DESC);