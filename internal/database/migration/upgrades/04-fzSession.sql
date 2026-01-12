-- v3 -> v4: Create fzSession table for multi-tenant support

-- Enable pgcrypto extension for gen_random_bytes
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- Add maxSessions and createdAt to fzUser
ALTER TABLE "fzUser" ADD COLUMN IF NOT EXISTS "maxSessions" INTEGER DEFAULT 5;
ALTER TABLE "fzUser" ADD COLUMN IF NOT EXISTS "createdAt" TIMESTAMP DEFAULT CURRENT_TIMESTAMP;

-- Create fzSession table
CREATE TABLE IF NOT EXISTS "fzSession" (
    "id" VARCHAR(64) PRIMARY KEY,
    "userId" VARCHAR(64) NOT NULL REFERENCES "fzUser"("id") ON DELETE CASCADE,
    "name" VARCHAR(255) NOT NULL,
    "jid" VARCHAR(255) DEFAULT '',
    "qrCode" TEXT DEFAULT '',
    "connected" INTEGER DEFAULT 0,
    "webhook" TEXT DEFAULT '',
    "events" TEXT DEFAULT '',
    "proxyUrl" TEXT DEFAULT '',
    "createdAt" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE("userId", "name")
);

CREATE INDEX IF NOT EXISTS "idxFzSessionUserId" ON "fzSession" ("userId");
CREATE INDEX IF NOT EXISTS "idxFzSessionConnected" ON "fzSession" ("connected");

-- Migrate existing user data to sessions (before dropping columns)
INSERT INTO "fzSession" ("id", "userId", "name", "jid", "qrCode", "connected", "webhook", "events", "proxyUrl")
SELECT 
    encode(gen_random_bytes(16), 'hex'),
    "id",
    'default',
    COALESCE("jid", ''),
    COALESCE("qrCode", ''),
    COALESCE("connected", 0),
    COALESCE("webhook", ''),
    COALESCE("events", ''),
    COALESCE("proxyUrl", '')
FROM "fzUser"
WHERE "jid" IS NOT NULL AND "jid" != ''
ON CONFLICT ("userId", "name") DO NOTHING;

-- Remove columns from fzUser that are now in fzSession
ALTER TABLE "fzUser" DROP COLUMN IF EXISTS "jid";
ALTER TABLE "fzUser" DROP COLUMN IF EXISTS "qrCode";
ALTER TABLE "fzUser" DROP COLUMN IF EXISTS "connected";
ALTER TABLE "fzUser" DROP COLUMN IF EXISTS "events";
ALTER TABLE "fzUser" DROP COLUMN IF EXISTS "proxyUrl";
ALTER TABLE "fzUser" DROP COLUMN IF EXISTS "expiration";
ALTER TABLE "fzUser" DROP COLUMN IF EXISTS "webhook";
