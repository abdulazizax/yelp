-- Drop Tables (Migration Down)
DROP TABLE IF EXISTS businesses_attachment;
DROP TABLE IF EXISTS businesses;
DROP TABLE IF EXISTS business_categories;


-- Drop Enums (Migration Down)
DROP TYPE IF EXISTS attachment_type;