BEGIN;

ALTER TABLE `guilds` 
    MODIFY `guildID` VARCHAR(25) NOT NULL,
    DROP `iid`,
    ADD PRIMARY KEY (`guildID`);

ALTER TABLE `permissions` 
    MODIFY `roleID` VARCHAR(25) NOT NULL,
    DROP `iid`,
    ADD PRIMARY KEY (`roleID`);

ALTER TABLE `reports` 
    MODIFY `id` VARCHAR(25) NOT NULL,
    DROP `iid`,
    ADD PRIMARY KEY (`id`);

ALTER TABLE `votes` 
    MODIFY `id` VARCHAR(25) NOT NULL,
    DROP `iid`,
    ADD PRIMARY KEY (`id`);

ALTER TABLE `tags` 
    MODIFY `id` VARCHAR(25) NOT NULL,
    DROP `iid`,
    ADD PRIMARY KEY (`id`);

ALTER TABLE `imagestore` 
    MODIFY `id` VARCHAR(25) NOT NULL,
    DROP `iid`,
    ADD PRIMARY KEY (`id`);

COMMIT;