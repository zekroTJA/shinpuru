ALTER TABLE `reports` 
    ADD `attachment` text NOT NULL; 

ALTER TABLE `guilds`
    ADD `joinMsg` text NOT NULL,
    ADD `leaveMsg` text NOT NULL;