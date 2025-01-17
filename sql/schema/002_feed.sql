-- +goose Up
CREATE TABLE feeds (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),             
    name VARCHAR(255) NOT NULL,        
    url TEXT UNIQUE NOT NULL,                 
    user_id UUID NOT NULL,              
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_user FOREIGN KEY(user_id) 
      REFERENCES users(id) 
      ON DELETE CASCADE
);

-- +goose Down
DROP TABLE feeds;
