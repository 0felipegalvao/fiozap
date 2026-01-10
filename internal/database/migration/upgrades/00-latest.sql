-- v0 -> v2 (compatible with v1+): Latest schema for fresh installs

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

CREATE TABLE IF NOT EXISTS "fzMessageHistory" (
    "id" SERIAL PRIMARY KEY,
    "userId" VARCHAR(64) NOT NULL,
    "chatJid" VARCHAR(255) NOT NULL,
    "senderJid" VARCHAR(255) NOT NULL,
    "messageId" VARCHAR(255) NOT NULL,
    "timestamp" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "messageType" VARCHAR(50) NOT NULL,
    "textContent" TEXT,
    "mediaLink" TEXT,
    "quotedMessageId" VARCHAR(255),
    UNIQUE("userId", "messageId")
);

CREATE INDEX IF NOT EXISTS "idxFzMessageHistoryUserChat" 
ON "fzMessageHistory" ("userId", "chatJid", "timestamp" DESC);
