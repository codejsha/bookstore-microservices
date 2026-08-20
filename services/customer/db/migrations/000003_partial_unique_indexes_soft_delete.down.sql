DROP INDEX IF EXISTS uq_review_user_book_uid;
ALTER TABLE customer_review
    ADD CONSTRAINT uq_review_user_book_uid UNIQUE (user_uid, book_uid);

DROP INDEX IF EXISTS uq_customer_wishlist_user_book_uid;
ALTER TABLE customer_wishlist
    ADD CONSTRAINT uq_customer_wishlist_user_book_uid UNIQUE (user_uid, book_uid);
