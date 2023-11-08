-- Create the "stock" database
CREATE DATABASE IF NOT EXISTS login;

-- Use the "stock" database
USE login;


CREATE TABLE IF NOT EXISTS user (
    Id INT NOT NULL PRIMARY KEY AUTO_INCREMENT,
    Name VARCHAR(255) NOT NULL,
    Password VARCHAR(255) NOT NULL,
    Email VARCHAR(255),
    UNIQUE(Email)
);