CREATE TABLE servies (
	id SERIAL PRIMARY KEY,
	cluster_id INT REFERENCES clusters(id) ON DELETE CASCADE,
	name VARCHAR(100) NOT NULL,
	image VARCHAR(255),
	replicas INT DEFAULT 1,
	status VARCHAR(50) DEFAULT 'pending',
	created_at TIMESTAMP DEFAULT NOW(),
	updated_at TIMESTAMP DEFAULT NOW()
);