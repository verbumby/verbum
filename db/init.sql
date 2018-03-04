CREATE TABLE dicts (
    id INT PRIMARY KEY AUTO_INCREMENT,
    title VARCHAR(200)
);

CREATE TABLE articles (
    id INT PRIMARY KEY AUTO_INCREMENT,
    title varchar(200),
    content VARCHAR(10000),
    dict_id INT NOT NULL DEFAULT(-1)
);
