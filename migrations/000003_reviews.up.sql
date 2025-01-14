CREATE TABLE reviews (
  id UUID PRIMARY KEY,
  business_id UUID REFERENCES businesses(id) ON DELETE CASCADE,
  user_id UUID REFERENCES users(id) ON DELETE CASCADE,
  rating INT NOT NULL CHECK (rating > 0 AND rating <= 5),  
  comment TEXT,
  created_at TIMESTAMP NOT NULL DEFAULT 'now()',
  updated_at TIMESTAMP NOT NULL DEFAULT 'now()'
);

CREATE TABLE reviews_attachments (
  id UUID PRIMARY KEY,
  review_id UUID REFERENCES reviews(id) ON DELETE CASCADE,
  filepath VARCHAR(255) NOT NULL,
  content_type attachment_type NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT 'now()',
  updated_at TIMESTAMP NOT NULL DEFAULT 'now()'
);