package com.codejsha.bookstore.payment.infrastructure.support

import io.temporal.activity.Activity
import io.temporal.activity.ActivityExecutionContext
import org.junit.jupiter.api.Test
import org.mockito.BDDMockito.willThrow
import org.mockito.Mockito.mock
import org.mockito.Mockito.mockStatic
import org.mockito.Mockito.verify
import org.mockito.Mockito.verifyNoMoreInteractions
import kotlin.test.assertFailsWith
import kotlin.test.assertSame

class TemporalActivityHeartbeatTest {

    private val heartbeat = TemporalActivityHeartbeat()

    @Test
    fun `beat_outsideActivityContext_skipsWithoutThrowing`() {
        heartbeat.beat("capture pay_001")
    }

    @Test
    fun `beat_insideActivityContext_reportsDetailToTemporal`() {
        val context = mock(ActivityExecutionContext::class.java)

        mockStatic(Activity::class.java).use { activity ->
            activity.`when`<ActivityExecutionContext> { Activity.getExecutionContext() }.thenReturn(context)

            heartbeat.beat("capture pay_001")
        }

        verify(context).heartbeat("capture pay_001")
        verifyNoMoreInteractions(context)
    }

    @Test
    fun `beat_heartbeatSignalsCancellation_propagatesFailureToActivity`() {
        val context = mock(ActivityExecutionContext::class.java)
        val cancellation = IllegalArgumentException("activity cancelled")
        willThrow(cancellation).given(context).heartbeat("refund ref_001")

        val thrown = mockStatic(Activity::class.java).use { activity ->
            activity.`when`<ActivityExecutionContext> { Activity.getExecutionContext() }.thenReturn(context)

            assertFailsWith<IllegalArgumentException> { heartbeat.beat("refund ref_001") }
        }

        assertSame(cancellation, thrown)
    }
}
