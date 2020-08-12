BEGIN;

ALTER TABLE `guilds` 
    ADD COLUMN `colorReaction` text NOT NULL;

COMMIT;