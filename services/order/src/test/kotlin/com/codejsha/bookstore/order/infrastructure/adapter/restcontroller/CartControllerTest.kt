package com.codejsha.bookstore.order.infrastructure.adapter.restcontroller

import com.codejsha.bookstore.generated.application.port.openapi.model.CartAddItemRequest
import com.codejsha.bookstore.generated.application.port.openapi.model.CartCheckoutRequest
import com.codejsha.bookstore.generated.application.port.openapi.model.CartUpdateItemRequest
import com.codejsha.bookstore.order.application.usecase.CartUseCase
import com.codejsha.bookstore.order.domain.aggregate.CartAggregate
import com.codejsha.bookstore.order.domain.aggregate.CartItemEntity
import com.codejsha.bookstore.order.domain.aggregate.OrderAggregate
import com.codejsha.bookstore.order.domain.constant.OrderStatus
import com.codejsha.bookstore.order.domain.model.command.CartAddItemCommand
import com.codejsha.bookstore.order.domain.model.command.CartCheckoutCommand
import com.codejsha.bookstore.order.infrastructure.support.auth.ForbiddenException
import com.codejsha.bookstore.order.infrastructure.support.auth.HttpPrincipalResolver
import com.codejsha.bookstore.order.support.OrderTestFixtures
import com.codejsha.platform.shared.data.ActorContext
import com.codejsha.platform.shared.data.ActorType
import kotlinx.coroutines.runBlocking
import org.junit.jupiter.api.AfterEach
import org.junit.jupiter.api.BeforeEach
import org.junit.jupiter.api.Test
import org.mockito.BDDMockito.given
import org.mockito.Mockito.mock
import org.mockito.Mockito.verify
import org.springframework.http.HttpStatus
import org.springframework.mock.web.MockHttpServletRequest
import org.springframework.web.context.request.RequestContextHolder
import org.springframework.web.context.request.ServletRequestAttributes
import tools.jackson.databind.ObjectMapper
import java.math.BigDecimal
import java.time.LocalDateTime
import java.util.UUID
import kotlin.test.assertEquals
import kotlin.test.assertFailsWith
import kotlin.test.assertNotNull

class CartControllerTest {

    private val userUid = OrderTestFixtures.USER_UID
    private val itemUid: UUID = UUID.fromString("99999999-9999-9999-9999-999999999999")
    private val controllerContext = ActorContext(actorId = 0L, actorType = ActorType.USER)
    private val resolver = HttpPrincipalResolver(ObjectMapper())

    @BeforeEach
    fun bindPrincipal() = bindSubject(userUid.toString())

    @AfterEach
    fun clearPrincipal() = RequestContextHolder.resetRequestAttributes()

    private fun bindSubject(sub: String) {
        val request = MockHttpServletRequest()
        request.addHeader("X-User-Id", sub)
        RequestContextHolder.setRequestAttributes(ServletRequestAttributes(request))
    }

    private fun cartItem(quantity: Int = 2, price: String = "10.50") = CartItemEntity(
        id = 1L,
        uid = itemUid,
        cartId = 1L,
        productId = 42L,
        productName = "Dune",
        quantity = quantity,
        currency = "USD",
        price = BigDecimal(price),
        createdAt = LocalDateTime.now(),
        updatedAt = null,
    )

    private fun cart(items: List<CartItemEntity> = listOf(cartItem())) = CartAggregate(
        id = 1L,
        uid = OrderTestFixtures.CART_UID,
        userUid = userUid,
        items = items,
        createdAt = LocalDateTime.now(),
        updatedAt = null,
    )

    @Test
    fun `cartGet returns the caller's own cart with derived totals`(): Unit = runBlocking {
        val useCase = mock(CartUseCase::class.java)
        val controller = CartController(useCase, resolver)
        given(useCase.getCart(userUid, controllerContext))
            .willReturn(cart(listOf(cartItem(quantity = 2, price = "10.50"), cartItem(quantity = 1, price = "5.00"))))

        val response = controller.cartGet()

        assertEquals(HttpStatus.OK, response.statusCode)
        val body = assertNotNull(response.body)
        assertEquals(userUid.toString(), body.userUid)
        assertEquals(3, body.totalItems)
        assertEquals(26.0, body.totalAmount)
        assertEquals(21.0, body.items[0].subtotal)
    }

