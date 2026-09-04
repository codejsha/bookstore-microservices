package com.codejsha.bookstore.order.domain.service

import com.codejsha.bookstore.order.application.port.repo.OrderAdjustmentRepo
import com.codejsha.bookstore.order.application.port.repo.OrderItemRepo
import com.codejsha.bookstore.order.application.port.repo.OrderRepo
import com.codejsha.bookstore.order.application.port.repo.OrderShippingRepo
import com.codejsha.bookstore.order.domain.constant.OrderStatus
import com.codejsha.bookstore.order.domain.model.OrderStateConflictException
import com.codejsha.bookstore.order.domain.model.command.OrderAdjustmentCreateCommand
import com.codejsha.bookstore.order.domain.model.command.OrderCreateCommand
import com.codejsha.bookstore.order.domain.model.command.OrderItemCreateCommand
import com.codejsha.bookstore.order.domain.model.command.OrderItemUpdateCommand
import com.codejsha.bookstore.order.domain.model.command.OrderShippingCreateCommand
import com.codejsha.bookstore.order.domain.model.command.OrderShippingUpdateCommand
import com.codejsha.bookstore.order.domain.model.command.OrderUpdateCommand
import com.codejsha.bookstore.order.domain.model.option.OrderQueryOption
import com.codejsha.bookstore.order.support.FakeTransactionRunner
import com.codejsha.bookstore.order.support.OrderTestFixtures
import kotlinx.coroutines.runBlocking
import org.junit.jupiter.api.Test
import org.mockito.BDDMockito.given
import org.mockito.Mockito.mock
import org.mockito.Mockito.verify
import org.springframework.data.domain.PageImpl
import org.springframework.data.domain.PageRequest
import org.springframework.data.domain.Pageable
import java.math.BigDecimal
import kotlin.test.assertEquals
import kotlin.test.assertFailsWith
import kotlin.test.assertNotNull
import kotlin.test.assertNull

class OrderServiceTest {

    private val ctx = OrderTestFixtures.DEFAULT_CONTEXT

    private fun newService(
        orderRepo: OrderRepo,
        itemRepo: OrderItemRepo,
        adjRepo: OrderAdjustmentRepo,
        shipRepo: OrderShippingRepo,
    ) = OrderService(orderRepo, itemRepo, adjRepo, shipRepo, FakeTransactionRunner())

    @Test
    fun `findAllOrders_whenRepoReturnsRows_mapsEachRowToAggregate`() {
        val orderRepo = mock(OrderRepo::class.java)
        val itemRepo = mock(OrderItemRepo::class.java)
        val adjRepo = mock(OrderAdjustmentRepo::class.java)
        val shipRepo = mock(OrderShippingRepo::class.java)
        val service = newService(orderRepo, itemRepo, adjRepo, shipRepo)

        val pageable = PageRequest.of(0, 10)
        val results = listOf(
            OrderTestFixtures.orderResult(id = 1L, status = OrderStatus.PENDING.value),
            OrderTestFixtures.orderResult(id = 2L, status = OrderStatus.PAID.value),
        )
        val option = OrderQueryOption()
        given(orderRepo.findAll(option, pageable, ctx))
            .willReturn(PageImpl(results, pageable, 2L))

        val page = runBlocking { service.findAllOrders(option, pageable, ctx) }

        assertEquals(2, page.content.size)
        assertEquals(OrderStatus.PENDING, page.content[0].status)
        assertEquals(OrderStatus.PAID, page.content[1].status)
    }

