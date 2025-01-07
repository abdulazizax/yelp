-- Drop Tables (Migration Down)
DROP TABLE IF EXISTS sessions;
DROP TABLE IF EXISTS users;

-- Drop Enums (Migration Down)
DROP TYPE IF EXISTS platform;
DROP TYPE IF EXISTS user_status;
DROP TYPE IF EXISTS gender;
DROP TYPE IF EXISTS user_role;
DROP TYPE IF EXISTS user_type;
