CREATE DATABASE keylogger;
\c keylogger
CREATE TABLE IF NOT EXISTS recordings (
  ip_address VARCHAR(12) NOT NULL, 
  time_stamp TIMESTAMP NOT NULL, 
  keystrokes VARCHAR(50)
);