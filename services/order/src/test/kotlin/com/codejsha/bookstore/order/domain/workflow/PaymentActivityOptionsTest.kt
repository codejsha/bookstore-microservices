package com.codejsha.bookstore.order.domain.workflow

import org.junit.jupiter.api.Test
import java.time.Duration
import kotlin.test.assertEquals
import kotlin.test.assertTrue

class PaymentActivityOptionsTest {

    @Test
    fun `placementPaymentActivityOptions_gatewayCallsHeartbeat_setsHeartbeatTimeout`() {
        val options = OrderPlacementWorkflowImpl.paymentActivityOptions()

        assertEquals(Duration.ofSeconds(30), options.heartbeatTimeout)
        assertEquals(Duration.ofSeconds(60), options.startToCloseTimeout)
        assertEquals("payment-task-queue", options.taskQueue)
    }

    @Test
    fun `placementCompensationPaymentActivityOptions_gatewayCallsHeartbeat_setsHeartbeatTimeout`() {
        val options = OrderPlacementWorkflowImpl.compensationPaymentActivityOptions()

        assertEquals(Duration.ofSeconds(30), options.heartbeatTimeout)
        assertEquals(Duration.ofSeconds(60), options.startToCloseTimeout)
    }

    @Test
    fun `cancellationPaymentActivityOptions_gatewayCallsHeartbeat_setsHeartbeatTimeout`() {
        val options = OrderCancellationWorkflowImpl.paymentActivityOptions()

        assertEquals(Duration.ofSeconds(30), options.heartbeatTimeout)
        assertEquals(Duration.ofSeconds(60), options.startToCloseTimeout)
    }

    @Test
    fun `paymentActivityOptions_hyperswitchTimeoutBudget_keepsHeartbeatBelowStartToClose`() {
        val options = OrderPlacementWorkflowImpl.paymentActivityOptions()
        val worstCaseGatewayCall = HYPERSWITCH_CONNECT_TIMEOUT.plus(HYPERSWITCH_READ_TIMEOUT)

        assertTrue(
            options.heartbeatTimeout > worstCaseGatewayCall,
            "heartbeat timeout ${options.heartbeatTimeout} must outlast a $worstCaseGatewayCall gateway call",
        )
        assertTrue(
            options.heartbeatTimeout < options.startToCloseTimeout,
            "heartbeat timeout must be shorter than the start-to-close budget",
        )
    }

    private companion object {
        val HYPERSWITCH_CONNECT_TIMEOUT: Duration = Duration.ofSeconds(2)
        val HYPERSWITCH_READ_TIMEOUT: Duration = Duration.ofSeconds(15)
    }
}
