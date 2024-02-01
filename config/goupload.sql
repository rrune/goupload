CREATE TABLE Users (
    Username varchar(255) NOT NULL UNIQUE,
    Password varchar(255) NOT NULL,
    Root tinyint(1) NOT NULL,
    Blind tinyint(1) NOT NULL,
    Onetime tinyint(1) NOT NULL,
    Restricted tinyint(1) NOT NULL
); 

CREATE TABLE Shorts (
    Short varchar(255) NOT NULL UNIQUE,
    Type varchar(255) NOT NULL,
    Author varchar(255) NOT NULL,
    Timestamp timestamp NOT NULL DEFAULT current_timestamp(),
    Ip varchar(255) NOT NULL DEFAULT '0.0.0.0',
    Restricted tinyint(1) NOT NULL,
    Downloads int(11) NOT NULL DEFAULT 0
)

CREATE TABLE Files (
    Short varchar(255) NOT NULL UNIQUE,
    Filename varchar(255) NOT NULL UNIQUE,
);

CREATE TABLE Texts (
    Short varchar(255) NOT NULL UNIQUE,
    Text text NOT NULL,
);