package com.codejsha.bookstore.payment.infrastructure.support

import io.grpc.ForwardingServerCall
import io.grpc.Metadata
import io.grpc.ServerCall
import io.grpc.ServerCallHandler
import io.grpc.ServerInterceptor
import io.grpc.Status
import io.micrometer.core.instrument.MeterRegistry
import io.micrometer.core.instrument.Timer
import org.springframework.stereotype.Component

@Component
class GrpcMetricsInterceptor(
    private val meterRegistry: MeterRegistry,
) : ServerInterceptor {
    override fun <ReqT, RespT> interceptCall(
        call: ServerCall<ReqT, RespT>,
        headers: Metadata,
        next: ServerCallHandler<ReqT, RespT>,
    ): ServerCall.Listener<ReqT> {
        val fullMethod = call.methodDescriptor.fullMethodName
        val service = fullMethod.substringBeforeLast('/', "")
        val method = fullMethod.substringAfterLast('/', fullMethod)
        val sample = Timer.start(meterRegistry)

        val monitoringCall =
            object : ForwardingServerCall.SimpleForwardingServerCall<ReqT, RespT>(call) {
                override fun close(status: Status, trailers: Metadata) {
                    sample.stop(
                        Timer
                            .builder(TIMER_NAME)
                            .tag("rpc.system", "grpc")
                            .tag("rpc.service", service)
                            .tag("rpc.method", method)
                            .tag("rpc.grpc.status_code", status.code.value().toString())
                            .register(meterRegistry),
                    )
                    super.close(status, trailers)
                }
            }

        return next.startCall(monitoringCall, headers)
    }

    companion object {
        private const val TIMER_NAME = "rpc.server.duration"
    }
}
