package com.codejsha.bookstore.order.domain.workflow

import com.codejsha.bookstore.order.support.FakeDeliveryActivities
import com.codejsha.bookstore.order.support.FakeNotificationActivities
import com.codejsha.bookstore.order.support.FakeOrderActivities
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
class OrderFulfillmentWorkflowImplTest {

    private val orderUid = "11111111-1111-1111-1111-111111111111"
    private val userUid = "22222222-2222-2222-2222-222222222222"
    private val shipmentUid = "33333333-3333-3333-3333-333333333333"
    private val trackingNumber = "TRACK-12345"
    private val taskQueue = "order-task-queue"

    private val destination = ShipmentDestination(
        addressLine = "1 Foo St",
        city = "Seoul",
        state = "Seoul",
        countryCode = "KR",
        postalCode = "12345",
    )

    private lateinit var env: TestWorkflowEnvironment
    private lateinit var orderActivities: FakeOrderActivities
    private lateinit var deliveryActivities: FakeDeliveryActivities
    private lateinit var notificationActivities: FakeNotificationActivities
    private lateinit var workflow: OrderFulfillmentWorkflow

    @BeforeEach
    fun setUp() {
        env = TestWorkflowEnvironment.newInstance()
        orderActivities = FakeOrderActivities()
        deliveryActivities = FakeDeliveryActivities()
        notificationActivities = FakeNotificationActivities()

        env.newWorker(taskQueue).apply {
            registerWorkflowImplementationTypes(OrderFulfillmentWorkflowImpl::class.java)
            registerActivitiesImplementations(orderActivities)
        }
        env.newWorker("delivery-task-queue").registerActivitiesImplementations(deliveryActivities)
        env.newWorker("notification-task-queue").registerActivitiesImplementations(notificationActivities)
        env.start()

        workflow = env.workflowClient.newWorkflowStub(
            OrderFulfillmentWorkflow::class.java,
            WorkflowOptions.newBuilder().setTaskQueue(taskQueue).build(),
        )
    }

    @AfterEach
    fun tearDown() {
        env.close()
    }

    @Test
    fun `creates shipment, notifies customer with tracking number, marks shipped`() {
        orderActivities.cancellationState = OrderCancellationState(status = "PAID", paymentUid = null)
        deliveryActivities.createShipmentReturn = CreateShipmentResult(
            shipmentUid = shipmentUid,
            trackingNumber = trackingNumber,
            status = "CREATED",
        )

        val result = workflow.fulfillOrder(
            FulfillOrderWorkflowRequest(orderUid = orderUid, userUid = userUid, destination = destination),
        )

        assertEquals("SHIPPED", result.status)
        assertEquals(orderUid, result.orderUid)
        assertEquals(shipmentUid, result.shipmentUid)

        assertEquals(
            listOf(
                CreateShipmentRequest(
                    orderUid = orderUid,
                    originAddress = "warehouse",
                    destinationAddress = destination.addressLine,
                    destinationCity = destination.city,
                    destinationState = destination.state,
                    destinationCountryCode = destination.countryCode,
                    destinationPostalCode = destination.postalCode,
                ),
            ),
            deliveryActivities.createShipmentRequests,
        )
        assertEquals(
            listOf(
                ShipmentNotificationRequest(
                    userUid = userUid,
                    orderUid = orderUid,
                    notificationType = "SHIPMENT_DISPATCHED",
                    trackingNumber = trackingNumber,
                ),
            ),
            notificationActivities.sentNotifications,
        )
        assertEquals(listOf(orderUid), orderActivities.shippedOrderUids)
    }

    @Test
    fun `refuses to fulfill an order that is not PAID — no shipment, notification, or status change`() {
        orderActivities.cancellationState = OrderCancellationState(status = "PENDING", paymentUid = null)

        assertFailsWith<WorkflowFailedException> {
            workflow.fulfillOrder(
                FulfillOrderWorkflowRequest(orderUid = orderUid, userUid = userUid, destination = destination),
            )
        }

        assertTrue(deliveryActivities.createShipmentRequests.isEmpty())
        assertTrue(notificationActivities.sentNotifications.isEmpty())
        assertTrue(orderActivities.shippedOrderUids.isEmpty())
    }
}
