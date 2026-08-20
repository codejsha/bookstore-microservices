ALTER TABLE customer_wishlist DROP CONSTRAINT uq_customer_wishlist_user_book_uid;
CREATE UNIQUE INDEX uq_customer_wishlist_user_book_uid
    ON customer_wishlist (user_uid, book_uid)
    WHERE deleted_at IS NULL;

ALTER TABLE customer_review DROP CONSTRAINT uq_review_user_book_uid;
CREATE UNIQUE INDEX uq_review_user_book_uid
    ON customer_review (user_uid, book_uid)
    WHERE deleted_at IS NULL;
