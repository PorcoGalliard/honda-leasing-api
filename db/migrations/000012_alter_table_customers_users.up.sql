-- Step 1: Tambah kolom user_id tanpa NOT NULL dulu
ALTER TABLE dealer.customers
ADD COLUMN user_id BIGINT UNIQUE;

-- Step 2: Update data existing customers dengan user_id yang sesuai (one-to-one mapping)
-- Map Kang Dian ke user_id 7 (Kang Dian - Customer)
UPDATE dealer.customers 
SET user_id = (SELECT user_id FROM account.users WHERE full_name LIKE '%Kang Dian%' LIMIT 1)
WHERE nama_lengkap = 'Kang Dian';

-- Map Winona ke user_id 8 (Winona - Customer)
UPDATE dealer.customers 
SET user_id = (SELECT user_id FROM account.users WHERE full_name LIKE '%Winona%' LIMIT 1)
WHERE nama_lengkap = 'Winona';

-- Step 3: Set constraint NOT NULL setelah semua data terisi
ALTER TABLE dealer.customers
ALTER COLUMN user_id SET NOT NULL;

-- Step 4: Tambah foreign key constraint
ALTER TABLE dealer.customers
ADD CONSTRAINT fk_customer_user
FOREIGN KEY (user_id)
REFERENCES account.users(user_id)
ON DELETE CASCADE;