-- Create the "stock" database
DROP DATABASE stock;
CREATE DATABASE IF NOT EXISTS stock;

-- Use the "stock" database
USE stock;

-- Create the stock table
CREATE TABLE IF NOT EXISTS stocks (
    stock_id VARCHAR(25) NOT NULL PRIMARY KEY,
    stock_name NVARCHAR(255) NOT NULL,
    open_time DATETIME,
    open_price DECIMAL(10, 3),
    close_time DATETIME,
    close_price DECIMAL(10, 3),
    price DECIMAL(10, 3)
);

-- Create the users table
CREATE TABLE IF NOT EXISTS users (
    user_id INT AUTO_INCREMENT NOT NULL PRIMARY KEY,
    user_name VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    full_name VARCHAR(255),
    email VARCHAR(255) NOT NULL,
    phone_number VARCHAR(20),
    account_balance DECIMAL(10, 3) NOT NULL
);

-- Create the command_type table
CREATE TABLE IF NOT EXISTS command_types (
    command_type_id VARCHAR(25) NOT NULL PRIMARY KEY,
    command_type_name NVARCHAR(255) NOT NULL
);

-- Create the user_info table
CREATE TABLE IF NOT EXISTS user_infos (
    user_info_id INT AUTO_INCREMENT NOT NULL PRIMARY KEY,
    user_id INT NOT NULL,
    stock_id VARCHAR(25) NOT NULL,
    stock_name NVARCHAR(255) NOT NULL,
    stock_quantity float,
    FOREIGN KEY (user_id) REFERENCES users(user_id),
    FOREIGN KEY (stock_id) REFERENCES stocks(stock_id)
);

-- Create the `order_info` table
CREATE TABLE IF NOT EXISTS orders (
    order_id INT AUTO_INCREMENT NOT NULL PRIMARY KEY,
    price DECIMAL(10, 3),
    quantity float,
    total float,
    implementation_date DATETIME,
    stock_id VARCHAR(25) NOT NULL,
    command_type_id VARCHAR(25) NOT NULL,
    user_id INT NOT NULL,
    command nvarchar(25) NOT NULL,
    `status` VARCHAR(50) NOT NULL,
    FOREIGN KEY (stock_id) REFERENCES stocks(stock_id),
    FOREIGN KEY (command_type_id) REFERENCES command_types(command_type_id),
    FOREIGN KEY (user_id) REFERENCES users(user_id)
);

-- Insert data into users table
INSERT INTO users (user_name, password, full_name, email, phone_number, account_balance)
VALUES
    ('nguyenvanA', '$2a$12$9t8Z5NRaajzcXVVfbXkVeeDld0Kzxa4M7Nx8o1wnkkEQ2TUw2V7cu', 'Nguyen Van A', 'nguyenvana@email.com', '1234567890', 1000.00),
    ('tranthiB', 'pass1234', 'Tran Thi B', 'tranthib@email.com', '9876543210', 1500.00);

-- Insert data into stockOrder table
INSERT INTO command_types (command_type_id, command_type_name)
VALUES
    ('LO', N'Lệnh giới hạn.'),
    ('MP', N'Lệnh thị trường.'),
    ('MTL', N'Lệnh thị trường giới hạn.'),
    ('MOK', N'Lệnh thị trường khớp toàn bộ hoặc huỷ.'),
    ('MAK', N'Lệnh thị trường khớp và huỷ.'),
    ('ATO', N'Giá mở cửa'),
    ('ATC', N'Giá đóng cửa.'),
    ('PLO', N'Lệnh khớp lệnh sau giờ.'),
    ('TCO', N'Lệnh điều kiện với thời gian.'),
    ('PRO', N'Lệnh tranh mua hoặc tranh bán.'),
    ('ST', N'Lệnh dừng.'),
    ('TS', N'Lệnh xu hướng.');

-- Insert data into stock table
INSERT INTO stocks (stock_id, stock_name, open_time, open_price, close_time, close_price, price)
VALUES
    ('A32', N'Công ty Cổ phần 32 (A32)', '2023-10-16', 145.00, '2023-10-17', 148.50, 15.500),
    ('AAA', N'Công ty Cổ phần Nhựa An Phát Xanh (AAA)', '2023-10-16', 2700.00, '2023-10-17', 2725.50, 13.515),
    ('AAM', N'CTCP Thủy Sản MeKong', '2023-10-16', 290.00, '2023-10-17', 295.25, 8.123),
    ('AAS', N'CTCP Chứng khoán SmartInvest', '2023-10-16', 3400.00, '2023-10-17', 3440.75, 10.187),
    ('AAT', N'CTCP Tập đoàn Tiên Sơn Thanh Hóa', '2023-10-16', 750.00, '2023-10-17', 762.50, 21.312);

-- Insert data into user_info table
INSERT INTO user_infos (user_id, stock_id, stock_name, stock_quantity)
VALUES
    (1, 'A32', N'Công ty Cổ phần 32 (A32)', 100),
    (1, 'AAA', N'Công ty Cổ phần Nhựa An Phát Xanh (AAA)', 50),
    (2, 'AAM', N'CTCP Thủy Sản MeKong', 75),
    (2, 'AAS', N'CTCP Chứng khoán SmartInvest', 30);

-- Insert data into `order_info` table
INSERT INTO orders (stock_id, price, quantity, total, implementation_date, command_type_id, user_id, `command`, status)
VALUES
    ('A32', 147, 10, 1470, '2023-10-18', 'LO', 1, N'MUA', 'SUCCESS'),
    ('AAA', 2730, 5, 13650, '2023-10-18', 'MP', 1, N'BÁN', 'SUCCESS'),
    ('AAM', 292, 15, 4380, '2023-10-18', 'MTL', 2, N'MUA', 'CANCEL'),
    ('AAS', 3420, 8, 27360, '2023-10-18', 'MAK', 2, N'BÁN', 'CANCEL');

-- Select orders with status 'CANCEL'
SELECT * FROM orders WHERE status = 'CANCEL';

-- Select user_id for the user with userName 'nguyenvanA' and password 'password123'
SELECT user_id FROM users WHERE user_name = 'nguyenvanA' AND password = 'password123';
