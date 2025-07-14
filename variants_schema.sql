-- Corrected variants table schema
CREATE TABLE public.variants (
    id serial4 NOT NULL,
    product_id int4 NULL,
    variant_name varchar(255) NULL,
    pics jsonb NULL, -- This is already correct for JSON array
    variant_type text NULL,
    CONSTRAINT variants_pkey PRIMARY KEY (id),
    CONSTRAINT variants_product_id_fkey FOREIGN KEY (product_id) REFERENCES public.products(id)
);

-- Add a check constraint to ensure pics is always an array
ALTER TABLE public.variants 
ADD CONSTRAINT variants_pics_array_check 
CHECK (pics IS NULL OR jsonb_typeof(pics) = 'array');

-- Example of how to insert data with JSON array
-- INSERT INTO variants (product_id, variant_name, pics, variant_type) 
-- VALUES (1, 'Red Color', '["image1.jpg", "image2.jpg", "image3.jpg"]'::jsonb, 'color');

-- Example of how to query the first image from the array
-- SELECT id, variant_name, pics->0 as first_image FROM variants; 