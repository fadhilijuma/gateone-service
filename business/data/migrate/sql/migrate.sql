-- Version: 1.01
-- Description: Create table users
CREATE TABLE users
(
    user_id       UUID        NOT NULL,
    region_id     UUID        NOT NULL,
    name          TEXT        NOT NULL,
    email         TEXT UNIQUE NOT NULL,
    roles         TEXT[]      NOT NULL,
    password_hash TEXT        NOT NULL,
    department    TEXT        NULL,
    enabled       BOOLEAN     NOT NULL,
    date_created  TIMESTAMP   NOT NULL,
    date_updated  TIMESTAMP   NOT NULL,

    PRIMARY KEY (user_id),
    FOREIGN KEY (region_id) REFERENCES regions (region_id) ON DELETE CASCADE
);
-- Version: 1.02
-- Description: Create table patients
CREATE TABLE patients
(
    patient_id   UUID      NOT NULL,
    user_id      UUID      NOT NULL,
    name         TEXT      NOT NULL,
    age          INT       NOT NULL,
    video_links  TEXT[]    NOT NULL,
    condition    TEXT      NOT NULL,
    healed       BOOLEAN   NOT NULL,
    date_created TIMESTAMP NOT NULL,
    date_updated TIMESTAMP NOT NULL,

    PRIMARY KEY (patient_id),
    FOREIGN KEY (user_id) REFERENCES users (user_id) ON DELETE CASCADE
);

-- Version: 1.03
-- Description: Create table regions
CREATE TABLE regions
(
    region_id    UUID      NOT NULL,
    name         TEXT      NOT NULL,
    description  TEXT      NOT NULL,
    date_created TIMESTAMP NOT NULL,
    date_updated TIMESTAMP NOT NULL,

    PRIMARY KEY (region_id)
);


-- Version: 1.04
-- Description: Create table conditions
CREATE TABLE conditions
(
    condition_id UUID      NOT NULL,
    name         TEXT      NOT NULL,
    description  TEXT      NOT NULL,
    date_created TIMESTAMP NOT NULL,
    date_updated TIMESTAMP NOT NULL,

    PRIMARY KEY (condition_id)
);
-- Version: 1.05
-- Description: Create table roles
CREATE TABLE roles
(
    role_id      UUID      NOT NULL,
    name         TEXT      NOT NULL,
    description  TEXT      NOT NULL,
    date_created TIMESTAMP NOT NULL,
    date_updated TIMESTAMP NOT NULL,

    PRIMARY KEY (role_id)
);