    @Test
    fun `findOrder_whenOrderExists_composesItemsAdjustmentsAndShipping`() {
        val orderRepo = mock(OrderRepo::class.java)
        val itemRepo = mock(OrderItemRepo::class.java)
        val adjRepo = mock(OrderAdjustmentRepo::class.java)
        val shipRepo = mock(OrderShippingRepo::class.java)
        val service = newService(orderRepo, itemRepo, adjRepo, shipRepo)

        given(orderRepo.findOne(OrderTestFixtures.ORDER_UID, ctx))
            .willReturn(OrderTestFixtures.orderResult())
        given(itemRepo.findAllByOrder(OrderTestFixtures.ORDER_UID, Pageable.unpaged(), ctx))
            .willReturn(PageImpl(listOf(OrderTestFixtures.orderItemResult())))
        given(adjRepo.findAllByOrder(OrderTestFixtures.ORDER_UID, Pageable.unpaged(), ctx))
            .willReturn(PageImpl(listOf(OrderTestFixtures.orderAdjustmentResult())))
        given(shipRepo.findByOrder(OrderTestFixtures.ORDER_UID, ctx))
            .willReturn(OrderTestFixtures.orderShippingResult())

        val agg = runBlocking { service.findOrder(OrderTestFixtures.ORDER_UID, ctx) }

        assertEquals(1, agg.items.size)
        assertEquals(1, agg.adjustments.size)
        assertNotNull(agg.shipping)
        assertEquals("Alice", agg.shipping.recipientName)
    }

    @Test
    fun `findOrder_whenShippingMissing_returnsOrderWithoutShipping`() {
        val orderRepo = mock(OrderRepo::class.java)
        val itemRepo = mock(OrderItemRepo::class.java)
        val adjRepo = mock(OrderAdjustmentRepo::class.java)
        val shipRepo = mock(OrderShippingRepo::class.java)
        val service = newService(orderRepo, itemRepo, adjRepo, shipRepo)

        given(orderRepo.findOne(OrderTestFixtures.ORDER_UID, ctx))
            .willReturn(OrderTestFixtures.orderResult())
        given(itemRepo.findAllByOrder(OrderTestFixtures.ORDER_UID, Pageable.unpaged(), ctx))
            .willReturn(PageImpl(emptyList()))
        given(adjRepo.findAllByOrder(OrderTestFixtures.ORDER_UID, Pageable.unpaged(), ctx))
            .willReturn(PageImpl(emptyList()))
        given(shipRepo.findByOrder(OrderTestFixtures.ORDER_UID, ctx)).willReturn(null)

        val agg = runBlocking { service.findOrder(OrderTestFixtures.ORDER_UID, ctx) }

        assertNull(agg.shipping)
    }

    @Test
    fun `placeOrder_whenCommandValid_createsOrderItemsAndShipping`() {
        val orderRepo = mock(OrderRepo::class.java)
        val itemRepo = mock(OrderItemRepo::class.java)
        val adjRepo = mock(OrderAdjustmentRepo::class.java)
        val shipRepo = mock(OrderShippingRepo::class.java)
        val service = newService(orderRepo, itemRepo, adjRepo, shipRepo)

        val command = OrderCreateCommand(
            userUid = OrderTestFixtures.USER_UID,
            currency = "KRW",
            itemsAmount = BigDecimal("10000.00"),
            discountAmount = BigDecimal.ZERO,
            shippingAmount = BigDecimal.ZERO,
            taxAmount = BigDecimal.ZERO,
            totalAmount = BigDecimal("10000.00"),
            idempotencyKey = "idem_001",
        )
        val itemCmd1 = OrderItemCreateCommand(
            productId = 200L, sku = null, productName = "Book A", options = null,
            quantity = 1, currency = "KRW", price = BigDecimal("5000.00"), taxRate = BigDecimal.ZERO,
        )
        val itemCmd2 = itemCmd1.copy(productId = 201L, productName = "Book B")
        val shippingCmd = OrderShippingCreateCommand(
            recipientName = "Alice", recipientPhone = "010-1",
            addressLine1 = "1 Foo St", addressLine2 = null,
            city = "Seoul", state = "Seoul", postalCode = "12345",
            country = "KR", shippingMethod = "STANDARD",
        )

        given(orderRepo.create(command, ctx)).willReturn(OrderTestFixtures.orderResult())
        given(itemRepo.create(OrderTestFixtures.ORDER_UID, itemCmd1, ctx))
            .willReturn(OrderTestFixtures.orderItemResult(id = 1L))
        given(itemRepo.create(OrderTestFixtures.ORDER_UID, itemCmd2, ctx))
            .willReturn(OrderTestFixtures.orderItemResult(id = 2L))
        given(shipRepo.create(OrderTestFixtures.ORDER_UID, shippingCmd, ctx))
            .willReturn(OrderTestFixtures.orderShippingResult())

        val agg = runBlocking { service.placeOrder(command, listOf(itemCmd1, itemCmd2), shippingCmd, ctx) }

        assertEquals(2, agg.items.size)
        assertNotNull(agg.shipping)
        verify(orderRepo).create(command, ctx)
        verify(itemRepo).create(OrderTestFixtures.ORDER_UID, itemCmd1, ctx)
        verify(itemRepo).create(OrderTestFixtures.ORDER_UID, itemCmd2, ctx)
        verify(shipRepo).create(OrderTestFixtures.ORDER_UID, shippingCmd, ctx)
    }

