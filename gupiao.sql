/*
Navicat MySQL Data Transfer

Source Server         : localhost
Source Server Version : 50717
Source Host           : localhost:3306
Source Database       : gupiao

Target Server Type    : MYSQL
Target Server Version : 50717
File Encoding         : 65001

Date: 2017-01-20 15:42:07
*/

SET FOREIGN_KEY_CHECKS=0;

-- ----------------------------
-- Table structure for baseinfo
-- ----------------------------
DROP TABLE IF EXISTS `baseinfo`;
CREATE TABLE `baseinfo` (
  `id` int(5) NOT NULL AUTO_INCREMENT,
  `code` varchar(10) NOT NULL DEFAULT '',
  `name` varchar(50) NOT NULL DEFAULT '',
  `jiaoyisuo` varchar(10) NOT NULL DEFAULT '',
  `a_or_b` enum('B','A') NOT NULL,
  `market_time` int(11) NOT NULL DEFAULT '0',
  `zong_gu_ben` double(30,2) NOT NULL DEFAULT '0.00',
  `liutong_gu_ben` double(30,2) NOT NULL DEFAULT '0.00',
  PRIMARY KEY (`id`),
  KEY `a_or_b` (`a_or_b`,`code`)
) ENGINE=InnoDB AUTO_INCREMENT=3136 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for rikxian
-- ----------------------------
DROP TABLE IF EXISTS `rikxian`;
CREATE TABLE `rikxian` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `code` varchar(8) NOT NULL DEFAULT '',
  `date` date NOT NULL,
  `date_int` int(11) NOT NULL DEFAULT '0',
  `kaipan` double(8,4) NOT NULL DEFAULT '0.0000',
  `shoupan` double(8,4) NOT NULL DEFAULT '0.0000',
  `zuigao` double(8,4) NOT NULL DEFAULT '0.0000',
  `zuidi` double(8,4) NOT NULL DEFAULT '0.0000',
  `zhangdiee` double(8,4) NOT NULL DEFAULT '0.0000',
  `zhangdiefu` double(8,4) NOT NULL DEFAULT '0.0000',
  `chengjiaoliang` int(11) NOT NULL DEFAULT '0',
  `chengjiaoe` double(18,4) NOT NULL DEFAULT '0.0000',
  `huanshoulv` double(8,4) NOT NULL DEFAULT '0.0000',
  `zongshizhi` bigint(20) NOT NULL DEFAULT '0',
  `liutongshizhi` bigint(20) NOT NULL DEFAULT '0',
  `rijun30` double(8,4) NOT NULL DEFAULT '0.0000',
  `rijun30_cha` double(8,4) NOT NULL DEFAULT '0.0000',
  PRIMARY KEY (`id`),
  KEY `code` (`code`,`date_int`)
) ENGINE=InnoDB AUTO_INCREMENT=1251323 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for rimingxi
-- ----------------------------
DROP TABLE IF EXISTS `rimingxi`;
CREATE TABLE `rimingxi` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` varchar(10) NOT NULL DEFAULT '',
  `date` varchar(14) NOT NULL DEFAULT '',
  `date_int` int(11) NOT NULL DEFAULT '0',
  `chengjiaojia` double(8,4) NOT NULL DEFAULT '0.0000',
  `zhengdiee` double(8,4) NOT NULL DEFAULT '0.0000',
  `chengjiaoshou` int(11) NOT NULL DEFAULT '0',
  `chengjiaoe` double(20,4) NOT NULL DEFAULT '0.0000',
  `buy_sall` varchar(2) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`),
  KEY `code` (`code`,`date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for screen
-- ----------------------------
DROP TABLE IF EXISTS `screen`;
CREATE TABLE `screen` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `Zhangdiefu` double(8,4) NOT NULL DEFAULT '0.0000',
  `rijun30_cha` double(8,4) NOT NULL DEFAULT '0.0000',
  `hit_num` int(11) NOT NULL DEFAULT '0',
  `total_num` int(11) NOT NULL DEFAULT '0',
  `percent` double(8,4) NOT NULL DEFAULT '0.0000',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
SET FOREIGN_KEY_CHECKS=1;
