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

CREATE TABLE "item" (
    "id" INTEGER NOT NULL,
    "name" TEXT NOT NULL,
    "cost_by_one" REAL NOT NULL,
    "one_is_called" TEXT NOT NULL,
    PRIMARY KEY ("id" AUTOINCREMENT)
);

CREATE TABLE "order_item_fulfillment" (
    "id" INTEGER NOT NULL,
    "order_id" INTEGER NOT NULL,
    "item_id" INTEGER NOT NULL,
    "quantity_fulfilled" REAL NOT NULL,
    PRIMARY KEY ("id" AUTOINCREMENT),
    FOREIGN KEY ("order_id") REFERENCES "order" ("id"),
    FOREIGN KEY ("item_id") REFERENCES "item" ("id"),
    UNIQUE ("order_id", "item_id")
);

CREATE TABLE "resource" (
    "id" INTEGER NOT NULL,
    "name" TEXT NOT NULL,
    "date_last_updated" TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "cost_by_one" REAL NOT NULL,
    "one_is_called" TEXT NOT NULL DEFAULT "Единица",
    "quantity" REAL NOT NULL DEFAULT 0,
    PRIMARY KEY ("id" AUTOINCREMENT),
    UNIQUE ("name", "cost_by_one")
);

CREATE TABLE "item_resource_need" (
    "id" INTEGER NOT NULL,
    "resource_id" INTEGER NOT NULL,
    "item_id" INTEGER NOT NULL,
    "quantity_needed" REAL NOT NULL,
    PRIMARY KEY ("id" AUTOINCREMENT),
    FOREIGN KEY ("item_id") REFERENCES "item" ("id"),
    FOREIGN KEY ("resource_id") REFERENCES "resource" ("id"),
    UNIQUE ("item_id", "resource_id")
);

CREATE TABLE "resource_resupply" (
    "id" INTEGER NOT NULL,
    "resource_id" INTEGER NOT NULL,
    "quantity_added" INTEGER NOT NULL,
    "date" TEXT NOT NULL,
    PRIMARY KEY ("id" AUTOINCREMENT),
    FOREIGN KEY ("resource_id") REFERENCES "resource" ("id")
);

CREATE TABLE "order_resource_spending" (
    "id" INTEGER NOT NULL,
    "order_id" INTEGER NOT NULL,
    "resource_id" INTEGER NOT NULL,
    "quantity_spent" REAL NOT NULL,
    "date" TEXT NOT NULL,
    PRIMARY KEY ("id" AUTOINCREMENT),
    FOREIGN KEY ("order_id") REFERENCES "order" ("id"),
    FOREIGN KEY ("resource_id") REFERENCES "resource" ("id"),
    UNIQUE ("order_id", "resource_id")
);

CREATE TABLE "user" (
    "id" INTEGER NOT NULL,
    "name" TEXT NOT NULL UNIQUE,
    "password" TEXT NOT NULL,
    "is_admin" INTEGER NOT NULL,
    PRIMARY KEY ("id" AUTOINCREMENT)
);
