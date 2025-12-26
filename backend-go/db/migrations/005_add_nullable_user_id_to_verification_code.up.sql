-- Make user_id nullable for delayed registration
ALTER TABLE public.verification_code ALTER COLUMN user_id DROP NOT NULL;
