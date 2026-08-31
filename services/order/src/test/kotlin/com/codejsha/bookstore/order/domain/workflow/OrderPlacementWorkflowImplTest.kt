package com.codejsha.bookstore.order.domain.workflow

import com.codejsha.bookstore.order.support.FakeInventoryActivities
import com.codejsha.bookstore.order.support.FakeOrderActivities
import com.codejsha.bookstore.order.support.FakePaymentActivities
import io.temporal.client.WorkflowFailedException
import io.temporal.client.WorkflowOptions
import io.temporal.testing.TestWorkflowEnvironment
import org.junit.jupiter.api.AfterEach
import org.junit.jupiter.api.BeforeEach
import org.junit.jupiter.api.Disabled
import org.junit.jupiter.api.Test
import java.math.BigDecimal
import kotlin.test.assertEquals
import kotlin.test.assertFailsWith
import kotlin.test.assertTrue

@Disabled(
    "TestWorkflowEnvironment never dispatches activity tasks in this environment; reproduced unchanged with temporal-testing 1.30.1 on both JDK 25 and JDK 21 (2026-07-30), so a lib upgrade alone does not fix it; see OrderCancellationWorkflowImplTest for the original diagnosis",
)
class OrderPlacementWorkflowImplTest {

    private val orderUid = "44444444-4444-4444-4444-444444444444"
    private val paymentUid = "99999999-9999-9999-9999-999999999999"
    private val userUid = "22222222-2222-2222-2222-222222222222"
    private val taskQueue = "order-task-queue"

    private val items = listOf(
        OrderItemWorkflowRequest(
            productId = 200L,
            productName = "Book A",
            quantity = 2,
            currency = "USD",
            price = BigDecimal("10.00"),
        ),
        OrderItemWorkflowRequest(
            productId = 201L,
            productName = "Book B",
            quantity = 1,
            currency = "USD",
            price = BigDecimal("5.00"),
        ),
    )
    private val expectedTotal = BigDecimal("25.00")

    private val request = PlaceOrderWorkflowRequest(
        userUid = userUid,
        currency = "USD",
        items = items,
        shipping = null,
        idempotencyKey = "idem-123",
    )

    private lateinit var env: TestWorkflowEnvironment
    private lateinit var orderActivities: FakeOrderActivities
    private lateinit var inventoryActivities: FakeInventoryActivities
    private lateinit var paymentActivities: FakePaymentActivities
    private lateinit var workflow: OrderPlacementWorkflow

    @BeforeEach
    fun setUp() {
        env = TestWorkflowEnvironment.newInstance()
        orderActivities = FakeOrderActivities().apply {
            createOrderReturn = CreateOrderResult(orderUid = orderUid, status = "PENDING", paymentUid = null)
        }
        inventoryActivities = FakeInventoryActivities()
        paymentActivities = FakePaymentActivities()

        env.newWorker(taskQueue).apply {
            registerWorkflowImplementationTypes(OrderPlacementWorkflowImpl::class.java)
            registerActivitiesImplementations(orderActivities)
        }
        env.newWorker("inventory-task-queue").registerActivitiesImplementations(inventoryActivities)
        env.newWorker("payment-task-queue").registerActivitiesImplementations(paymentActivities)
        env.start()

        workflow = env.workflowClient.newWorkflowStub(
            OrderPlacementWorkflow::class.java,
            WorkflowOptions.newBuilder()
                .setTaskQueue(taskQueue)
                .setWorkflowId("place-${request.idempotencyKey}")
                .build(),
        )
    }

    @AfterEach
    fun tearDown() {
        env.close()
    }

