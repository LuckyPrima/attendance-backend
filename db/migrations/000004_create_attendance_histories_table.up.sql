CREATE TABLE attendance_histories (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    employee_id BIGINT UNSIGNED NOT NULL,
    attendance_id BIGINT UNSIGNED NOT NULL,
    date_attendance DATETIME(3) NOT NULL,
    attendance_type TINYINT NOT NULL COMMENT '1 = Clock In, 2 = Clock Out',
    description VARCHAR(255),
    created_at DATETIME(3),
    updated_at DATETIME(3),
    CONSTRAINT fk_history_employee FOREIGN KEY (employee_id) REFERENCES employees(id) ON DELETE CASCADE,
    CONSTRAINT fk_history_attendance FOREIGN KEY (attendance_id) REFERENCES attendances(id) ON DELETE CASCADE
);