    @Test
    fun `cartAddItem forwards the item to the caller's cart`(): Unit = runBlocking {
        val useCase = mock(CartUseCase::class.java)
        val controller = CartController(useCase, resolver)
        val expected = CartAddItemCommand(
            productId = 42L,
            productName = "Dune",
            quantity = 2,
            currency = "USD",
            price = BigDecimal.valueOf(10.5),
        )
        given(useCase.addItem(userUid, expected, controllerContext)).willReturn(cartItem())

        val response = controller.cartAddItem(
            CartAddItemRequest(productId = 42L, productName = "Dune", quantity = 2, currency = "USD", price = 10.5),
        )

        assertEquals(HttpStatus.NO_CONTENT, response.statusCode)
        verify(useCase).addItem(userUid, expected, controllerContext)
    }

    @Test
    fun `cartUpdateItem scopes the item to the caller`(): Unit = runBlocking {
        val useCase = mock(CartUseCase::class.java)
        val controller = CartController(useCase, resolver)
        given(useCase.updateItemQuantity(userUid, itemUid, 5, controllerContext))
            .willReturn(cartItem(quantity = 5))

        val response = controller.cartUpdateItem(itemUid.toString(), CartUpdateItemRequest(quantity = 5))

        assertEquals(HttpStatus.OK, response.statusCode)
        assertEquals(5, assertNotNull(response.body).quantity)
        verify(useCase).updateItemQuantity(userUid, itemUid, 5, controllerContext)
    }

    @Test
    fun `cartRemoveItem scopes the item to the caller`(): Unit = runBlocking {
        val useCase = mock(CartUseCase::class.java)
        val controller = CartController(useCase, resolver)

        val response = controller.cartRemoveItem(itemUid.toString())

        assertEquals(HttpStatus.NO_CONTENT, response.statusCode)
        verify(useCase).removeItem(userUid, itemUid, controllerContext)
    }

    @Test
    fun `cartClear clears the caller's cart`(): Unit = runBlocking {
        val useCase = mock(CartUseCase::class.java)
        val controller = CartController(useCase, resolver)

        val response = controller.cartClear()

        assertEquals(HttpStatus.NO_CONTENT, response.statusCode)
        verify(useCase).clearCart(userUid, controllerContext)
    }

    private fun placedOrder() = OrderAggregate(
        id = 1L,
        uid = OrderTestFixtures.ORDER_UID,
        userUid = userUid,
        orderNumber = "ORD-20260101-0001",
        status = OrderStatus.PENDING,
        currency = "USD",
        itemsAmount = BigDecimal("26.00"),
        discountAmount = BigDecimal.ZERO,
        shippingAmount = BigDecimal.ZERO,
        taxAmount = BigDecimal.ZERO,
        totalAmount = BigDecimal("26.00"),
        idempotencyKey = "key-1",
        paymentUid = null,
        items = emptyList(),
        adjustments = emptyList(),
        shipping = null,
        createdAt = LocalDateTime.of(2026, 1, 1, 0, 0),
        updatedAt = null,
    )

    @Test
    fun `cartCheckout returns the placed order`(): Unit = runBlocking {
        val useCase = mock(CartUseCase::class.java)
        val controller = CartController(useCase, resolver)
        val expected = CartCheckoutCommand(currency = "USD", idempotencyKey = "key-1", shipping = null)
        given(useCase.checkout(userUid, expected, controllerContext))
            .willReturn(placedOrder())

        val response = controller.cartCheckout(
            CartCheckoutRequest(currency = "USD", idempotencyKey = "key-1", shipping = null),
        )

        assertEquals(HttpStatus.OK, response.statusCode)
        assertEquals(userUid.toString(), assertNotNull(response.body).userUid)
        verify(useCase).checkout(userUid, expected, controllerContext)
    }

    @Test
    fun `a non-uuid subject owns no cart`(): Unit = runBlocking {
        val useCase = mock(CartUseCase::class.java)
        val controller = CartController(useCase, resolver)
        bindSubject("service-account-batch")

        assertFailsWith<ForbiddenException> { controller.cartGet() }
    }
}
