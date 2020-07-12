create table if not exists users(
  id varchar(36) NOT NULL PRIMARY KEY,
  name varchar(255) NOT NULL,
  email varchar(255) NOT NULL,
  password varchar(255) NOT NULL,
  avatar varchar(255) NULL,
  UNIQUE KEY email (email)
)