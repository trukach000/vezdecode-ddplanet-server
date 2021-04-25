ALTER TABLE `support_requests`
    ADD COLUMN `status` ENUM('new', 'closed') NOT NULL DEFAULT 'new',
    ADD COLUMN `ts_closed` BIGINT UNSIGNED NOT NULL DEFAULT 0,
    ADD INDEX(`status`);