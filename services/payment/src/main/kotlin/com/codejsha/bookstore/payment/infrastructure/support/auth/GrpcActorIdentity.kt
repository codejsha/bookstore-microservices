package com.codejsha.bookstore.payment.infrastructure.support.auth

import io.grpc.Context
import io.grpc.Metadata
import io.grpc.Status
import io.grpc.StatusRuntimeException
import java.util.UUID

data class ActorIdentity(
    val uid: UUID,
    val admin: Boolean,
) {
    val customerUid: String get() = uid.toString()

    companion object {
        val UID_KEY: Metadata.Key<String> = Metadata.Key.of("x-actor-uid", Metadata.ASCII_STRING_MARSHALLER)
        val ADMIN_KEY: Metadata.Key<String> = Metadata.Key.of("x-actor-admin", Metadata.ASCII_STRING_MARSHALLER)

        val CONTEXT_KEY: Context.Key<ActorIdentity> = Context.key("payment-actor-identity")

        fun fromMetadata(headers: Metadata): ActorIdentity? {
            val uid = headers.get(UID_KEY)
                ?.trim()
                ?.takeIf { it.isNotEmpty() }
                ?.let { runCatching { UUID.fromString(it) }.getOrNull() }
                ?: return null
            val admin = headers.get(ADMIN_KEY)?.trim().equals("true", ignoreCase = true)
            return ActorIdentity(uid, admin)
        }

        fun require(): ActorIdentity =
            CONTEXT_KEY.get()
                ?: throw StatusRuntimeException(
                    Status.UNAUTHENTICATED.withDescription("missing or invalid actor identity"),
                )
    }
}
