
Below steps are necessary

1. create table imageinfo for storing data
 CREATE TABLE `imageinfo` (
        `id` INTEGER PRIMARY KEY AUTOINCREMENT,
        `title` VARCHAR(64) NULL,
        `filename` VARCHAR(64) NULL,
        `description` VARCHAR(64) NULL
    );
2.build and run