package com.codejsha.bookstore.order.domain.service

import com.codejsha.bookstore.order.application.port.repo.CartItemRepo
import com.codejsha.bookstore.order.application.port.repo.CartRepo
import com.codejsha.bookstore.order.application.port.repo.OrderItemRepo
import com.codejsha.bookstore.order.application.port.repo.OrderRepo
import com.codejsha.bookstore.order.application.port.repo.OrderShippingRepo
import com.codejsha.bookstore.order.domain.model.CartStateConflictException
import com.codejsha.bookstore.order.domain.model.command.CartAddItemCommand
import com.codejsha.bookstore.order.domain.model.command.CartCheckoutCommand
import com.codejsha.bookstore.order.domain.model.command.OrderCreateCommand
import com.codejsha.bookstore.order.domain.model.command.OrderItemCreateCommand
import com.codejsha.bookstore.order.support.FakeDistributedLock
import com.codejsha.bookstore.order.support.FakeIdempotencyService
import com.codejsha.bookstore.order.support.FakeTransactionRunner
import com.codejsha.bookstore.order.support.OrderTestFixtures
import kotlinx.coroutines.runBlocking
import org.junit.jupiter.api.Test
import org.mockito.BDDMockito.given
import org.mockito.Mockito.inOrder
import org.mockito.Mockito.mock
import org.mockito.Mockito.never
import org.mockito.Mockito.verify
import org.mockito.Mockito.verifyNoInteractions
import java.math.BigDecimal
import kotlin.test.assertEquals
import kotlin.test.assertFailsWith
import kotlin.test.assertNotNull
import kotlin.test.assertNull

class CartServiceTest {

    private val ctx = OrderTestFixtures.DEFAULT_CONTEXT
    private val userUid = OrderTestFixtures.USER_UID

    @Test
    fun `getCart_whenCartMissing_createsCart`() {
        val cartRepo = mock(CartRepo::class.java)
        val cartItemRepo = mock(CartItemRepo::class.java)
        val service = newService(cartRepo, cartItemRepo)

        given(cartRepo.findByUser(userUid, ctx)).willReturn(null)
        given(cartRepo.createForUser(userUid, ctx))
            .willReturn(OrderTestFixtures.cartResult(id = 1L, userUid = userUid))
        given(cartItemRepo.findAllByCart(1L, ctx)).willReturn(emptyList())

        val agg = runBlocking { service.getCart(userUid, ctx) }

        assertEquals(userUid, agg.userUid)
        assertEquals(0, agg.items.size)
        verify(cartRepo).createForUser(userUid, ctx)
    }

    @Test
    fun `getCart_whenCartExists_returnsCartWithItems`() {
        val cartRepo = mock(CartRepo::class.java)
        val cartItemRepo = mock(CartItemRepo::class.java)
        val service = newService(cartRepo, cartItemRepo)

        given(cartRepo.findByUser(userUid, ctx))
            .willReturn(OrderTestFixtures.cartResult(id = 1L))
        given(cartItemRepo.findAllByCart(1L, ctx))
            .willReturn(listOf(OrderTestFixtures.cartItemResult(quantity = 3)))

        val agg = runBlocking { service.getCart(userUid, ctx) }

        assertEquals(1, agg.items.size)
        assertEquals(3, agg.items[0].quantity)
    }

