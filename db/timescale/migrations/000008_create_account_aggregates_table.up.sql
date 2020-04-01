CREATE TABLE IF NOT EXISTS account_aggregates
(
--     Domain entity
    id                                 uuid                     NOT NULL DEFAULT uuid_generate_v4(),
    created_at                         TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at                         TIMESTAMP WITH TIME ZONE NOT NULL,

-- Aggregate
    started_at_height                  NUMERIC                  NOT NULL,
    started_at                         TIMESTAMP WITH TIME ZONE NOT NULL,

-- Custom
    public_key                         TEXT,
    last_general_balance               DECIMAL(65, 0)           NOT NULL,
    last_general_nonce                 NUMERIC                  NOT NULL,
    last_escrow_active_balance         DECIMAL(65, 0)           NOT NULL,
    last_escrow_active_total_shares    DECIMAL(65, 0)           NOT NULL,
    last_escrow_debonding_balance      DECIMAL(65, 0)           NOT NULL,
    last_escrow_debonding_total_shares DECIMAL(65, 0)           NOT NULL,

    PRIMARY KEY (id)
);

-- Hypertable

-- Indexes
CREATE index idx_account_aggregates_public_key on account_aggregates (public_key);