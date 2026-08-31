package com.codejsha.bookstore.order.domain.workflow

import com.codejsha.bookstore.order.application.port.repo.OrderItemRepo
import com.codejsha.bookstore.order.application.port.repo.OrderRepo
import com.codejsha.bookstore.order.application.port.repo.OrderShippingRepo
import com.codejsha.bookstore.order.domain.constant.OrderStatus
import com.codejsha.bookstore.order.domain.model.command.OrderCreateCommand
import com.codejsha.bookstore.order.domain.model.command.OrderItemCreateCommand
import com.codejsha.bookstore.order.support.FakeTransactionRunner
import com.codejsha.bookstore.order.support.OrderTestFixtures
import com.codejsha.platform.shared.data.ActorContext
import com.codejsha.platform.shared.data.ActorType
import io.temporal.failure.ApplicationFailure
import org.junit.jupiter.api.Test
import org.mockito.BDDMockito.given
import org.mockito.Mockito.mock
import org.mockito.Mockito.verify
import org.mockito.Mockito.verifyNoInteractions
import org.mockito.Mockito.verifyNoMoreInteractions
import org.springframework.data.domain.PageImpl
import org.springframework.data.domain.Pageable
import java.math.BigDecimal
import kotlin.test.assertEquals
import kotlin.test.assertFailsWith
import kotlin.test.assertTrue

class OrderActivitiesImplTest {

    private val systemContext = ActorContext(actorId = 0L, actorType = ActorType.SYSTEM)

    private val userUid = OrderTestFixtures.USER_UID
    private val userContext = ActorContext(actorId = 0L, actorType = ActorType.USER)

    private fun newActivities(
        orderRepo: OrderRepo,
        itemRepo: OrderItemRepo,
        shipRepo: OrderShippingRepo,
    ) = OrderActivitiesImpl(orderRepo, itemRepo, shipRepo, FakeTransactionRunner())

    // ─── createOrder (idempotency) ───────────────────────────────────────────

    @Test
    fun `createOrder_whenIdempotencyKeyAlreadyPresent_returnsExistingOrderWithoutInserting`() {
        val orderRepo = mock(OrderRepo::class.java)
        val itemRepo = mock(OrderItemRepo::class.java)
        val shipRepo = mock(OrderShippingRepo::class.java)
        val activities = newActivities(orderRepo, itemRepo, shipRepo)

        val existing = OrderTestFixtures.orderResult(uid = OrderTestFixtures.ORDER_UID)
        given(orderRepo.findByIdempotencyKey("place-key-1", userContext)).willReturn(existing)

        val result = activities.createOrder(createRequest(idempotencyKey = "place-key-1"))

        assertEquals(OrderTestFixtures.ORDER_UID.toString(), result.orderUid)
        assertEquals(existing.status, result.status)
        verify(orderRepo).findByIdempotencyKey("place-key-1", userContext)
        verifyNoMoreInteractions(orderRepo)
        verifyNoInteractions(itemRepo, shipRepo)
    }

    @Test
    fun `createOrder_whenIdempotencyKeyUnseen_insertsOrderAndItems`() {
        val orderRepo = mock(OrderRepo::class.java)
        val itemRepo = mock(OrderItemRepo::class.java)
        val shipRepo = mock(OrderShippingRepo::class.java)
        val activities = newActivities(orderRepo, itemRepo, shipRepo)

        val created = OrderTestFixtures.orderResult(uid = OrderTestFixtures.ORDER_UID)
        given(orderRepo.findByIdempotencyKey("place-key-2", userContext)).willReturn(null)
        given(orderRepo.create(expectedOrderCommand("place-key-2"), userContext)).willReturn(created)
        given(itemRepo.create(OrderTestFixtures.ORDER_UID, expectedItemCommand(), userContext))
            .willReturn(OrderTestFixtures.orderItemResult())

        val result = activities.createOrder(createRequest(idempotencyKey = "place-key-2"))

        assertEquals(OrderTestFixtures.ORDER_UID.toString(), result.orderUid)
        verify(orderRepo).create(expectedOrderCommand("place-key-2"), userContext)
        verify(itemRepo).create(OrderTestFixtures.ORDER_UID, expectedItemCommand(), userContext)
    }