    @Test
    fun `addItem_whenProductAlreadyInCart_mergesQuantity`() {
        val cartRepo = mock(CartRepo::class.java)
        val cartItemRepo = mock(CartItemRepo::class.java)
        val service = newService(cartRepo, cartItemRepo)

        given(cartRepo.findByUser(userUid, ctx))
            .willReturn(OrderTestFixtures.cartResult(id = 1L))
        given(cartItemRepo.findByProduct(1L, 200L, ctx))
            .willReturn(OrderTestFixtures.cartItemResult(quantity = 2))
        given(cartItemRepo.updateQuantity(1L, OrderTestFixtures.CART_ITEM_UID, 5, ctx))
            .willReturn(OrderTestFixtures.cartItemResult(quantity = 5))

        val command = CartAddItemCommand(
            productId = 200L, productName = "Book A",
            quantity = 3, currency = "KRW", price = BigDecimal("5000.00"),
        )
        val item = runBlocking { service.addItem(userUid, command, ctx) }

        assertEquals(5, item.quantity)
        verify(cartItemRepo).updateQuantity(1L, OrderTestFixtures.CART_ITEM_UID, 5, ctx)
        verify(cartItemRepo, never()).create(1L, command, ctx)
    }

    @Test
    fun `addItem_whenProductNotInCart_createsItem`() {
        val cartRepo = mock(CartRepo::class.java)
        val cartItemRepo = mock(CartItemRepo::class.java)
        val service = newService(cartRepo, cartItemRepo)

        given(cartRepo.findByUser(userUid, ctx))
            .willReturn(OrderTestFixtures.cartResult(id = 1L))
        given(cartItemRepo.findByProduct(1L, 200L, ctx)).willReturn(null)

        val command = CartAddItemCommand(
            productId = 200L, productName = "Book A",
            quantity = 3, currency = "KRW", price = BigDecimal("5000.00"),
        )
        given(cartItemRepo.create(1L, command, ctx))
            .willReturn(OrderTestFixtures.cartItemResult(quantity = 3))

        val item = runBlocking { service.addItem(userUid, command, ctx) }

        assertEquals(3, item.quantity)
        verify(cartItemRepo).create(1L, command, ctx)
    }

    @Test
    fun `updateItemQuantity_whenQuantityZero_removesItem`() {
        val cartRepo = mock(CartRepo::class.java)
        val cartItemRepo = mock(CartItemRepo::class.java)
        val service = newService(cartRepo, cartItemRepo)

        given(cartRepo.findByUser(userUid, ctx))
            .willReturn(OrderTestFixtures.cartResult(id = 1L))
        given(cartItemRepo.findOne(1L, OrderTestFixtures.CART_ITEM_UID, ctx))
            .willReturn(OrderTestFixtures.cartItemResult(quantity = 5))

        val item = runBlocking { service.updateItemQuantity(userUid, OrderTestFixtures.CART_ITEM_UID, 0, ctx) }

        assertEquals(0, item.quantity)
        verify(cartItemRepo).delete(1L, OrderTestFixtures.CART_ITEM_UID, ctx)
        verify(cartItemRepo, never()).updateQuantity(1L, OrderTestFixtures.CART_ITEM_UID, 0, ctx)
    }

    @Test
    fun `updateItemQuantity_whenQuantityPositive_updatesItem`() {
        val cartRepo = mock(CartRepo::class.java)
        val cartItemRepo = mock(CartItemRepo::class.java)
        val service = newService(cartRepo, cartItemRepo)

        given(cartRepo.findByUser(userUid, ctx))
            .willReturn(OrderTestFixtures.cartResult(id = 1L))
        given(cartItemRepo.updateQuantity(1L, OrderTestFixtures.CART_ITEM_UID, 7, ctx))
            .willReturn(OrderTestFixtures.cartItemResult(quantity = 7))

        val item = runBlocking { service.updateItemQuantity(userUid, OrderTestFixtures.CART_ITEM_UID, 7, ctx) }

        assertEquals(7, item.quantity)
        verify(cartItemRepo, never()).delete(1L, OrderTestFixtures.CART_ITEM_UID, ctx)
    }

    @Test
    fun `clearCart_whenCartMissing_writesNothing`() {
        val cartRepo = mock(CartRepo::class.java)
        val cartItemRepo = mock(CartItemRepo::class.java)
        val service = newService(cartRepo, cartItemRepo)

        given(cartRepo.findByUser(userUid, ctx)).willReturn(null)

        runBlocking { service.clearCart(userUid, ctx) }

        verifyNoInteractions(cartItemRepo)
    }

