/* *****************************************************************************
// Setup the preferences
// ****************************************************************************/
SET NAMES utf8 COLLATE 'utf8_unicode_ci';
/* SET foreign_key_checks = 1; */
SET time_zone = '+00:00';
SET sql_mode = 'NO_AUTO_VALUE_ON_ZERO';
/*SET default_storage_engine = InnoDB; */
SET CHARACTER SET utf8;

/* *****************************************************************************
// Remove old database
// ****************************************************************************/
DROP DATABASE IF EXISTS thego;

/* *****************************************************************************
// Create new database
// ****************************************************************************/
CREATE DATABASE thego DEFAULT CHARSET = utf8 COLLATE = utf8_unicode_ci;
USE thego;

/* *****************************************************************************
// Create the tables
// ****************************************************************************/
create table instance(instance_id varchar(40), count int);

