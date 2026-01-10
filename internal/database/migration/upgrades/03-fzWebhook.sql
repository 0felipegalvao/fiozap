-- v2 -> v3: Create fzWebhook table

CREATE TABLE IF NOT EXISTS "fzWebhook" (
    "id" SERIAL PRIMARY KEY,
    "userId" VARCHAR(64) NOT NULL,
    "eventType" VARCHAR(50) NOT NULL,
    "payload" JSONB NOT NULL,
    "status" VARCHAR(20) DEFAULT 'pending',
    "attempts" INTEGER DEFAULT 0,
    "lastAttempt" TIMESTAMP,
    "createdAt" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS "idxFzWebhookPending" 
ON "fzWebhook" ("status", "createdAt") WHERE "status" = 'pending';
