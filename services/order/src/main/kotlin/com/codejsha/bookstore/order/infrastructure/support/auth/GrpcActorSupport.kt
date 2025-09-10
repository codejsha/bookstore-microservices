package com.codejsha.bookstore.order.infrastructure.support.auth

import io.grpc.Context
import io.grpc.Contexts
import io.grpc.Metadata
import io.grpc.ServerCall
import io.grpc.ServerCallHandler
import io.grpc.ServerInterceptor
import io.grpc.Status
import org.springframework.stereotype.Component
import java.util.UUID

val ACTOR_UID_METADATA_KEY: Metadata.Key<String> =
    Metadata.Key.of("x-actor-uid", Metadata.ASCII_STRING_MARSHALLER)

val ACTOR_ADMIN_METADATA_KEY: Metadata.Key<String> =
    Metadata.Key.of("x-actor-admin", Metadata.ASCII_STRING_MARSHALLER)

internal val ACTOR_UID_CONTEXT_KEY: Context.Key<UUID> = Context.key("order-actor-uid")
internal val ACTOR_ADMIN_CONTEXT_KEY: Context.Key<Boolean> = Context.key("order-actor-admin")

data class GrpcActor(val uid: UUID, val isAdmin: Boolean)

fun currentGrpcActor(): GrpcActor =
    GrpcActor(
        uid = requireNotNull(ACTOR_UID_CONTEXT_KEY.get()) { "actor uid missing from gRPC context" },
        isAdmin = ACTOR_ADMIN_CONTEXT_KEY.get() ?: false,
    )

@Component
class ActorAuthInterceptor : ServerInterceptor {
    override fun <ReqT, RespT> interceptCall(
        call: ServerCall<ReqT, RespT>,
        headers: Metadata,
        next: ServerCallHandler<ReqT, RespT>,
    ): ServerCall.Listener<ReqT> {
        val actorUid = headers.get(ACTOR_UID_METADATA_KEY)
            ?.takeIf { it.isNotBlank() }
            ?.let { runCatching { UUID.fromString(it) }.getOrNull() }
        if (actorUid == null) {
            call.close(
                Status.UNAUTHENTICATED.withDescription("missing or invalid actor identity"),
                Metadata(),
            )
            return object : ServerCall.Listener<ReqT>() {}
        }
        val isAdmin = headers.get(ACTOR_ADMIN_METADATA_KEY) == "true"
        val ctx = Context.current()
            .withValue(ACTOR_UID_CONTEXT_KEY, actorUid)
            .withValue(ACTOR_ADMIN_CONTEXT_KEY, isAdmin)
        return Contexts.interceptCall(ctx, call, headers, next)
    }
}
