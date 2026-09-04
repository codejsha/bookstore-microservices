ALTER TABLE payment_methods
    DROP FOREIGN KEY fk_pm_customer;

ALTER TABLE mandates
    DROP FOREIGN KEY fk_mandate_customer;

ALTER TABLE mandates
    DROP FOREIGN KEY fk_mandate_pm;

ALTER TABLE mandates
    RENAME INDEX fk_mandate_pm TO idx_mandate_pm;

ALTER TABLE payments
    DROP FOREIGN KEY fk_payment_customer;

ALTER TABLE payments
    DROP FOREIGN KEY fk_payment_pm;

ALTER TABLE payments
    RENAME INDEX fk_payment_pm TO idx_payment_pm;

ALTER TABLE payments
    DROP FOREIGN KEY fk_payment_mandate;

ALTER TABLE payment_attempts
    DROP FOREIGN KEY fk_attempt_payment;

ALTER TABLE refunds
    DROP FOREIGN KEY fk_refund_payment;
