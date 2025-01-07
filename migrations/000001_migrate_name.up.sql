-- Enums
CREATE TYPE user_type AS ENUM ('user', 'admin');
CREATE TYPE user_role AS ENUM ('user', 'admin', 'business_owner', 'super_admin');
CREATE TYPE gender AS ENUM ('male', 'female');
CREATE TYPE user_status AS ENUM ('active', 'blocked', 'inverify');
CREATE TYPE platform AS ENUM ('web', 'mobile', 'admin_web');

-- Tables
CREATE TABLE users (
    id UUID PRIMARY KEY,
    user_type user_type NOT NULL,
    user_role user_role NOT NULL,
    full_name VARCHAR(50) NOT NULL,
    username varchar(50) UNIQUE NOT NULL,
    email VARCHAR(50) UNIQUE NOT NULL,
    password VARCHAR(256) NOT NULL,
    bio TEXT,
    gender gender NOT NULL DEFAULT 'male',
    profile_picture VARCHAR(255),
    status user_status NOT NULL DEFAULT 'inverify',
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE session (
  id uuid PRIMARY KEY,
  user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  user_agent text NOT NULL,
  platform platform NOT NULL,
  ip_address varchar(64) NOT NULL,
  is_active boolean NOT NULL DEFAULT true,
  expires_at timestamp NOT NULL,
  last_active_at timestamp NOT NULL DEFAULT now(),
  created_at timestamp NOT NULL DEFAULT now(),
  updated_at timestamp NOT NULL DEFAULT now()
);