    @Test
    fun `clearCart_whenCartExists_deletesEveryItem`() {
        val cartRepo = mock(CartRepo::class.java)
        val cartItemRepo = mock(CartItemRepo::class.java)
        val service = newService(cartRepo, cartItemRepo)

        given(cartRepo.findByUser(userUid, ctx))
            .willReturn(OrderTestFixtures.cartResult(id = 1L))

        runBlocking { service.clearCart(userUid, ctx) }

        verify(cartItemRepo).deleteAllByCart(1L, ctx)
    }

    @Test
    fun `checkout_whenIdempotencyKeyAlreadyProcessed_returnsExistingOrder`() {
        val cartRepo = mock(CartRepo::class.java)
        val cartItemRepo = mock(CartItemRepo::class.java)
        val orderRepo = mock(OrderRepo::class.java)
        val orderItemRepo = mock(OrderItemRepo::class.java)
        val orderShipRepo = mock(OrderShippingRepo::class.java)
        val lock = FakeDistributedLock()
        val idempotency = FakeIdempotencyService()
        val service = CartService(cartRepo, cartItemRepo, orderRepo, orderItemRepo, orderShipRepo, lock, idempotency, FakeTransactionRunner())

        idempotency.markProcessed("$userUid:idem_001", OrderTestFixtures.ORDER_UID.toString())
        given(orderRepo.findOne(OrderTestFixtures.ORDER_UID, ctx))
            .willReturn(OrderTestFixtures.orderResult())

        val command = CartCheckoutCommand(currency = "KRW", idempotencyKey = "idem_001", shipping = null)
        val agg = runBlocking { service.checkout(userUid, command, ctx) }

        assertEquals(OrderTestFixtures.ORDER_UID, agg.uid)
        assertEquals(0, lock.invocationCount)
        verify(cartRepo, never()).findByUser(userUid, ctx)
    }

    @Test
    fun `checkout_whenCartHasItems_buildsOrderDeletesCartAndMarksIdempotency`() {
        val cartRepo = mock(CartRepo::class.java)
        val cartItemRepo = mock(CartItemRepo::class.java)
        val orderRepo = mock(OrderRepo::class.java)
        val orderItemRepo = mock(OrderItemRepo::class.java)
        val orderShipRepo = mock(OrderShippingRepo::class.java)
        val lock = FakeDistributedLock()
        val idempotency = FakeIdempotencyService()
        val service = CartService(cartRepo, cartItemRepo, orderRepo, orderItemRepo, orderShipRepo, lock, idempotency, FakeTransactionRunner())

        given(cartRepo.findByUser(userUid, ctx))
            .willReturn(OrderTestFixtures.cartResult(id = 1L))
        val items = listOf(
            OrderTestFixtures.cartItemResult(productId = 200L, quantity = 2, price = BigDecimal("5000.00")),
            OrderTestFixtures.cartItemResult(productId = 201L, quantity = 1, price = BigDecimal("3000.00")),
        )
        given(cartItemRepo.findAllByCart(1L, ctx)).willReturn(items)

        val expectedOrderCmd = OrderCreateCommand(
            userUid = userUid,
            currency = "KRW",
            itemsAmount = BigDecimal("13000.00"),
            discountAmount = BigDecimal.ZERO,
            shippingAmount = BigDecimal.ZERO,
            taxAmount = BigDecimal.ZERO,
            totalAmount = BigDecimal("13000.00"),
            idempotencyKey = "$userUid:idem_001",
        )
        given(orderRepo.create(expectedOrderCmd, ctx))
            .willReturn(OrderTestFixtures.orderResult(itemsAmount = BigDecimal("13000.00"), totalAmount = BigDecimal("13000.00")))

        items.forEachIndexed { idx, ci ->
            val expectedItemCmd = OrderItemCreateCommand(
                productId = ci.productId,
                sku = null,
                productName = ci.productName,
                options = null,
                quantity = ci.quantity,
                currency = ci.currency,
                price = ci.price,
                taxRate = BigDecimal.ZERO,
            )
            given(orderItemRepo.create(OrderTestFixtures.ORDER_UID, expectedItemCmd, ctx))
                .willReturn(OrderTestFixtures.orderItemResult(id = (idx + 1).toLong()))
        }

        val command = CartCheckoutCommand(currency = "KRW", idempotencyKey = "idem_001", shipping = null)
        val agg = runBlocking { service.checkout(userUid, command, ctx) }

        assertEquals(2, agg.items.size)
        assertEquals(BigDecimal("13000.00"), agg.totalAmount)
        assertEquals(1, lock.invocationCount)
        assertEquals("checkout:$userUid:idem_001", lock.lastKey)
        assertNull(lock.lastTtl, "checkout relies on watchdog auto-renewal, not a fixed lease")
        assertEquals(FakeDistributedLock.TEST_TOKEN, lock.lastUnlockedToken)
        val cascade = inOrder(cartItemRepo, cartRepo)
        cascade.verify(cartItemRepo).deleteAllByCart(1L, ctx)
        cascade.verify(cartRepo).delete(userUid, ctx)
        assertEquals(OrderTestFixtures.ORDER_UID.toString(), idempotency.getResult("$userUid:idem_001"))
    }

