
-- 用户信息
CREATE TABLE IF NOT EXISTS `t_wowoo_user_info` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '自增ID',
  `uuid` char(32) NOT NULL COMMENT 'UUID',
  `openid` varchar(64) NOT NULL COMMENT '微信OpenID',
  `unionid` varchar(64) NOT NULL DEFAULT '' COMMENT '微信UnionID',
  `session_key` varchar(64) NOT NULL DEFAULT '' COMMENT 'session_key',
  `status` tinyint(4) NOT NULL DEFAULT '0' COMMENT '0:正常；1:被删除；2:被禁用',
  `nickname` varchar(255) NOT NULL DEFAULT '' COMMENT '昵称',
  `avatar_url` varchar(255) NOT NULL DEFAULT '' COMMENT '头像URL',
  `mobile` varchar(32) NOT NULL DEFAULT '' COMMENT '用户手机号码',
  `inviter` varchar(255) NOT NULL DEFAULT '' COMMENT '邀请者',
  `gender` int(11) NOT NULL DEFAULT '0' COMMENT '性别：0-未知；1-男；2-女',
  `language` varchar(16) NOT NULL DEFAULT '' COMMENT '语言',
  `city` varchar(32) NOT NULL DEFAULT '' COMMENT '城市',
  `province` varchar(32) NOT NULL DEFAULT '' COMMENT '省份',
  `country` varchar(32) NOT NULL DEFAULT '' COMMENT '国家',
  `create_time` datetime NOT NULL COMMENT '记录创建时间',
  `last_login_time` datetime NOT NULL COMMENT '最后登录时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_user_uuid` (`uuid`),
  UNIQUE KEY `uniq_user_openid` (`openid`),
  KEY `uniq_user_unionid` (`unionid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户信息';

-- 签到表
CREATE TABLE IF NOT EXISTS `t_wowoo_attendance` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `uuid` varchar(32) NOT NULL DEFAULT '' COMMENT '用户uuid',
  `attendance_time` datetime DEFAULT NULL COMMENT '签到时间',
  `seq_day` int(11) NOT NULL DEFAULT '1' COMMENT '连续签到天数,超过七天重置',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `uniq_uuid` (`uuid`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='签到表';


-- 用户关系
CREATE TABLE IF NOT EXISTS `t_wowoo_user_relation` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `uid` char(32) NOT NULL DEFAULT '' COMMENT '用户id',
  `parent_id` char(32) NOT NULL DEFAULT '' COMMENT '一级邀请人id',
  `grand_parent_id` char(32) NOT NULL DEFAULT '' COMMENT '二级邀请人id',
  `create_time` datetime NOT NULL COMMENT '创建时间',
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `user_relation_wowoo_uniq_uid` (`uid`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户关系';

