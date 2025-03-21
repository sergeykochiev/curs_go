CREATE TABLE "order" (
    "id" SERIAL NOT NULL,
    "name" VARCHAR(255) NOT NULL,
    "client_name" VARCHAR(255) NOT NULL,
    "client_phone" VARCHAR(255) NOT NULL,
    "date_created" TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL,
    "creator_id" BIGINT NOT NULL,
    "date_ended" TIMESTAMP(0) WITHOUT TIME ZONE NULL,
    "ended" BOOLEAN NOT NULL DEFAULT 'false'
);

ALTER TABLE "order" ADD PRIMARY KEY ("id");

CREATE TABLE "resource" (
    "id" SERIAL NOT NULL,
    "name" VARCHAR(255) NOT NULL,
    "date_last_updated" TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL,
    "cost_by_one" BIGINT NOT NULL,
    "quantity" BIGINT NOT NULL
);

ALTER TABLE "resource" ADD PRIMARY KEY ("id");

ALTER TABLE "resource" ADD CONSTRAINT "resource_name_and_cost_by_one_unique" UNIQUE ("name", "cost_by_one");

CREATE TABLE "resource_resupply" (
    "id" SERIAL NOT NULL,
    "resource_id" BIGINT NOT NULL,
    "quantity_added" BIGINT NOT NULL,
    "date" TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL
);

ALTER TABLE "resource_resupply" ADD PRIMARY KEY ("id");

CREATE TABLE "resource_spending" (
    "id" SERIAL NOT NULL,
    "order_id" BIGINT NOT NULL,
    "resource_id" BIGINT NOT NULL,
    "quantity_spent" BIGINT NOT NULL,
    "date" TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL
);

ALTER TABLE "resource_spending" ADD PRIMARY KEY ("id");

CREATE TABLE "user" (
    "id" SERIAL NOT NULL,
    "name" VARCHAR(255) NOT NULL,
    "password" VARCHAR(255) NOT NULL,
    "is_admin" BOOLEAN NOT NULL
);

ALTER TABLE "user" ADD PRIMARY KEY ("id");

ALTER TABLE "user" ADD CONSTRAINT "user_name_unique" UNIQUE ("name");

ALTER TABLE "resource_spending" ADD CONSTRAINT "resource_spending_order_id_foreign" FOREIGN KEY ("order_id") REFERENCES "order" ("id");

ALTER TABLE "resource_spending" ADD CONSTRAINT "resource_spending_resource_id_foreign" FOREIGN KEY ("resource_id") REFERENCES "resource" ("id");

ALTER TABLE "order" ADD CONSTRAINT "order_creator_id_foreign" FOREIGN KEY ("creator_id") REFERENCES "user" ("id");

ALTER TABLE "resource_resupply" ADD CONSTRAINT "resource_resupply_resource_id_foreign" FOREIGN KEY ("resource_id") REFERENCES "resource" ("id");
