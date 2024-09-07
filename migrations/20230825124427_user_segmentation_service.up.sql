CREATE TABLE users (
    user_id VARCHAR(40) PRIMARY KEY,
    created_at TIMESTAMP not null default now()
);

CREATE TABLE segments (
    slug VARCHAR PRIMARY KEY,
    created_at TIMESTAMP not null default now()
);

CREATE TABLE user_segments (
    user_id VARCHAR(40) REFERENCES users(user_id),
    segment_slug VARCHAR REFERENCES segments(slug) ON DELETE CASCADE,
    PRIMARY KEY (user_id, segment_slug)
);

CREATE TABLE api_keys (
    id SERIAL PRIMARY KEY,
    hash_key VARCHAR
);

CREATE TABLE history (
    history_id SERIAL PRIMARY KEY,
    user_id VARCHAR(40) REFERENCES users(user_id),
    segment_slug VARCHAR,
    type VARCHAR(15),
    created_at TIMESTAMP not null default now()
);

CREATE TABLE tasks_delete (
    task_id SERIAL PRIMARY KEY,
    user_id VARCHAR(40) REFERENCES users(user_id),
    segment_slug VARCHAR,
    deadline TIMESTAMP not null,
    done BOOLEAN default false
);