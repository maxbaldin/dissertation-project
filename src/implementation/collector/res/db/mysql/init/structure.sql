CREATE DATABASE IF NOT EXISTS collector;
CREATE USER 'collector'@'%' IDENTIFIED BY '!VB3{&uC6uwA9M#P';
GRANT ALL PRIVILEGES ON collector.* TO 'collector'@'%';


CREATE TABLE IF NOT EXISTS `inbound_traffic` (
             `date` date NOT NULL,
             `hour` int unsigned NOT NULL DEFAULT '0',
             `process_name` varchar(30) NOT NULL,
             `hostname` varchar(30) NOT NULL,
             `source_ip` varbinary(16) NOT NULL,
             `source_port` int unsigned NOT NULL,
             `target_ip` varbinary(16) NOT NULL,
             `target_port` int unsigned NOT NULL,
             `packets` int unsigned NOT NULL,
             `size` bigint unsigned DEFAULT NULL,
             UNIQUE KEY `inbound_traffic_pk` (`date`,`hour`,`process_name`,`hostname`,`source_ip`,`source_port`,`target_ip`,`target_port`),
             KEY `inbound_traffic_date_hour_index` (`date` DESC,`hour` DESC)
           ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `outbound_traffic` (
             `date` date NOT NULL,
             `hour` int unsigned NOT NULL DEFAULT '0',
             `process_name` varchar(30) NOT NULL,
             `hostname` varchar(30) NOT NULL,
             `source_ip` varbinary(16) NOT NULL,
             `source_port` int unsigned NOT NULL,
             `target_ip` varbinary(16) NOT NULL,
             `target_port` int unsigned NOT NULL,
             `packets` int unsigned NOT NULL,
             `size` bigint unsigned DEFAULT NULL,
             UNIQUE KEY `outbound_traffic_pk` (`date`,`hour`,`process_name`,`hostname`,`source_ip`,`source_port`,`target_ip`,`target_port`),
             KEY `outbound_traffic_date_hour_index` (`date` DESC,`hour` DESC)
           ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE OR REPLACE VIEW `known_nodes` AS
SELECT ip
FROM (
         SELECT inet_ntoa(source_ip) as ip
         FROM collector.outbound_traffic
         GROUP BY source_ip
         UNION
         SELECT inet_ntoa(target_ip) as ip
         FROM collector.inbound_traffic
         GROUP BY target_ip
     ) as nodes
GROUP BY ip;