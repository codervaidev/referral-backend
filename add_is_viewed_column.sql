-- Migration to add is_viewed column to gems_history table
-- Run this SQL command on your PostgreSQL database

ALTER TABLE gems_history 
ADD COLUMN is_viewed BOOLEAN DEFAULT FALSE;

-- Optional: Update existing records to set is_viewed to false explicitly
UPDATE gems_history 
SET is_viewed = FALSE 
WHERE is_viewed IS NULL; 