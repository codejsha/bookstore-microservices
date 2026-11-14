package com.codejsha.bookstore.payment.infrastructure.support

import io.grpc.Metadata
import io.grpc.MethodDescriptor
import io.grpc.ServerCall
import io.grpc.Status
import java.io.ByteArrayInputStream
import java.io.InputStream

object GrpcInterceptorTestSupport {

    private object StringMarshaller : MethodDescriptor.Marshaller<String> {
        override fun stream(value: String): InputStream = ByteArrayInputStream(value.toByteArray())

        override fun parse(stream: InputStream): String = String(stream.readAllBytes())
    }

    fun unaryMethod(fullMethodName: String): MethodDescriptor<String, String> =
        MethodDescriptor.newBuilder(StringMarshaller, StringMarshaller)
            .setType(MethodDescriptor.MethodType.UNARY)
            .setFullMethodName(fullMethodName)
            .build()

    class RecordingServerCall(
        private val descriptor: MethodDescriptor<String, String>,
    ) : ServerCall<String, String>() {
        var closedStatus: Status? = null
        var closedTrailers: Metadata? = null
        val sentMessages = mutableListOf<String>()

        override fun request(numMessages: Int) = Unit

        override fun sendHeaders(headers: Metadata) = Unit

        override fun sendMessage(message: String) {
            sentMessages.add(message)
        }

        override fun close(status: Status, trailers: Metadata) {
            closedStatus = status
            closedTrailers = trailers
        }

        override fun isCancelled(): Boolean = false

        override fun getMethodDescriptor(): MethodDescriptor<String, String> = descriptor
    }
}
