
CREATE DATABASE IF NOT EXISTS ddplanet;

-- create the users for each database
CREATE USER IF NOT EXISTS 'ddplanet'@'%' IDENTIFIED BY 'password';
GRANT CREATE, ALTER, INDEX, LOCK TABLES, REFERENCES, UPDATE, DELETE, DROP, SELECT, INSERT ON `ddplanet`.* TO 'ddplanet'@'%';

FLUSH PRIVILEGES;

USE ddplanet;


CREATE TABLE IF NOT EXISTS `support_requests` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `ts_created` BIGINT UNSIGNED NOT NULL,
    `first_name` VARCHAR(255) NOT NULL,
	`second_name` VARCHAR(255) NOT NULL,
    `last_name` VARCHAR(255) NOT NULL,
    `phone` VARCHAR(30) NOT NULL,
    `message` LONGTEXT NOT NULL,
    PRIMARY KEY (`id`),
    INDEX (`ts_created`)
) ENGINE = InnoDB CHARSET=utf8 COLLATE utf8_general_ci;
