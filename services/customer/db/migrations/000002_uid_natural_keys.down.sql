DROP INDEX IF EXISTS idx_review_book_uid;
DROP INDEX IF EXISTS idx_review_user_uid;
ALTER TABLE customer_review DROP CONSTRAINT uq_review_user_book_uid;
ALTER TABLE customer_review
    ADD CONSTRAINT uq_review_user_book UNIQUE (user_id, book_id);

DROP INDEX IF EXISTS idx_point_user_uid;
ALTER TABLE customer_point DROP CONSTRAINT uq_customer_point_user_uid;
ALTER TABLE customer_point
    ADD CONSTRAINT uq_customer_point_user UNIQUE (user_id);

DROP INDEX IF EXISTS idx_wishlist_user_uid;
ALTER TABLE customer_wishlist DROP CONSTRAINT uq_customer_wishlist_user_book_uid;
ALTER TABLE customer_wishlist
    ADD CONSTRAINT uq_customer_wishlist_user_book UNIQUE (user_id, book_id);
