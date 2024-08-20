ALTER TABLE users
   DROP COLUMN IF EXISTS new_column,
   ALTER COLUMN existing_column TYPE VARCHAR(255);
