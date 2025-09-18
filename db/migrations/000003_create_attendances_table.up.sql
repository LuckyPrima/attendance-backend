CREATE TABLE attendances (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    employee_id BIGINT UNSIGNED NOT NULL,
    clock_in DATETIME(3) NULL,
    clock_out DATETIME(3) NULL,
    created_at DATETIME(3),
    updated_at DATETIME(3),
    CONSTRAINT fk_attendance_employee FOREIGN KEY (employee_id) REFERENCES employees(id) ON DELETE CASCADE
);