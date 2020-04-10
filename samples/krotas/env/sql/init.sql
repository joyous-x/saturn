
CREATE TABLE `t_mytest_table` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `uid` varchar(40) NOT NULL DEFAULT '' COMMENT 'uid',
  `content` text NOT NULL COMMENT 'content',
  `status` int(10) unsigned NOT NULL DEFAULT '0' COMMENT 'status',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后一次更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_uid_status` (`uid`,`status`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COMMENT='t_mytest_table';