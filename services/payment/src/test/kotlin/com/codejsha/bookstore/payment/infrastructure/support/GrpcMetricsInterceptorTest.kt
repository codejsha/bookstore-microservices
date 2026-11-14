package com.codejsha.bookstore.payment.infrastructure.support

import com.codejsha.bookstore.payment.infrastructure.support.GrpcInterceptorTestSupport.RecordingServerCall
import com.codejsha.bookstore.payment.infrastructure.support.GrpcInterceptorTestSupport.unaryMethod
import io.grpc.Metadata
import io.grpc.ServerCall
import io.grpc.ServerCallHandler
import io.grpc.Status
import io.micrometer.core.instrument.simple.SimpleMeterRegistry
import org.junit.jupiter.api.Test
import kotlin.test.assertEquals
import kotlin.test.assertNotNull
import kotlin.test.assertNull
import kotlin.test.assertSame
import kotlin.test.assertTrue

class GrpcMetricsInterceptorTest {

    private val registry = SimpleMeterRegistry()
    private val interceptor = GrpcMetricsInterceptor(registry)

    private class CapturingHandler : ServerCallHandler<String, String> {
        var call: ServerCall<String, String>? = null
        var headers: Metadata? = null
        val listener = object : ServerCall.Listener<String>() {}

        override fun startCall(call: ServerCall<String, String>, headers: Metadata): ServerCall.Listener<String> {
            this.call = call
            this.headers = headers
            return listener
        }
    }

    private fun timer(service: String, method: String, statusCode: String) =
        registry.find("rpc.server.duration")
            .tag("rpc.system", "grpc")
            .tag("rpc.service", service)
            .tag("rpc.method", method)
            .tag("rpc.grpc.status_code", statusCode)
            .timer()

    @Test
    fun `interceptCall_callClosedOk_recordsTimerTaggedWithServiceMethodAndStatusZero`() {
        val call = RecordingServerCall(unaryMethod("payment.v1.PaymentService/FindPayment"))
        val handler = CapturingHandler()
        val trailers = Metadata()

        interceptor.interceptCall(call, Metadata(), handler)
        assertNotNull(handler.call).close(Status.OK, trailers)

        val timer = assertNotNull(timer("payment.v1.PaymentService", "FindPayment", "0"))
        assertEquals(1L, timer.count())
        assertSame(Status.OK, call.closedStatus)
        assertSame(trailers, call.closedTrailers)
    }

    @Test
    fun `interceptCall_callClosedWithError_tagsNumericGrpcStatusCodeAndForwardsClose`() {
        val call = RecordingServerCall(unaryMethod("payment.v1.PaymentService/ListRefunds"))
        val handler = CapturingHandler()
        val status = Status.NOT_FOUND.withDescription("refund not found")

        interceptor.interceptCall(call, Metadata(), handler)
        assertNotNull(handler.call).close(status, Metadata())

        assertEquals(1L, assertNotNull(timer("payment.v1.PaymentService", "ListRefunds", "5")).count())
        assertNull(timer("payment.v1.PaymentService", "ListRefunds", "0"))
        assertSame(status, call.closedStatus)
    }

    @Test
    fun `interceptCall_callStillOpen_recordsNoTimer`() {
        val call = RecordingServerCall(unaryMethod("payment.v1.PaymentService/FindPayment"))

        interceptor.interceptCall(call, Metadata(), CapturingHandler())

        assertTrue(registry.find("rpc.server.duration").timers().isEmpty())
        assertNull(call.closedStatus)
    }

    @Test
    fun `interceptCall_anyCall_passesHeadersAndListenerThroughAndForwardsMessages`() {
        val call = RecordingServerCall(unaryMethod("payment.v1.PaymentService/FindPayment"))
        val handler = CapturingHandler()
        val headers = Metadata()

        val listener = interceptor.interceptCall(call, headers, handler)
        assertNotNull(handler.call).sendMessage("payload")

        assertSame(handler.listener, listener)
        assertSame(headers, handler.headers)
        assertEquals(listOf("payload"), call.sentMessages)
    }

    @Test
    fun `interceptCall_twoCallsOfOneMethod_accumulateOnOneTimer`() {
        repeat(2) {
            val call = RecordingServerCall(unaryMethod("payment.v1.PaymentService/FindPayment"))
            val handler = CapturingHandler()
            interceptor.interceptCall(call, Metadata(), handler)
            assertNotNull(handler.call).close(Status.OK, Metadata())
        }

        assertEquals(2L, assertNotNull(timer("payment.v1.PaymentService", "FindPayment", "0")).count())
        assertEquals(1, registry.find("rpc.server.duration").timers().size)
    }

    @Test
    fun `interceptCall_methodNameWithoutServicePrefix_usesEmptyServiceAndWholeNameAsMethod`() {
        val call = RecordingServerCall(unaryMethod("Ping"))
        val handler = CapturingHandler()

        interceptor.interceptCall(call, Metadata(), handler)
        assertNotNull(handler.call).close(Status.OK, Metadata())

        assertEquals(1L, assertNotNull(timer("", "Ping", "0")).count())
    }
}
