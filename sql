CREATE TABLE `rikxian` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `code` varchar(8) NOT NULL,
  `date` date NOT NULL,
  `date_int` int(11) NOT NULL,
  `kaipan` double(8,4) NOT NULL,
  `shoupan` double(8,4) NOT NULL,
  `zuigao` double(8,4) NOT NULL,
  `zuidi` double(8,4) NOT NULL,
  `zhangdiee` double(8,4) NOT NULL,
  `zhangdiefu` double(8,4) NOT NULL,
  `chengjiaoliang` int(11) NOT NULL,
  `chengjiaoe` double(18,4) NOT NULL,
  `huanshoulv` double(8,4) NOT NULL,
  `zongshizhi` bigint(20) NOT NULL,
  `liutongshizhi` bigint(20) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `code` (`code`,`date_int`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `baseinfo` (
  `id` int(5) NOT NULL AUTO_INCREMENT,
  `code` varchar(10) NOT NULL,
  `name` varchar(50) NOT NULL,
  `jiaoyisuo` varchar(10) NOT NULL,
  `a_or_b` enum('B','A') NOT NULL,
  `market_time` int(11) NOT NULL,
  `zong_gu_ben` double(30,2) NOT NULL,
  `liutong_gu_ben` double(30,2) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `a_or_b` (`a_or_b`,`code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;



CREATE TABLE `chengjiaoliang` (
`id`  int(11) NOT NULL AUTO_INCREMENT ,
`code`  varchar(8) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' ,
`date`  varchar(14) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' ,
`chaoda_buy`  int(11) NOT NULL DEFAULT 0 ,
`chaoda_buy_shou`  int(11) NOT NULL DEFAULT 0 ,
`chaoda_sall`  int(11) NOT NULL DEFAULT 0 ,
`chaoda_sall_shou`  int(11) NOT NULL DEFAULT 0 ,
`da_buy`  int(11) NOT NULL DEFAULT 0 ,
`da_buy_shou`  int(11) NOT NULL DEFAULT 0 ,
`da_sall`  int(11) NOT NULL DEFAULT 0 ,
`da_sall_shou`  int(11) NOT NULL DEFAULT 0 ,
`zhong_buy`  int(11) NOT NULL DEFAULT 0 ,
`zhong_buy_shou`  int(11) NOT NULL DEFAULT 0 ,
`zhong_sall`  int(11) NOT NULL DEFAULT 0 ,
`zhong_sall_shou`  int(11) NOT NULL DEFAULT 0 ,
`xiao_buy`  int(11) NOT NULL DEFAULT 0 ,
`xiao_buy_shou`  int(11) NOT NULL DEFAULT 0 ,
`xiao_sall`  int(11) NOT NULL DEFAULT 0 ,
`xiao_sall_shou`  int(11) NOT NULL DEFAULT 0 ,
`zhangdiefu`  double(8,4) NOT NULL DEFAULT 0.0000 ,
PRIMARY KEY (`id`),
INDEX `code` (`code`, `date`) USING BTREE
)
ENGINE=InnoDB
DEFAULT CHARACTER SET=utf8 COLLATE=utf8_general_ci
AUTO_INCREMENT=1664
ROW_FORMAT=DYNAMIC
;
