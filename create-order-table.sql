CREATE TABLE IF NOT EXISTS `orders` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'Order ID',
  `distance` BIGINT UNSIGNED NOT NULL COMMENT 'order distance from origin to destination. Unit is in meters',
  `status` varchar(8) NOT NULL DEFAULT 'unassign' COMMENT 'Status of the order either unassign or taken.',
  `create_date` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Date when the order was created',
  `update_date` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Date when the order was last updated',
  PRIMARY KEY (`id`),
  UNIQUE INDEX `id` (`id`)
) AUTO_INCREMENT = 40000000
COMMENT = 'Order information';