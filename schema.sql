CREATE TABLE "order" (
    "id" INTEGER NOT NULL,
    "name" TEXT NOT NULL,
    "client_name" TEXT NOT NULL,
    "client_phone" TEXT NOT NULL,
    "date_created" TEXT NOT NULL,
    "company_name" TEXT NULL,
    "creator_id" INTEGER NOT NULL,
    "date_ended" TEXT NULL,
    "ended" INTEGER NOT NULL DEFAULT 0,
    PRIMARY KEY ("id" AUTOINCREMENT),
    FOREIGN KEY ("creator_id") REFERENCES "user" ("id")
);

CREATE TABLE "resource" (
    "id" INTEGER NOT NULL,
    "name" TEXT NOT NULL,
    "date_last_updated" TEXT NOT NULL,
    "cost_by_one" REAL NOT NULL,
    "one_is_called" TEXT NOT NULL DEFAULT "Единица"
    "quantity" INTEGER NOT NULL,
    PRIMARY KEY ("id" AUTOINCREMENT),
    UNIQUE ("name", "cost_by_one")
);

CREATE TABLE "resource_resupply" (
    "id" INTEGER NOT NULL,
    "resource_id" INTEGER NOT NULL,
    "quantity_added" INTEGER NOT NULL,
    "date" TEXT NOT NULL,
    PRIMARY KEY ("id" AUTOINCREMENT),
    FOREIGN KEY ("resource_id") REFERENCES "resource" ("id")
);

CREATE TABLE "resource_spending" (
    "id" INTEGER NOT NULL,
    "order_id" INTEGER NOT NULL,
    "resource_id" INTEGER NOT NULL,
    "quantity_spent" INTEGER NOT NULL,
    "date" TEXT NOT NULL,
    PRIMARY KEY ("id" AUTOINCREMENT),
    FOREIGN KEY ("order_id") REFERENCES "order" ("id"),
    FOREIGN KEY ("resource_id") REFERENCES "resource" ("id")
);

CREATE TABLE "user" (
    "id" INTEGER NOT NULL,
    "name" TEXT NOT NULL UNIQUE,
    "password" TEXT NOT NULL,
    "is_admin" INTEGER NOT NULL,
    PRIMARY KEY ("id" AUTOINCREMENT)
);
