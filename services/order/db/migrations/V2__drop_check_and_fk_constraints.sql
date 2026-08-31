ALTER TABLE orders
    DROP CHECK chk_orders_status;
ALTER TABLE orders
    DROP CHECK chk_orders_currency;
ALTER TABLE orders
    DROP CHECK chk_orders_amounts;

ALTER TABLE order_item
    DROP FOREIGN KEY fk_order_item__order;
ALTER TABLE order_item
    DROP CHECK chk_order_item_currency;
ALTER TABLE order_item
    DROP CHECK chk_order_item_price;
ALTER TABLE order_item
    DROP CHECK chk_order_item_subtotal;

ALTER TABLE order_adjustment
    DROP FOREIGN KEY fk_order_adjustment__order;
ALTER TABLE order_adjustment
    DROP CHECK chk_order_adjustment_type;

ALTER TABLE order_shipping
    DROP FOREIGN KEY fk_order_shipping__order;

ALTER TABLE cart_item
    DROP FOREIGN KEY fk_cart_item__cart;
ALTER TABLE cart_item
    DROP CHECK chk_cart_item_quantity;
ALTER TABLE cart_item
    DROP CHECK chk_cart_item_price;
