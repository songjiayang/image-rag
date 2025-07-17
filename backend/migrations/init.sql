-- Image RAG Service Database Schema
-- MySQL initialization script

CREATE DATABASE IF NOT EXISTS image_rag CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE image_rag;

-- Records table for image metadata
CREATE TABLE IF NOT EXISTS records (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_name (name),
    INDEX idx_created_at (created_at)
);

-- Images table for storing image file references
CREATE TABLE IF NOT EXISTS images (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    record_id BIGINT NOT NULL,
    filename VARCHAR(255) NOT NULL,
    path VARCHAR(500) NOT NULL,
    vector_id VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_record_id (record_id),
    INDEX idx_vector_id (vector_id),
    INDEX idx_filename (filename),
    FOREIGN KEY (record_id) REFERENCES records(id) ON DELETE CASCADE
);

-- Sample data for testing
INSERT INTO records (name, description) VALUES
('Sample Cat', 'A cute domestic cat'),
('Sample Dog', 'A friendly golden retriever'),
('Sample Landscape', 'Beautiful mountain landscape');

INSERT INTO images (record_id, filename, path, vector_id) VALUES
(1, 'cat1.jpg', './uploads/cat1.jpg', 'vec_cat_001'),
(1, 'cat2.jpg', './uploads/cat2.jpg', 'vec_cat_002'),
(2, 'dog1.jpg', './uploads/dog1.jpg', 'vec_dog_001'),
(3, 'landscape1.jpg', './uploads/landscape1.jpg', 'vec_landscape_001');