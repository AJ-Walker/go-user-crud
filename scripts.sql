DROP TABLE IF EXISTS users;

CREATE TABLE users (
    id INT AUTO_INCREMENT NOT NULL,
    user_id VARCHAR(255) NOT NULL,
    name VARCHAR(128) NOT NULL,
    email VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    PRIMARY KEY (`id`)
);

INSERT INTO users (user_id,name,email,password)
VALUES 
('c1823e28-89f4-40e3-b0d1-714ad098ac58', 'Abhay Jha', 'abhay.jha@gmail.com','$2a$10$BzpOx9bhqXEA6IyPKmg3fegp4M39eSAX7u.u.vFkiEAoJHWfvsjJS'),
('f2bf6e5a-13ca-4090-be5b-ec12df6f9109', 'John Doe', 'john.doe@gmail.com','$2a$10$KCjYT4N.GgJWz1RtgOKAA.fa9JaQYssBco6lp0FnB.NIid5eUJGWq');

