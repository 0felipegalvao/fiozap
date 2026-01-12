-- v1 -> v2: Create fzMessage table

CREATE TABLE IF NOT EXISTS "fzMessage" (
    "id" SERIAL PRIMARY KEY,
    "userId" VARCHAR(64) NOT NULL,
    "sessionId" VARCHAR(64),
    "chatJid" VARCHAR(255) NOT NULL,
    "senderJid" VARCHAR(255) NOT NULL,
    "messageId" VARCHAR(255) NOT NULL,
    "timestamp" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "messageType" VARCHAR(50) NOT NULL,
    "textContent" TEXT,
    "mediaLink" TEXT,
    "quotedMessageId" VARCHAR(255),
    UNIQUE("sessionId", "messageId")
);

CREATE INDEX IF NOT EXISTS "idxFzMessageSessionChat" 
ON "fzMessage" ("sessionId", "chatJid", "timestamp" DESC);
