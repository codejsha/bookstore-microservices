package com.codejsha.bookstore.order.infrastructure.adapter.restcontroller

import com.codejsha.bookstore.order.domain.aggregate.CartItemEntity
import com.codejsha.bookstore.order.domain.workflow.OrderItemWorkflowRequest
import com.codejsha.bookstore.order.domain.workflow.PlaceOrderWorkflowRequest
import com.codejsha.bookstore.order.infrastructure.support.auth.BadRequestException
import org.junit.jupiter.api.Test
import java.math.BigDecimal
import java.time.LocalDateTime
import java.util.UUID
import kotlin.test.assertEquals
import kotlin.test.assertFailsWith

class RepriceFromCartTest {

    private fun cartItem(productId: Long, price: String, currency: String = "USD") = CartItemEntity(
        id = productId,
        uid = UUID.randomUUID(),
        cartId = 1L,
        productId = productId,
        productName = "product-$productId",
        quantity = 99,
        currency = currency,
        price = BigDecimal(price),
        createdAt = LocalDateTime.of(2026, 1, 1, 0, 0),
        updatedAt = null,
    )

    private fun requestItem(productId: Long, price: String, currency: String = "USD", quantity: Int = 1) =
        OrderItemWorkflowRequest(
            productId = productId,
            productName = "client-name-$productId",
            quantity = quantity,
            currency = currency,
            price = BigDecimal(price),
        )

    private fun request(items: List<OrderItemWorkflowRequest>) = PlaceOrderWorkflowRequest(
        userUid = "22222222-2222-2222-2222-222222222222",
        currency = "USD",
        items = items,
        shipping = null,
        idempotencyKey = "idem-1",
    )

    @Test
    fun `overrides client prices and currency with the cart's authoritative values`() {
        val cart = listOf(cartItem(productId = 10L, price = "10.00"), cartItem(productId = 11L, price = "5.50"))
        val req = request(
            listOf(
                requestItem(productId = 10L, price = "0.01", currency = "XXX", quantity = 3),
                requestItem(productId = 11L, price = "0.02", currency = "XXX", quantity = 2),
            ),
        )

        val result = repriceFromCart(req, cart)

        assertEquals(BigDecimal("10.00"), result.items[0].price)
        assertEquals("USD", result.items[0].currency)
        assertEquals(3, result.items[0].quantity)
        assertEquals(BigDecimal("5.50"), result.items[1].price)
        assertEquals("USD", result.items[1].currency)
        assertEquals(2, result.items[1].quantity)
    }

    @Test
    fun `rejects an item whose product is not in the cart`() {
        val cart = listOf(cartItem(productId = 10L, price = "10.00"))
        val req = request(
            listOf(
                requestItem(productId = 10L, price = "10.00"),
                requestItem(productId = 999L, price = "10.00"),
            ),
        )

        assertFailsWith<BadRequestException> { repriceFromCart(req, cart) }
    }

    @Test
    fun `rejects an empty item list so the saga never charges a zero total`() {
        assertFailsWith<BadRequestException> { repriceFromCart(request(emptyList()), emptyList()) }
    }
}