    @Test
    fun `placeOrder_whenStockAndPaymentSucceed_recordsPaymentAndConfirmsOrder`() {
        paymentActivities.processPaymentReturn = ProcessPaymentResult(paymentUid = paymentUid, status = "succeeded")

        val result = workflow.placeOrder(request)

        assertEquals(orderUid, result.orderUid)
        assertEquals(paymentUid, result.paymentUid)
        assertEquals("PAID", result.status)

        assertEquals(1, orderActivities.createOrderRequests.size)
        val created = orderActivities.createOrderRequests.first()
        assertEquals(0, expectedTotal.compareTo(created.totalAmount))
        assertEquals(0, expectedTotal.compareTo(created.itemsAmount))
        assertEquals("place-${request.idempotencyKey}", created.idempotencyKey)

        assertEquals(
            listOf(
                ReserveStockRequest(
                    orderUid = orderUid,
                    items = listOf(
                        StockReservationItem(productId = 200L, quantity = 2),
                        StockReservationItem(productId = 201L, quantity = 1),
                    ),
                ),
            ),
            inventoryActivities.reserveStockRequests,
        )
        assertEquals(1, paymentActivities.processPaymentRequests.size)
        assertEquals(0, expectedTotal.compareTo(paymentActivities.processPaymentRequests.first().amount))

        assertEquals(listOf(orderUid to paymentUid), orderActivities.recordedPayments)
        assertEquals(listOf(orderUid), orderActivities.confirmedOrderUids)
        assertTrue(inventoryActivities.releaseStockRequests.isEmpty())
        assertTrue(paymentActivities.refundedPaymentUids.isEmpty())
        assertTrue(orderActivities.cancelledOrderUids.isEmpty())
    }

    @Test
    fun `placeOrder_whenPaymentUnsuccessful_releasesStockRefundsAndCancelsOrder`() {
        paymentActivities.processPaymentReturn = ProcessPaymentResult(paymentUid = paymentUid, status = "failed")

        assertFailsWith<WorkflowFailedException> { workflow.placeOrder(request) }

        assertEquals(1, inventoryActivities.reserveStockRequests.size)
        assertEquals(1, inventoryActivities.releaseStockRequests.size)
        assertEquals(listOf(paymentUid), paymentActivities.refundedPaymentUids)
        assertEquals(listOf(orderUid), orderActivities.refundedOrderUids)
        assertTrue(orderActivities.cancelledOrderUids.isEmpty())
        assertTrue(orderActivities.confirmedOrderUids.isEmpty())
    }

    @Test
    fun `placeOrder_whenStockReservationFails_compensatesWithoutChargingPayment`() {
        inventoryActivities.reserveStockFailure = IllegalStateException("out of stock")

        assertFailsWith<WorkflowFailedException> { workflow.placeOrder(request) }

        assertTrue(inventoryActivities.reserveStockRequests.isNotEmpty())
        assertEquals(1, inventoryActivities.releaseStockRequests.size)
        assertTrue(paymentActivities.processPaymentRequests.isEmpty())
        assertTrue(paymentActivities.refundedPaymentUids.isEmpty())
        assertEquals(listOf(orderUid), orderActivities.cancelledOrderUids)
        assertTrue(orderActivities.confirmedOrderUids.isEmpty())
    }

    @Test
    fun `placeOrder_whenOrderAlreadyPaid_returnsIdempotentlyWithoutTouchingStockOrPayment`() {
        orderActivities.createOrderReturn =
            CreateOrderResult(orderUid = orderUid, status = "PAID", paymentUid = paymentUid)

        val result = workflow.placeOrder(request)

        assertEquals(PlaceOrderWorkflowResult(orderUid = orderUid, paymentUid = paymentUid, status = "PAID"), result)
        assertTrue(inventoryActivities.reserveStockRequests.isEmpty())
        assertTrue(paymentActivities.processPaymentRequests.isEmpty())
        assertTrue(paymentActivities.refundedPaymentUids.isEmpty())
        assertTrue(orderActivities.confirmedOrderUids.isEmpty())
        assertTrue(orderActivities.refundedOrderUids.isEmpty())
    }

    @Test
    fun `placeOrder_whenOrderFinalized_failsWithoutCompensating`() {
        orderActivities.createOrderReturn =
            CreateOrderResult(orderUid = orderUid, status = "CANCELLED", paymentUid = null)

        assertFailsWith<WorkflowFailedException> { workflow.placeOrder(request) }

        assertTrue(inventoryActivities.reserveStockRequests.isEmpty())
        assertTrue(inventoryActivities.releaseStockRequests.isEmpty())
        assertTrue(paymentActivities.processPaymentRequests.isEmpty())
        assertTrue(paymentActivities.refundedPaymentUids.isEmpty())
        assertTrue(orderActivities.cancelledOrderUids.isEmpty())
        assertTrue(orderActivities.refundedOrderUids.isEmpty())
    }
}
