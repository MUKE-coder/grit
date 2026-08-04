-- Encore requires a migration to consider the database defined. The table is
-- created here with exactly the shape in ../../seed/schema.sql, which every
-- other framework in this benchmark shares — same column types, same index.
-- The harness truncates and reloads the same 10,000 rows before every run
-- regardless, so this only has to agree, not populate.
CREATE TABLE products (
    id          character varying(36) NOT NULL,
    name        character varying(255),
    sku         character varying(255),
    description text,
    price       numeric(12,2),
    stock       bigint,
    active      boolean,
    version     bigint DEFAULT 1 NOT NULL,
    created_at  timestamp with time zone,
    updated_at  timestamp with time zone,
    deleted_at  timestamp with time zone
);

ALTER TABLE ONLY products ADD CONSTRAINT products_pkey PRIMARY KEY (id);
CREATE INDEX idx_products_deleted_at ON products USING btree (deleted_at);