    @Test
    fun `placeOrder_whenShippingNull_skipsShipping`() {
        val orderRepo = mock(OrderRepo::class.java)
        val itemRepo = mock(OrderItemRepo::class.java)
        val adjRepo = mock(OrderAdjustmentRepo::class.java)
        val shipRepo = mock(OrderShippingRepo::class.java)
        val service = newService(orderRepo, itemRepo, adjRepo, shipRepo)

        val command = OrderCreateCommand(
            userUid = OrderTestFixtures.USER_UID, currency = "KRW",
            itemsAmount = BigDecimal.ZERO, discountAmount = BigDecimal.ZERO,
            shippingAmount = BigDecimal.ZERO, taxAmount = BigDecimal.ZERO,
            totalAmount = BigDecimal.ZERO, idempotencyKey = "idem",
        )
        given(orderRepo.create(command, ctx)).willReturn(OrderTestFixtures.orderResult())

        val agg = runBlocking { service.placeOrder(command, emptyList(), shipping = null, context = ctx) }

        assertNull(agg.shipping)
        assertEquals(0, agg.items.size)
    }

    @Test
    fun `cancelOrder_whenStatusPending_marksOrderCancelled`() {
        val orderRepo = mock(OrderRepo::class.java)
        val itemRepo = mock(OrderItemRepo::class.java)
        val adjRepo = mock(OrderAdjustmentRepo::class.java)
        val shipRepo = mock(OrderShippingRepo::class.java)
        val service = newService(orderRepo, itemRepo, adjRepo, shipRepo)

        given(orderRepo.findOne(OrderTestFixtures.ORDER_UID, ctx))
            .willReturn(OrderTestFixtures.orderResult(status = OrderStatus.PENDING.value))

        val expectedUpdate = OrderUpdateCommand(
            userUid = null,
            status = OrderStatus.CANCELLED.value,
            currency = null, itemsAmount = null,
            discountAmount = null, shippingAmount = null,
            taxAmount = null, totalAmount = null,
        )
        given(orderRepo.update(OrderTestFixtures.ORDER_UID, expectedUpdate, ctx))
            .willReturn(OrderTestFixtures.orderResult(status = OrderStatus.CANCELLED.value))

        val agg = runBlocking { service.cancelOrder(OrderTestFixtures.ORDER_UID, ctx) }

        assertEquals(OrderStatus.CANCELLED, agg.status)
        verify(orderRepo).update(OrderTestFixtures.ORDER_UID, expectedUpdate, ctx)
    }

    @Test
    fun `cancelOrder_whenStatusNotPending_throws`() {
        val orderRepo = mock(OrderRepo::class.java)
        val itemRepo = mock(OrderItemRepo::class.java)
        val adjRepo = mock(OrderAdjustmentRepo::class.java)
        val shipRepo = mock(OrderShippingRepo::class.java)
        val service = newService(orderRepo, itemRepo, adjRepo, shipRepo)

        given(orderRepo.findOne(OrderTestFixtures.ORDER_UID, ctx))
            .willReturn(OrderTestFixtures.orderResult(status = OrderStatus.PAID.value))

        val ex = assertFailsWith<OrderStateConflictException> {
            runBlocking { service.cancelOrder(OrderTestFixtures.ORDER_UID, ctx) }
        }
        assertEquals(true, ex.message?.contains("PENDING"))
    }

