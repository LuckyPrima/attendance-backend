CREATE TABLE departments (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    departement_name VARCHAR(255) NOT NULL,
    max_clock_in TIME NOT NULL,
    max_clock_out TIME NOT NULL,
    created_at DATETIME(3),
    updated_at DATETIME(3)
);
