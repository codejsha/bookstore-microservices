ALTER TABLE customer_review
    ADD CONSTRAINT chk_review_rating CHECK (rating BETWEEN 1 AND 5);

ALTER TABLE customer_point_history
    ADD CONSTRAINT chk_point_history_type
        CHECK (change_type IN ('EARN', 'SPEND', 'ADJUST', 'EXPIRE'));
