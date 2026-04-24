INSERT INTO persons (firstname,lastname,dateOfBirth,age)
VALUES ($1, $2, $3)
RETURNING id;