    @Test
    fun `addItem_whenOrderNotPending_throws`() {
        val orderRepo = mock(OrderRepo::class.java)
        val itemRepo = mock(OrderItemRepo::class.java)
        val adjRepo = mock(OrderAdjustmentRepo::class.java)
        val shipRepo = mock(OrderShippingRepo::class.java)
        val service = newService(orderRepo, itemRepo, adjRepo, shipRepo)

        given(orderRepo.findOne(OrderTestFixtures.ORDER_UID, ctx))
            .willReturn(OrderTestFixtures.orderResult(status = OrderStatus.SHIPPED.value))

        val cmd = OrderItemCreateCommand(
            productId = 1L, sku = null, productName = null, options = null,
            quantity = 1, currency = "KRW", price = BigDecimal.ONE, taxRate = BigDecimal.ZERO,
        )
        assertFailsWith<OrderStateConflictException> {
            runBlocking { service.addItem(OrderTestFixtures.ORDER_UID, cmd, ctx) }
        }
    }

    @Test
    fun `addItem_whenOrderPending_createsItem`() {
        val orderRepo = mock(OrderRepo::class.java)
        val itemRepo = mock(OrderItemRepo::class.java)
        val adjRepo = mock(OrderAdjustmentRepo::class.java)
        val shipRepo = mock(OrderShippingRepo::class.java)
        val service = newService(orderRepo, itemRepo, adjRepo, shipRepo)

        given(orderRepo.findOne(OrderTestFixtures.ORDER_UID, ctx))
            .willReturn(OrderTestFixtures.orderResult())

        val cmd = OrderItemCreateCommand(
            productId = 200L, sku = null, productName = "Book A", options = null,
            quantity = 1, currency = "KRW", price = BigDecimal("5000.00"), taxRate = BigDecimal.ZERO,
        )
        given(itemRepo.create(OrderTestFixtures.ORDER_UID, cmd, ctx))
            .willReturn(OrderTestFixtures.orderItemResult())

        val item = runBlocking { service.addItem(OrderTestFixtures.ORDER_UID, cmd, ctx) }

        assertEquals(OrderTestFixtures.ORDER_ITEM_UID, item.uid)
    }

    // Both item mutations share the PENDING guard, so they are asserted together.
    @Test
    fun `updateItemAndRemoveItem_whenOrderNotPending_throw`() {
        val orderRepo = mock(OrderRepo::class.java)
        val itemRepo = mock(OrderItemRepo::class.java)
        val adjRepo = mock(OrderAdjustmentRepo::class.java)
        val shipRepo = mock(OrderShippingRepo::class.java)
        val service = newService(orderRepo, itemRepo, adjRepo, shipRepo)

        given(orderRepo.findOne(OrderTestFixtures.ORDER_UID, ctx))
            .willReturn(OrderTestFixtures.orderResult(status = OrderStatus.DELIVERED.value))

        val updateCmd = OrderItemUpdateCommand(
            productId = null, sku = null, productName = null, options = null,
            quantity = 5, currency = null, price = null, taxRate = null,
        )

        assertFailsWith<OrderStateConflictException> {
            runBlocking { service.updateItem(OrderTestFixtures.ORDER_UID, OrderTestFixtures.ORDER_ITEM_UID, updateCmd, ctx) }
        }
        assertFailsWith<OrderStateConflictException> {
            runBlocking { service.removeItem(OrderTestFixtures.ORDER_UID, OrderTestFixtures.ORDER_ITEM_UID, ctx) }
        }
    }

    @Test
    fun `setShipping_whenShippingMissing_createsShipping`() {
        val orderRepo = mock(OrderRepo::class.java)
        val itemRepo = mock(OrderItemRepo::class.java)
        val adjRepo = mock(OrderAdjustmentRepo::class.java)
        val shipRepo = mock(OrderShippingRepo::class.java)
        val service = newService(orderRepo, itemRepo, adjRepo, shipRepo)

        val cmd = OrderShippingCreateCommand(
            recipientName = "Alice", recipientPhone = "010-1",
            addressLine1 = "1", addressLine2 = null,
            city = "Seoul", state = "Seoul", postalCode = "12345",
            country = "KR", shippingMethod = "STANDARD",
        )
        given(orderRepo.findOne(OrderTestFixtures.ORDER_UID, ctx))
            .willReturn(OrderTestFixtures.orderResult())
        given(shipRepo.findByOrder(OrderTestFixtures.ORDER_UID, ctx)).willReturn(null)
        given(shipRepo.create(OrderTestFixtures.ORDER_UID, cmd, ctx))
            .willReturn(OrderTestFixtures.orderShippingResult())

        val shipping = runBlocking { service.setShipping(OrderTestFixtures.ORDER_UID, cmd, ctx) }

        assertEquals(OrderTestFixtures.ORDER_SHIPPING_UID, shipping.uid)
        verify(shipRepo).create(OrderTestFixtures.ORDER_UID, cmd, ctx)
    }

