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
import kotlin.test.assertEquals
import kotlin.test.assertFailsWith
import kotlin.test.assertTrue

@Disabled(
    "TestWorkflowEnvironment never dispatches activity tasks in this environment; reproduced unchanged with temporal-testing 1.30.1 on both JDK 25 and JDK 21 (2026-07-30), so a lib upgrade alone does not fix it; see OrderCancellationWorkflowImplTest for the original diagnosis",
)
class OrderCancellationWorkflowImplTest {

    private val orderUid = "11111111-1111-1111-1111-111111111111"
    private val paymentUid = "99999999-9999-9999-9999-999999999999"
    private val taskQueue = "order-task-queue"

    private lateinit var env: TestWorkflowEnvironment
    private lateinit var orderActivities: FakeOrderActivities
    private lateinit var inventoryActivities: FakeInventoryActivities
    private lateinit var paymentActivities: FakePaymentActivities
    private lateinit var workflow: OrderCancellationWorkflow

    @BeforeEach
    fun setUp() {
        env = TestWorkflowEnvironment.newInstance()
        orderActivities = FakeOrderActivities()
        inventoryActivities = FakeInventoryActivities()
        paymentActivities = FakePaymentActivities()

        env.newWorker(taskQueue).apply {
            registerWorkflowImplementationTypes(OrderCancellationWorkflowImpl::class.java)
            registerActivitiesImplementations(orderActivities)
        }
        env.newWorker("inventory-task-queue").registerActivitiesImplementations(inventoryActivities)
        env.newWorker("payment-task-queue").registerActivitiesImplementations(paymentActivities)
        env.start()

        workflow = env.workflowClient.newWorkflowStub(
            OrderCancellationWorkflow::class.java,
            WorkflowOptions.newBuilder().setTaskQueue(taskQueue).build(),
        )
    }

    @AfterEach
    fun tearDown() {
        env.close()
    }

    @Test
    fun `cancelOrder_whenOrderPaid_releasesStockRefundsPaymentAndReturnsRefunded`() {
        val items = listOf(
            StockReservationItem(productId = 200L, quantity = 2),
            StockReservationItem(productId = 201L, quantity = 1),
        )
        orderActivities.cancellationState = OrderCancellationState(status = "PAID", paymentUid = paymentUid)
        orderActivities.orderItems = items

        val result = workflow.cancelOrder(
            CancelOrderWorkflowRequest(orderUid = orderUid, reason = "user requested"),
        )

        assertEquals("REFUNDED", result.status)
        assertEquals(orderUid, result.orderUid)

        assertEquals(
            listOf(ReserveStockRequest(orderUid = orderUid, items = items)),
            inventoryActivities.releaseStockRequests,
        )
        assertEquals(listOf(paymentUid), paymentActivities.refundedPaymentUids)
        assertEquals(listOf(orderUid), orderActivities.refundedOrderUids)
        assertTrue(orderActivities.cancelledOrderUids.isEmpty())
    }

    @Test
    fun `cancelOrder_whenOrderPending_releasesStockAndReturnsCancelledWithoutRefund`() {
        val items = listOf(StockReservationItem(productId = 200L, quantity = 2))
        orderActivities.cancellationState = OrderCancellationState(status = "PENDING", paymentUid = null)
        orderActivities.orderItems = items

        val result = workflow.cancelOrder(
            CancelOrderWorkflowRequest(orderUid = orderUid, reason = null),
        )

        assertEquals("CANCELLED", result.status)
        assertEquals(orderUid, result.orderUid)

        assertEquals(
            listOf(ReserveStockRequest(orderUid = orderUid, items = items)),
            inventoryActivities.releaseStockRequests,
        )
        assertEquals(listOf(orderUid), orderActivities.cancelledOrderUids)
        assertTrue(paymentActivities.refundedPaymentUids.isEmpty())
        assertTrue(orderActivities.refundedOrderUids.isEmpty())
    }

    @Test
    fun `cancelOrder_whenOrderHasNoItems_skipsStockRelease`() {
        orderActivities.cancellationState = OrderCancellationState(status = "PAID", paymentUid = paymentUid)
        orderActivities.orderItems = emptyList()

        val result = workflow.cancelOrder(
            CancelOrderWorkflowRequest(orderUid = orderUid, reason = null),
        )

        assertEquals("REFUNDED", result.status)
        assertTrue(inventoryActivities.releaseStockRequests.isEmpty())
        assertEquals(listOf(paymentUid), paymentActivities.refundedPaymentUids)
        assertEquals(listOf(orderUid), orderActivities.refundedOrderUids)
    }

    @Test
    fun `cancelOrder_whenPaymentNotRefundable_restoresPaidAndFailsWithoutReleasingStock`() {
        orderActivities.cancellationState = OrderCancellationState(status = "PAID", paymentUid = paymentUid)
        orderActivities.orderItems = listOf(StockReservationItem(productId = 200L, quantity = 2))
        paymentActivities.refundPaymentReturn = false

        assertFailsWith<WorkflowFailedException> {
            workflow.cancelOrder(CancelOrderWorkflowRequest(orderUid = orderUid, reason = null))
        }

        assertEquals(listOf(paymentUid), paymentActivities.refundedPaymentUids)
        assertTrue(inventoryActivities.releaseStockRequests.isEmpty(), "stock must stay reserved when the refund fails")
        assertEquals(listOf(orderUid), orderActivities.refundedOrderUids, "the order is claimed before the refund")
        assertEquals(listOf(orderUid), orderActivities.restoredPaidOrderUids, "and restored to PAID when the refund is refused")
        assertTrue(orderActivities.cancelledOrderUids.isEmpty())
    }

    @Test
    fun `cancelOrder_whenOrderShipped_writesNothing`() {
        orderActivities.cancellationState = OrderCancellationState(status = "SHIPPED", paymentUid = paymentUid)

        assertFailsWith<WorkflowFailedException> {
            workflow.cancelOrder(CancelOrderWorkflowRequest(orderUid = orderUid, reason = null))
        }

        assertTrue(inventoryActivities.releaseStockRequests.isEmpty())
        assertTrue(paymentActivities.refundedPaymentUids.isEmpty())
        assertTrue(orderActivities.refundedOrderUids.isEmpty())
        assertTrue(orderActivities.cancelledOrderUids.isEmpty())
    }
}
