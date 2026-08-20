ALTER TABLE customer_wishlist DROP CONSTRAINT uq_customer_wishlist_user_book;
ALTER TABLE customer_wishlist
    ADD CONSTRAINT uq_customer_wishlist_user_book_uid UNIQUE (user_uid, book_uid);
CREATE INDEX idx_wishlist_user_uid ON customer_wishlist (user_uid);

ALTER TABLE customer_point DROP CONSTRAINT uq_customer_point_user;
ALTER TABLE customer_point
    ADD CONSTRAINT uq_customer_point_user_uid UNIQUE (user_uid);
CREATE INDEX idx_point_user_uid ON customer_point (user_uid);

ALTER TABLE customer_review DROP CONSTRAINT uq_review_user_book;
ALTER TABLE customer_review
    ADD CONSTRAINT uq_review_user_book_uid UNIQUE (user_uid, book_uid);
CREATE INDEX idx_review_user_uid ON customer_review (user_uid);
CREATE INDEX idx_review_book_uid ON customer_review (book_uid);
