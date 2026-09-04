package com.codejsha.bookstore.order.infrastructure.adapter.restcontroller

import com.codejsha.bookstore.order.domain.model.command.InvalidCommandException
import com.codejsha.bookstore.order.domain.workflow.OrderItemWorkflowRequest
import com.codejsha.bookstore.order.domain.workflow.PlaceOrderWorkflowRequest
import com.codejsha.bookstore.order.domain.workflow.ShippingWorkflowRequest
import org.junit.jupiter.api.Test
import java.math.BigDecimal
import kotlin.test.assertFailsWith

class ValidatePlaceOrderRequestTest {

    private fun shipping() = ShippingWorkflowRequest(
        recipientName = "Alice",
        recipientPhone = "010-1234-5678",
        addressLine1 = "1 Foo St",
        addressLine2 = null,
        city = "Seoul",
        state = "Seoul",
        postalCode = "12345",
        country = "KR",
        shippingMethod = "STANDARD",
    )

    private fun request(
        userUid: String = "22222222-2222-2222-2222-222222222222",
        currency: String = "USD",
        idempotencyKey: String = "idem-1",
        items: List<OrderItemWorkflowRequest> = listOf(
            OrderItemWorkflowRequest(
                productId = 10L,
                productName = "Book",
                quantity = 1,
                currency = "USD",
                price = BigDecimal("10.00"),
            ),
        ),
        shipping: ShippingWorkflowRequest? = shipping(),
    ) = PlaceOrderWorkflowRequest(
        userUid = userUid,
        currency = currency,
        items = items,
        shipping = shipping,
        idempotencyKey = idempotencyKey,
    )

    @Test
    fun `validatePlaceOrderRequest_whenRequestWellFormed_passes`() {
        validatePlaceOrderRequest(request())
        validatePlaceOrderRequest(request(shipping = null))
    }

    @Test
    fun `validatePlaceOrderRequest_whenUserUidMalformed_throwsBeforeStartingTheWorkflow`() {
        assertFailsWith<InvalidCommandException> { validatePlaceOrderRequest(request(userUid = "not-a-uuid")) }
    }

    @Test
    fun `validatePlaceOrderRequest_whenCurrencyMalformedOrIdempotencyKeyOverLimit_throws`() {
        assertFailsWith<InvalidCommandException> { validatePlaceOrderRequest(request(currency = "usd")) }
        assertFailsWith<InvalidCommandException> { validatePlaceOrderRequest(request(idempotencyKey = " ")) }
        assertFailsWith<InvalidCommandException> {
            validatePlaceOrderRequest(request(idempotencyKey = "k".repeat(101)))
        }
    }

    @Test
    fun `validatePlaceOrderRequest_whenShippingOverColumnWidth_throws`() {
        assertFailsWith<InvalidCommandException> {
            validatePlaceOrderRequest(request(shipping = shipping().copy(country = "KOR")))
        }
        assertFailsWith<InvalidCommandException> {
            validatePlaceOrderRequest(request(shipping = shipping().copy(state = "")))
        }
        assertFailsWith<InvalidCommandException> {
            validatePlaceOrderRequest(request(shipping = shipping().copy(city = "c".repeat(101))))
        }
        assertFailsWith<InvalidCommandException> {
            validatePlaceOrderRequest(request(shipping = shipping().copy(postalCode = "1".repeat(21))))
        }
    }

    @Test
    fun `validatePlaceOrderRequest_whenItemQuantityNotPositiveOrNameOverLimit_throws`() {
        val item = OrderItemWorkflowRequest(
            productId = 10L,
            productName = "Book",
            quantity = 1,
            currency = "USD",
            price = BigDecimal("10.00"),
        )
        assertFailsWith<InvalidCommandException> {
            validatePlaceOrderRequest(request(items = listOf(item.copy(quantity = 0))))
        }
        assertFailsWith<InvalidCommandException> {
            validatePlaceOrderRequest(request(items = listOf(item.copy(productName = "p".repeat(256)))))
        }
        assertFailsWith<InvalidCommandException> {
            validatePlaceOrderRequest(request(items = listOf(item.copy(price = BigDecimal("-1")))))
        }
    }
}
