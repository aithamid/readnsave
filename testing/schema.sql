CREATE TABLE Users (
    userId INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(100) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    valid BOOLEAN DEFAULT FALSE,
    public BOOLEAN DEFAULT TRUE
);

CREATE TABLE Books (
    bookId INT AUTO_INCREMENT PRIMARY KEY,
    bookname VARCHAR(100) NOT NULL,
    pages INT DEFAULT 0,
    totalpages INT NOT NULL,
    userId INT,
    FOREIGN KEY (userId) REFERENCES Users(userId)
);

CREATE TABLE Followers (
    id INT AUTO_INCREMENT PRIMARY KEY,
    userId1 INT,
    userId2 INT,
    approve BOOLEAN DEFAULT FALSE,
    FOREIGN KEY (userId1) REFERENCES Users(userId),
    FOREIGN KEY (userId2) REFERENCES Users(userId)
);

CREATE TABLE History (
    id INT AUTO_INCREMENT PRIMARY KEY,
    userId INT,
    bookId INT,
    pagesAdded INT NOT NULL,
    datetime DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (userId) REFERENCES Users(userId),
    FOREIGN KEY (bookId) REFERENCES Books(bookId)
);