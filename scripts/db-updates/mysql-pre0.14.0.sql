ALTER TABLE `permissions`
	DROP COLUMN `permission`,
    ADD COLUMN `permission` text NOT NULL;