-- Campus Errand Platform - MySQL 8 schema
-- Money is stored as integer cents. Keep the same rule in application code.

CREATE DATABASE IF NOT EXISTS campus
    CHARACTER SET utf8mb4
    COLLATE utf8mb4_0900_ai_ci;

USE campus;
SET NAMES utf8mb4;

CREATE TABLE IF NOT EXISTS schools (
    school_id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    code VARCHAR(32) NOT NULL,
    name VARCHAR(100) NOT NULL,
    status TINYINT(1) NOT NULL DEFAULT 1,
    created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    PRIMARY KEY (school_id),
    UNIQUE KEY uk_schools_code (code),
    KEY idx_schools_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS users (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    username VARCHAR(50) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role TINYINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '0=publisher,1=runner,2=admin',
    school_id BIGINT UNSIGNED NULL,
    phone VARCHAR(20) NOT NULL DEFAULT '',
    status TINYINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '0=active,1=banned,2=disabled',
    ban_reason VARCHAR(500) NOT NULL DEFAULT '',
    violation_times BIGINT UNSIGNED NOT NULL DEFAULT 0,
    banned_at DATETIME(3) NULL,
    banned_until DATETIME(3) NULL,
    created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    PRIMARY KEY (id),
    UNIQUE KEY uk_users_username (username),
    KEY idx_users_school_role (school_id, role),
    KEY idx_users_phone (phone),
    KEY idx_users_status (status),
    KEY idx_users_banned_until (banned_until),
    CONSTRAINT fk_users_school FOREIGN KEY (school_id) REFERENCES schools (school_id)
        ON UPDATE CASCADE ON DELETE RESTRICT,
    CONSTRAINT chk_users_role CHECK (role BETWEEN 0 AND 2),
    CONSTRAINT chk_users_status CHECK (status BETWEEN 0 AND 2)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS runner_profiles (
    user_id BIGINT UNSIGNED NOT NULL,
    completed_order_count BIGINT UNSIGNED NOT NULL DEFAULT 0,
    positive_count BIGINT UNSIGNED NOT NULL DEFAULT 0,
    negative_count BIGINT UNSIGNED NOT NULL DEFAULT 0,
    created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    PRIMARY KEY (user_id),
    CONSTRAINT fk_runner_profiles_user FOREIGN KEY (user_id) REFERENCES users (id)
        ON UPDATE CASCADE ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS user_profiles (
    user_id BIGINT UNSIGNED NOT NULL,
    published_order_count BIGINT UNSIGNED NOT NULL DEFAULT 0,
    positive_count BIGINT UNSIGNED NOT NULL DEFAULT 0,
    negative_count BIGINT UNSIGNED NOT NULL DEFAULT 0,
    created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    PRIMARY KEY (user_id),
    CONSTRAINT fk_user_profiles_user FOREIGN KEY (user_id) REFERENCES users (id)
        ON UPDATE CASCADE ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS order_categories (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    school_id BIGINT UNSIGNED NOT NULL,
    code VARCHAR(32) NOT NULL,
    name VARCHAR(50) NOT NULL,
    description VARCHAR(255) NOT NULL DEFAULT '',
    icon_url VARCHAR(512) NOT NULL DEFAULT '',
    sort_order INT NOT NULL DEFAULT 0,
    is_active TINYINT(1) NOT NULL DEFAULT 1,
    created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    PRIMARY KEY (id),
    UNIQUE KEY uk_order_categories_school_code (school_id, code),
    KEY idx_order_categories_list (school_id, is_active, sort_order),
    CONSTRAINT fk_order_categories_school FOREIGN KEY (school_id) REFERENCES schools (school_id)
        ON UPDATE CASCADE ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS orders (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    order_no VARCHAR(32) NOT NULL,
    school_id BIGINT UNSIGNED NOT NULL,
    category_id BIGINT UNSIGNED NOT NULL,
    title VARCHAR(100) NOT NULL,
    description TEXT NOT NULL,
    publisher_id BIGINT UNSIGNED NOT NULL,
    runner_id BIGINT UNSIGNED NULL,
    remuneration_cents BIGINT NOT NULL,
    start_place VARCHAR(255) NOT NULL,
    end_place VARCHAR(255) NOT NULL,
    contact_name VARCHAR(50) NOT NULL,
    user_phone VARCHAR(20) NOT NULL,
    status TINYINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '0=published,1=accepted,2=delivering,3=arrived,4=confirmed,5=canceled',
    deadline DATETIME(3) NOT NULL,
    accept_time DATETIME(3) NULL,
    runner_finished_at DATETIME(3) NULL,
    confirmed_at DATETIME(3) NULL,
    auto_confirm_at DATETIME(3) NULL,
    canceled_by_id BIGINT UNSIGNED NULL,
    canceled_by_role TINYINT UNSIGNED NULL COMMENT '0=publisher,1=runner,2=admin',
    canceled_at DATETIME(3) NULL,
    cancel_reason VARCHAR(500) NOT NULL DEFAULT '',
    created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    PRIMARY KEY (id),
    UNIQUE KEY uk_orders_order_no (order_no),
    KEY idx_orders_school_status_created (school_id, status, created_at),
    KEY idx_orders_publisher_created (publisher_id, created_at),
    KEY idx_orders_runner_status (runner_id, status),
    KEY idx_orders_category (category_id),
    KEY idx_orders_deadline (deadline),
    KEY idx_orders_auto_confirm (auto_confirm_at),
    KEY idx_orders_canceled_by (canceled_by_id),
    CONSTRAINT fk_orders_school FOREIGN KEY (school_id) REFERENCES schools (school_id)
        ON UPDATE CASCADE ON DELETE RESTRICT,
    CONSTRAINT fk_orders_category FOREIGN KEY (category_id) REFERENCES order_categories (id)
        ON UPDATE CASCADE ON DELETE RESTRICT,
    CONSTRAINT fk_orders_publisher FOREIGN KEY (publisher_id) REFERENCES users (id)
        ON UPDATE CASCADE ON DELETE RESTRICT,
    CONSTRAINT fk_orders_runner FOREIGN KEY (runner_id) REFERENCES users (id)
        ON UPDATE CASCADE ON DELETE RESTRICT,
    CONSTRAINT fk_orders_canceled_by FOREIGN KEY (canceled_by_id) REFERENCES users (id)
        ON UPDATE CASCADE ON DELETE RESTRICT,
    CONSTRAINT chk_orders_remuneration CHECK (remuneration_cents >= 0),
    CONSTRAINT chk_orders_status CHECK (status BETWEEN 0 AND 5),
    CONSTRAINT chk_orders_canceled_role CHECK (canceled_by_role IS NULL OR canceled_by_role BETWEEN 0 AND 2)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS runner_ratings (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    order_id BIGINT UNSIGNED NOT NULL,
    content VARCHAR(500) NOT NULL DEFAULT '',
    user_id BIGINT UNSIGNED NOT NULL,
    runner_id BIGINT UNSIGNED NOT NULL,
    is_positive TINYINT(1) NOT NULL,
    created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    PRIMARY KEY (id),
    UNIQUE KEY uk_runner_ratings_order (order_id),
    KEY idx_runner_ratings_user (user_id),
    KEY idx_runner_ratings_runner (runner_id),
    CONSTRAINT fk_runner_ratings_order FOREIGN KEY (order_id) REFERENCES orders (id)
        ON UPDATE CASCADE ON DELETE RESTRICT,
    CONSTRAINT fk_runner_ratings_user FOREIGN KEY (user_id) REFERENCES users (id)
        ON UPDATE CASCADE ON DELETE RESTRICT,
    CONSTRAINT fk_runner_ratings_runner FOREIGN KEY (runner_id) REFERENCES users (id)
        ON UPDATE CASCADE ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS user_ratings (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    order_id BIGINT UNSIGNED NOT NULL,
    content VARCHAR(500) NOT NULL DEFAULT '',
    user_id BIGINT UNSIGNED NOT NULL,
    runner_id BIGINT UNSIGNED NOT NULL,
    is_positive TINYINT(1) NOT NULL,
    created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    PRIMARY KEY (id),
    UNIQUE KEY uk_user_ratings_order (order_id),
    KEY idx_user_ratings_user (user_id),
    KEY idx_user_ratings_runner (runner_id),
    CONSTRAINT fk_user_ratings_order FOREIGN KEY (order_id) REFERENCES orders (id)
        ON UPDATE CASCADE ON DELETE RESTRICT,
    CONSTRAINT fk_user_ratings_user FOREIGN KEY (user_id) REFERENCES users (id)
        ON UPDATE CASCADE ON DELETE RESTRICT,
    CONSTRAINT fk_user_ratings_runner FOREIGN KEY (runner_id) REFERENCES users (id)
        ON UPDATE CASCADE ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS order_status_logs (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    order_id BIGINT UNSIGNED NOT NULL,
    from_status TINYINT UNSIGNED NULL,
    to_status TINYINT UNSIGNED NOT NULL,
    operator_type VARCHAR(16) NOT NULL COMMENT 'publisher,runner,admin,system',
    operator_id BIGINT UNSIGNED NULL,
    reason VARCHAR(500) NOT NULL DEFAULT '',
    created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    PRIMARY KEY (id),
    KEY idx_order_status_logs_order_created (order_id, created_at),
    KEY idx_order_status_logs_operator (operator_id),
    CONSTRAINT fk_order_status_logs_order FOREIGN KEY (order_id) REFERENCES orders (id)
        ON UPDATE CASCADE ON DELETE RESTRICT,
    CONSTRAINT fk_order_status_logs_operator FOREIGN KEY (operator_id) REFERENCES users (id)
        ON UPDATE CASCADE ON DELETE RESTRICT,
    CONSTRAINT chk_order_status_logs_from CHECK (from_status IS NULL OR from_status BETWEEN 0 AND 5),
    CONSTRAINT chk_order_status_logs_to CHECK (to_status BETWEEN 0 AND 5),
    CONSTRAINT chk_order_status_logs_operator CHECK (operator_type IN ('publisher','runner','admin','system'))
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS order_payments (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    order_id BIGINT UNSIGNED NOT NULL,
    amount_cents BIGINT NOT NULL,
    method TINYINT UNSIGNED NULL COMMENT '0=wechat,1=alipay,2=cash,3=other',
    status TINYINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '0=pending,1=reported,2=confirmed,3=disputed,4=refunded,5=canceled',
    reported_by_id BIGINT UNSIGNED NULL,
    reported_at DATETIME(3) NULL,
    confirmed_by_id BIGINT UNSIGNED NULL,
    confirmed_at DATETIME(3) NULL,
    note VARCHAR(500) NOT NULL DEFAULT '',
    created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    PRIMARY KEY (id),
    UNIQUE KEY uk_order_payments_order (order_id),
    KEY idx_order_payments_status (status),
    KEY idx_order_payments_reported_by (reported_by_id),
    KEY idx_order_payments_confirmed_by (confirmed_by_id),
    CONSTRAINT fk_order_payments_order FOREIGN KEY (order_id) REFERENCES orders (id)
        ON UPDATE CASCADE ON DELETE RESTRICT,
    CONSTRAINT fk_order_payments_reported_by FOREIGN KEY (reported_by_id) REFERENCES users (id)
        ON UPDATE CASCADE ON DELETE RESTRICT,
    CONSTRAINT fk_order_payments_confirmed_by FOREIGN KEY (confirmed_by_id) REFERENCES users (id)
        ON UPDATE CASCADE ON DELETE RESTRICT,
    CONSTRAINT chk_order_payments_amount CHECK (amount_cents >= 0),
    CONSTRAINT chk_order_payments_method CHECK (method IS NULL OR method BETWEEN 0 AND 3),
    CONSTRAINT chk_order_payments_status CHECK (status BETWEEN 0 AND 5)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS order_payment_proofs (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    payment_id BIGINT UNSIGNED NOT NULL,
    uploaded_by_id BIGINT UNSIGNED NULL,
    file_url VARCHAR(512) NOT NULL,
    created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    PRIMARY KEY (id),
    KEY idx_order_payment_proofs_payment (payment_id),
    KEY idx_order_payment_proofs_uploader (uploaded_by_id),
    CONSTRAINT fk_order_payment_proofs_payment FOREIGN KEY (payment_id) REFERENCES order_payments (id)
        ON UPDATE CASCADE ON DELETE CASCADE,
    CONSTRAINT fk_order_payment_proofs_uploader FOREIGN KEY (uploaded_by_id) REFERENCES users (id)
        ON UPDATE CASCADE ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS user_violations (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    order_id BIGINT UNSIGNED NULL,
    reporter_id BIGINT UNSIGNED NULL,
    target_user_id BIGINT UNSIGNED NOT NULL,
    report_type VARCHAR(50) NOT NULL,
    description TEXT NOT NULL,
    status TINYINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '0=reviewing,1=confirmed,2=rejected,3=closed',
    violation_type VARCHAR(50) NULL,
    penalty_type TINYINT UNSIGNED NULL COMMENT '0=warning,1=suspend,2=ban',
    penalty_start_at DATETIME(3) NULL,
    penalty_end_at DATETIME(3) NULL,
    decision VARCHAR(1000) NULL,
    handled_by BIGINT UNSIGNED NULL,
    resolved_at DATETIME(3) NULL,
    created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    PRIMARY KEY (id),
    KEY idx_user_violations_order (order_id),
    KEY idx_user_violations_reporter (reporter_id),
    KEY idx_user_violations_target_status (target_user_id, status),
    KEY idx_user_violations_handled_by (handled_by),
    CONSTRAINT fk_user_violations_order FOREIGN KEY (order_id) REFERENCES orders (id)
        ON UPDATE CASCADE ON DELETE SET NULL,
    CONSTRAINT fk_user_violations_reporter FOREIGN KEY (reporter_id) REFERENCES users (id)
        ON UPDATE CASCADE ON DELETE SET NULL,
    CONSTRAINT fk_user_violations_target FOREIGN KEY (target_user_id) REFERENCES users (id)
        ON UPDATE CASCADE ON DELETE RESTRICT,
    CONSTRAINT fk_user_violations_handled_by FOREIGN KEY (handled_by) REFERENCES users (id)
        ON UPDATE CASCADE ON DELETE SET NULL,
    CONSTRAINT chk_user_violations_status CHECK (status BETWEEN 0 AND 3),
    CONSTRAINT chk_user_violations_penalty CHECK (penalty_type IS NULL OR penalty_type BETWEEN 0 AND 2)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;