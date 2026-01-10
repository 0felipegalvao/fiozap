-- v0 -> v1: Create fzUser table

CREATE TABLE IF NOT EXISTS "fzUser" (
    "id" VARCHAR(64) PRIMARY KEY,
    "name" VARCHAR(255) NOT NULL,
    "token" VARCHAR(255) NOT NULL UNIQUE,
    "webhook" TEXT DEFAULT '',
    "jid" VARCHAR(255) DEFAULT '',
    "qrCode" TEXT DEFAULT '',
    "connected" INTEGER DEFAULT 0,
    "expiration" BIGINT DEFAULT 0,
    "events" TEXT DEFAULT '',
    "proxyUrl" TEXT DEFAULT ''
);