    @Test
    fun `checkout_whenCartMissing_throws`() {
        val cartRepo = mock(CartRepo::class.java)
        val cartItemRepo = mock(CartItemRepo::class.java)
        val orderRepo = mock(OrderRepo::class.java)
        val orderItemRepo = mock(OrderItemRepo::class.java)
        val orderShipRepo = mock(OrderShippingRepo::class.java)
        val service = CartService(
            cartRepo, cartItemRepo, orderRepo, orderItemRepo, orderShipRepo,
            FakeDistributedLock(), FakeIdempotencyService(), FakeTransactionRunner(),
        )

        given(cartRepo.findByUser(userUid, ctx)).willReturn(null)

        val ex = assertFailsWith<CartStateConflictException> {
            runBlocking { service.checkout(userUid, CartCheckoutCommand("KRW", "idem", null), ctx) }
        }
        assertNotNull(ex.message)
    }

    @Test
    fun `checkout_whenCartEmpty_throws`() {
        val cartRepo = mock(CartRepo::class.java)
        val cartItemRepo = mock(CartItemRepo::class.java)
        val orderRepo = mock(OrderRepo::class.java)
        val orderItemRepo = mock(OrderItemRepo::class.java)
        val orderShipRepo = mock(OrderShippingRepo::class.java)
        val service = CartService(
            cartRepo, cartItemRepo, orderRepo, orderItemRepo, orderShipRepo,
            FakeDistributedLock(), FakeIdempotencyService(), FakeTransactionRunner(),
        )

        given(cartRepo.findByUser(userUid, ctx))
            .willReturn(OrderTestFixtures.cartResult(id = 1L))
        given(cartItemRepo.findAllByCart(1L, ctx)).willReturn(emptyList())

        val ex = assertFailsWith<CartStateConflictException> {
            runBlocking { service.checkout(userUid, CartCheckoutCommand("KRW", "idem", null), ctx) }
        }
        assertEquals(true, ex.message?.contains("empty"))
    }

    private fun newService(
        cartRepo: CartRepo,
        cartItemRepo: CartItemRepo,
    ): CartService = CartService(
        cartRepo,
        cartItemRepo,
        mock(OrderRepo::class.java),
        mock(OrderItemRepo::class.java),
        mock(OrderShippingRepo::class.java),
        FakeDistributedLock(),
        FakeIdempotencyService(),
        FakeTransactionRunner(),
    )
}