    @Test
    fun `setShipping_whenShippingExists_updatesShipping`() {
        val orderRepo = mock(OrderRepo::class.java)
        val itemRepo = mock(OrderItemRepo::class.java)
        val adjRepo = mock(OrderAdjustmentRepo::class.java)
        val shipRepo = mock(OrderShippingRepo::class.java)
        val service = newService(orderRepo, itemRepo, adjRepo, shipRepo)

        val cmd = OrderShippingCreateCommand(
            recipientName = "Bob", recipientPhone = "010-2",
            addressLine1 = "2 Bar St", addressLine2 = null,
            city = "Busan", state = "Busan", postalCode = "67890",
            country = "KR", shippingMethod = "EXPRESS",
        )
        val expectedUpdate = OrderShippingUpdateCommand(
            recipientName = cmd.recipientName,
            recipientPhone = cmd.recipientPhone,
            addressLine1 = cmd.addressLine1,
            addressLine2 = cmd.addressLine2,
            city = cmd.city,
            state = cmd.state,
            postalCode = cmd.postalCode,
            country = cmd.country,
            shippingMethod = cmd.shippingMethod,
        )

        given(orderRepo.findOne(OrderTestFixtures.ORDER_UID, ctx))
            .willReturn(OrderTestFixtures.orderResult())
        given(shipRepo.findByOrder(OrderTestFixtures.ORDER_UID, ctx))
            .willReturn(OrderTestFixtures.orderShippingResult())
        given(shipRepo.update(OrderTestFixtures.ORDER_UID, expectedUpdate, ctx))
            .willReturn(OrderTestFixtures.orderShippingResult())

        runBlocking { service.setShipping(OrderTestFixtures.ORDER_UID, cmd, ctx) }

        verify(shipRepo).update(OrderTestFixtures.ORDER_UID, expectedUpdate, ctx)
    }

    // removeAdjustment is asserted in the same pass: it delegates straight to the repo.
    @Test
    fun `applyAdjustment_whenCommandValid_createsAdjustment`() {
        val orderRepo = mock(OrderRepo::class.java)
        val itemRepo = mock(OrderItemRepo::class.java)
        val adjRepo = mock(OrderAdjustmentRepo::class.java)
        val shipRepo = mock(OrderShippingRepo::class.java)
        val service = newService(orderRepo, itemRepo, adjRepo, shipRepo)

        val cmd = OrderAdjustmentCreateCommand(
            type = "COUPON", label = "WELCOME",
            amount = BigDecimal("-1000.00"), meta = null,
        )
        given(orderRepo.findOne(OrderTestFixtures.ORDER_UID, ctx))
            .willReturn(OrderTestFixtures.orderResult())
        given(adjRepo.create(OrderTestFixtures.ORDER_UID, cmd, ctx))
            .willReturn(OrderTestFixtures.orderAdjustmentResult())

        val adjustment = runBlocking { service.applyAdjustment(OrderTestFixtures.ORDER_UID, cmd, ctx) }
        assertEquals(OrderTestFixtures.ORDER_ADJUSTMENT_UID, adjustment.uid)

        runBlocking { service.removeAdjustment(OrderTestFixtures.ORDER_UID, OrderTestFixtures.ORDER_ADJUSTMENT_UID, ctx) }
        verify(adjRepo).delete(OrderTestFixtures.ORDER_UID, OrderTestFixtures.ORDER_ADJUSTMENT_UID, ctx)
    }
}
