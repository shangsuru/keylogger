#!/bin/bash
psql postgresql://postgres:password@localhost:5432/ -c "CREATE DATABASE keylogger"
psql postgresql://postgres:password@localhost:5432/keylogger -c "CREATE TABLE IF NOT EXISTS recordings (
                                                                  ip_address VARCHAR(12) NOT NULL, 
                                                                  time_stamp TIMESTAMP NOT NULL, 
                                                                  keystrokes VARCHAR(50));"