CREATE TABLE metrics (
    id SERIAL PRIMARY KEY,
    cluster_id INT REFERENCES clusters(id) ON DELETE CASCADE,
    cpu_usage FLOAT,
    memory_usage FLOAT,
    disk_usage FLOAT,
    created_at TIMESTAMP DEFAULT NOW()
);

