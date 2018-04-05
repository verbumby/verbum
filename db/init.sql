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

CREATE TABLE tasks (
    id INT PRIMARY KEY AUTO_INCREMENT,
    title varchar(200)
);

CREATE TABLE tasks_articles_rel (
    task_id INT,
    article_id INT,
    status varchar(20),
    PRIMARY KEY (task_id, article_id)
);
-- INSERT INTO tasks (id,title) VALUES (1,'RVBLR: Адфарматаваць і Праставіць Headwords');
-- INSERT INTO tasks_articles_rel (task_id, article_id, status) SELECT 1, articles.id, 'PENDING' FROM articles;