    @Test
    fun `createOrder_whenUserUidMalformed_throwsNonRetryableApplicationFailure`() {
        val orderRepo = mock(OrderRepo::class.java)
        val itemRepo = mock(OrderItemRepo::class.java)
        val shipRepo = mock(OrderShippingRepo::class.java)
        val activities = newActivities(orderRepo, itemRepo, shipRepo)

        given(orderRepo.findByIdempotencyKey("place-key-3", userContext)).willReturn(null)

        val failure = assertFailsWith<ApplicationFailure> {
            activities.createOrder(createRequest(idempotencyKey = "place-key-3").copy(userUid = "not-a-uuid"))
        }

        assertTrue(failure.isNonRetryable)
        assertEquals("InvalidCommand", failure.type)
        verifyNoInteractions(itemRepo, shipRepo)
    }

    @Test
    fun `createOrder_whenCommandViolatesDomainConstraint_throwsNonRetryableApplicationFailure`() {
        val orderRepo = mock(OrderRepo::class.java)
        val itemRepo = mock(OrderItemRepo::class.java)
        val shipRepo = mock(OrderShippingRepo::class.java)
        val activities = newActivities(orderRepo, itemRepo, shipRepo)

        given(orderRepo.findByIdempotencyKey("place-key-4", userContext)).willReturn(null)

        val failure = assertFailsWith<ApplicationFailure> {
            activities.createOrder(createRequest(idempotencyKey = "place-key-4").copy(currency = "usd"))
        }

        assertTrue(failure.isNonRetryable)
        assertEquals("InvalidCommand", failure.type)
        verifyNoInteractions(itemRepo, shipRepo)
    }

    private fun createRequest(idempotencyKey: String) = CreateOrderRequest(
        userUid = userUid.toString(),
        currency = "USD",
        itemsAmount = BigDecimal("49.99"),
        totalAmount = BigDecimal("49.99"),
        idempotencyKey = idempotencyKey,
        items = listOf(
            OrderItemWorkflowRequest(
                productId = 200L,
                productName = "Book",
                quantity = 1,
                currency = "USD",
                price = BigDecimal("49.99"),
            ),
        ),
        shipping = null,
    )

    private fun expectedOrderCommand(idempotencyKey: String) = OrderCreateCommand(
        userUid = userUid,
        currency = "USD",
        itemsAmount = BigDecimal("49.99"),
        discountAmount = BigDecimal.ZERO,
        shippingAmount = BigDecimal.ZERO,
        taxAmount = BigDecimal.ZERO,
        totalAmount = BigDecimal("49.99"),
        idempotencyKey = idempotencyKey,
    )

    private fun expectedItemCommand() = OrderItemCreateCommand(
        productId = 200L,
        sku = null,
        productName = "Book",
        options = null,
        quantity = 1,
        currency = "USD",
        price = BigDecimal("49.99"),
        taxRate = BigDecimal.ZERO,
    )

    // ─── loadOrderItems ──────────────────────────────────────────────────────

    @Test
    fun `loadOrderItems_whenOrderHasItems_mapsEachItemToStockReservationItem`() {
        val orderRepo = mock(OrderRepo::class.java)
        val itemRepo = mock(OrderItemRepo::class.java)
        val shipRepo = mock(OrderShippingRepo::class.java)
        val activities = newActivities(orderRepo, itemRepo, shipRepo)

        val items = listOf(
            OrderTestFixtures.orderItemResult(id = 1L, productId = 200L, quantity = 2),
            OrderTestFixtures.orderItemResult(id = 2L, productId = 201L, quantity = 5),
        )
        given(
            itemRepo.findAllByOrder(OrderTestFixtures.ORDER_UID, Pageable.unpaged(), systemContext),
        ).willReturn(PageImpl(items))

        val result = activities.loadOrderItems(OrderTestFixtures.ORDER_UID.toString())

        assertEquals(2, result.size)
        assertEquals(StockReservationItem(productId = 200L, quantity = 2), result[0])
        assertEquals(StockReservationItem(productId = 201L, quantity = 5), result[1])
    }

    @Test
    fun `loadOrderItems_whenOrderHasNoItems_returnsEmptyList`() {
        val orderRepo = mock(OrderRepo::class.java)
        val itemRepo = mock(OrderItemRepo::class.java)
        val shipRepo = mock(OrderShippingRepo::class.java)
        val activities = newActivities(orderRepo, itemRepo, shipRepo)

        given(
            itemRepo.findAllByOrder(OrderTestFixtures.ORDER_UID, Pageable.unpaged(), systemContext),
        ).willReturn(PageImpl(emptyList()))

        val result = activities.loadOrderItems(OrderTestFixtures.ORDER_UID.toString())

        assertEquals(0, result.size)
    }

    // ─── markOrderRefunded ───────────────────────────────────────────────────

