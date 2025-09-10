package com.codejsha.bookstore.payment.infrastructure.support.auth

import io.grpc.Context
import io.grpc.Contexts
import io.grpc.Metadata
import io.grpc.ServerCall
import io.grpc.ServerCallHandler
import io.grpc.ServerInterceptor
import io.grpc.Status
import org.springframework.stereotype.Component

@Component
class ActorIdentityServerInterceptor : ServerInterceptor {
    override fun <ReqT, RespT> interceptCall(
        call: ServerCall<ReqT, RespT>,
        headers: Metadata,
        next: ServerCallHandler<ReqT, RespT>,
    ): ServerCall.Listener<ReqT> {
        val identity = ActorIdentity.fromMetadata(headers)
        if (identity == null) {
            call.close(
                Status.UNAUTHENTICATED.withDescription("missing or invalid actor identity"),
                Metadata(),
            )
            return object : ServerCall.Listener<ReqT>() {}
        }
        val context = Context.current().withValue(ActorIdentity.CONTEXT_KEY, identity)
        return Contexts.interceptCall(context, call, headers, next)
    }
}
