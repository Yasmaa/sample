CREATE TABLE IF NOT EXISTS users(
    id int PRIMARY KEY,
    name VARCHAR (50) NOT NULL,
    email VARCHAR (300) UNIQUE NOT NULL,
    age int not null
    );


