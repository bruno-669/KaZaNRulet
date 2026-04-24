CREATE TABLE IF NOT EXISTS persons (
    id SERIAL PRIMARY KEY,
    firstname TEXT NOT NULL,
    lastname TEXT NOT NULL,
    dateOfBirth DATE NOT NULL,
    age INT NOT NULL,
    weight INT,
    height INT,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS doctors (
    id SERIAL PRIMARY KEY,
    person_id INT NOT NULL,
    specialty TEXT NOT NULL,
    startOfWork DATE NOT NULL,
    experience INT,
    created_at TIMESTAMP DEFAULT NOW(),
    FOREIGN KEY (person_id) REFERENCES persons (id)
    ON DELETE RESTRICT ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS patients (
    id SERIAL PRIMARY KEY,
    person_id INT NOT NULL,
    condition TEXT NOT NULL DEFAULT 'healthy',
    diagnosis TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    FOREIGN KEY (person_id) REFERENCES persons (id)
        ON DELETE RESTRICT ON UPDATE CASCADE
);