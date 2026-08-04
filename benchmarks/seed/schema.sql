-- The canonical products table. Every framework in this benchmark gets exactly
-- this, created from exactly this file.
--
-- Taken verbatim from what `grit generate resource Product` produces, because
-- the alternative — letting each ORM run its own migration — quietly changes
-- what is being measured. Eloquent, Django and GORM all disagree about integer
-- widths, timestamp precision and which columns get indexed, and a benchmark
-- where one side has an index the other lacks measures the index.

DROP TABLE IF EXISTS products;

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

-- Soft deletes are in the scaffold's model, so every framework filters on this
-- column and every framework gets the same index for it.
CREATE INDEX idx_products_deleted_at ON products USING btree (deleted_at);