    @Test
    fun `markOrderRefunded_whenOrderPaid_transitionsStatusToRefunded`() {
        val orderRepo = mock(OrderRepo::class.java)
        val itemRepo = mock(OrderItemRepo::class.java)
        val shipRepo = mock(OrderShippingRepo::class.java)
        val activities = newActivities(orderRepo, itemRepo, shipRepo)

        val from = setOf(OrderStatus.PENDING.value, OrderStatus.PAID.value, OrderStatus.CANCELLED.value)
        given(orderRepo.transitionStatus(OrderTestFixtures.ORDER_UID, from, OrderStatus.REFUNDED.value, systemContext))
            .willReturn(true)

        activities.markOrderRefunded(OrderTestFixtures.ORDER_UID.toString())

        verify(orderRepo).transitionStatus(
            OrderTestFixtures.ORDER_UID,
            from,
            OrderStatus.REFUNDED.value,
            systemContext,
        )
    }

    // ─── markOrderShipped ────────────────────────────────────────────────────

    @Test
    fun `markOrderShipped_whenOrderPaid_transitionsStatusToShipped`() {
        val orderRepo = mock(OrderRepo::class.java)
        val itemRepo = mock(OrderItemRepo::class.java)
        val shipRepo = mock(OrderShippingRepo::class.java)
        val activities = newActivities(orderRepo, itemRepo, shipRepo)

        given(
            orderRepo.transitionStatus(
                OrderTestFixtures.ORDER_UID,
                setOf(OrderStatus.PAID.value),
                OrderStatus.SHIPPED.value,
                systemContext,
            ),
        ).willReturn(true)

        activities.markOrderShipped(OrderTestFixtures.ORDER_UID.toString())

        verify(orderRepo).transitionStatus(
            OrderTestFixtures.ORDER_UID,
            setOf(OrderStatus.PAID.value),
            OrderStatus.SHIPPED.value,
            systemContext,
        )
    }

    // ─── status transition conflicts ─────────────────────────────────────────

    @Test
    fun `confirmOrder_whenOrderNoLongerPending_throwsNonRetryableApplicationFailure`() {
        val orderRepo = mock(OrderRepo::class.java)
        val itemRepo = mock(OrderItemRepo::class.java)
        val shipRepo = mock(OrderShippingRepo::class.java)
        val activities = newActivities(orderRepo, itemRepo, shipRepo)

        given(
            orderRepo.transitionStatus(
                OrderTestFixtures.ORDER_UID,
                setOf(OrderStatus.PENDING.value),
                OrderStatus.PAID.value,
                systemContext,
            ),
        ).willReturn(false)
        given(orderRepo.findOne(OrderTestFixtures.ORDER_UID, systemContext))
            .willReturn(OrderTestFixtures.orderResult(status = OrderStatus.CANCELLED.value))

        val failure = assertFailsWith<ApplicationFailure> {
            activities.confirmOrder(OrderTestFixtures.ORDER_UID.toString())
        }

        assertTrue(failure.isNonRetryable)
        assertEquals("OrderStateConflict", failure.type)
    }

    @Test
    fun `confirmOrder_whenOrderAlreadyAtTargetStatus_writesNothing`() {
        val orderRepo = mock(OrderRepo::class.java)
        val activities = newActivities(orderRepo, mock(OrderItemRepo::class.java), mock(OrderShippingRepo::class.java))

        given(
            orderRepo.transitionStatus(
                OrderTestFixtures.ORDER_UID,
                setOf(OrderStatus.PENDING.value),
                OrderStatus.PAID.value,
                systemContext,
            ),
        ).willReturn(false)
        given(orderRepo.findOne(OrderTestFixtures.ORDER_UID, systemContext))
            .willReturn(OrderTestFixtures.orderResult(status = OrderStatus.PAID.value))

        activities.confirmOrder(OrderTestFixtures.ORDER_UID.toString())
    }

    @Test
    fun `restoreOrderPaid_whenOrderShippedOrRefunded_transitionsStatusBackToPaid`() {
        val orderRepo = mock(OrderRepo::class.java)
        val activities = newActivities(orderRepo, mock(OrderItemRepo::class.java), mock(OrderShippingRepo::class.java))
        val from = setOf(OrderStatus.SHIPPED.value, OrderStatus.REFUNDED.value)
        given(orderRepo.transitionStatus(OrderTestFixtures.ORDER_UID, from, OrderStatus.PAID.value, systemContext))
            .willReturn(true)

        activities.restoreOrderPaid(OrderTestFixtures.ORDER_UID.toString())

        verify(orderRepo).transitionStatus(OrderTestFixtures.ORDER_UID, from, OrderStatus.PAID.value, systemContext)
    }
}
