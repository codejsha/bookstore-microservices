package com.codejsha.bookstore.order.domain.aggregate

import com.codejsha.bookstore.order.domain.constant.OrderAdjustmentType
import com.codejsha.bookstore.order.domain.constant.OrderStatus
import com.codejsha.bookstore.order.support.OrderTestFixtures
import org.junit.jupiter.api.Test
import kotlin.test.assertEquals
import kotlin.test.assertFailsWith

class OrderMappingTest {

    @Test
    fun `OrderResult toAggregate converts status string and zeroes nested collections`() {
        val result = OrderTestFixtures.orderResult(status = OrderStatus.PAID.value)

        val agg = result.toAggregate()

        assertEquals(OrderStatus.PAID, agg.status)
        assertEquals(0, agg.items.size)
        assertEquals(0, agg.adjustments.size)
        assertEquals(null, agg.shipping)
        assertEquals(result.totalAmount, agg.totalAmount)
        assertEquals(result.idempotencyKey, agg.idempotencyKey)
    }

    @Test
    fun `OrderResult toAggregate throws when status is unknown`() {
        val result = OrderTestFixtures.orderResult(status = "MYSTERY")

        assertFailsWith<NoSuchElementException> { result.toAggregate() }
    }

    @Test
    fun `OrderItemResult toEntity copies all fields`() {
        val result = OrderTestFixtures.orderItemResult(quantity = 3)

        val entity = result.toEntity()

        assertEquals(result.uid, entity.uid)
        assertEquals(3, entity.quantity)
        assertEquals(result.price, entity.price)
        assertEquals(result.subtotal, entity.subtotal)
    }

    @Test
    fun `OrderAdjustmentResult toEntity converts type string to enum`() {
        val result = OrderTestFixtures.orderAdjustmentResult(type = "COUPON")
        val entity = result.toEntity()

        assertEquals(OrderAdjustmentType.COUPON, entity.type)
        assertEquals(result.amount, entity.amount)
    }

    @Test
    fun `OrderAdjustmentResult toEntity throws when type unknown`() {
        val result = OrderTestFixtures.orderAdjustmentResult(type = "MYSTERY")

        assertFailsWith<NoSuchElementException> { result.toEntity() }
    }

    @Test
    fun `OrderShippingResult toEntity copies all fields`() {
        val result = OrderTestFixtures.orderShippingResult()
        val entity = result.toEntity()

        assertEquals(result.recipientName, entity.recipientName)
        assertEquals(result.country, entity.country)
        assertEquals(result.shippingMethod, entity.shippingMethod)
    }

    @Test
    fun `CartResult toAggregate yields empty items list`() {
        val cart = OrderTestFixtures.cartResult().toAggregate()

        assertEquals(0, cart.items.size)
    }

    @Test
    fun `CartItemResult toEntity preserves quantity, price, currency`() {
        val item = OrderTestFixtures.cartItemResult(quantity = 5).toEntity()

        assertEquals(5, item.quantity)
        assertEquals("KRW", item.currency)
    }
}

class OrderEnumsTest {

    @Test
    fun `OrderStatus round-trip and unknown lookup`() {
        OrderStatus.entries.forEach { assertEquals(it, OrderStatus.fromValue(it.value)) }
        assertFailsWith<NoSuchElementException> { OrderStatus.fromValue("UNKNOWN") }
    }

    @Test
    fun `OrderAdjustmentType round-trip`() {
        OrderAdjustmentType.entries.forEach {
            assertEquals(it, OrderAdjustmentType.fromValue(it.value))
        }
    }
}
