package com.codejsha.bookstore.payment.infrastructure.support.auth

import com.codejsha.bookstore.payment.infrastructure.support.GrpcInterceptorTestSupport.RecordingServerCall
import com.codejsha.bookstore.payment.infrastructure.support.GrpcInterceptorTestSupport.unaryMethod
import io.grpc.Metadata
import io.grpc.ServerCall
import io.grpc.ServerCallHandler
import io.grpc.Status
import io.grpc.StatusRuntimeException
import org.junit.jupiter.api.Test
import java.util.UUID
import kotlin.test.assertEquals
import kotlin.test.assertFailsWith
import kotlin.test.assertFalse
import kotlin.test.assertNotNull
import kotlin.test.assertNull
import kotlin.test.assertSame
import kotlin.test.assertTrue

class ActorIdentityServerInterceptorTest {

    private val interceptor = ActorIdentityServerInterceptor()

    private class IdentityCapturingHandler : ServerCallHandler<String, String> {
        var started = false
        var headers: Metadata? = null
        var identityAtStart: ActorIdentity? = null
        var identityAtHalfClose: ActorIdentity? = null

        override fun startCall(call: ServerCall<String, String>, headers: Metadata): ServerCall.Listener<String> {
            started = true
            this.headers = headers
            identityAtStart = ActorIdentity.CONTEXT_KEY.get()
            return object : ServerCall.Listener<String>() {
                override fun onHalfClose() {
                    identityAtHalfClose = ActorIdentity.CONTEXT_KEY.get()
                }
            }
        }
    }

    private fun newCall() = RecordingServerCall(unaryMethod("payment.v1.PaymentService/FindPayment"))

    private fun metadata(uid: String? = null, admin: String? = null) = Metadata().apply {
        uid?.let { put(ActorIdentity.UID_KEY, it) }
        admin?.let { put(ActorIdentity.ADMIN_KEY, it) }
    }

    @Test
    fun `interceptCall_validActorUid_exposesIdentityToHandlerAndListenerCallbacks`() {
        val call = newCall()
        val handler = IdentityCapturingHandler()
        val headers = metadata(uid = ACTOR_UID.toString())

        val listener = interceptor.interceptCall(call, headers, handler)
        listener.onHalfClose()

        assertEquals(ActorIdentity(ACTOR_UID, admin = false), handler.identityAtStart)
        assertEquals(ActorIdentity(ACTOR_UID, admin = false), handler.identityAtHalfClose)
        assertSame(headers, handler.headers)
        assertNull(call.closedStatus)
    }

    @Test
    fun `interceptCall_validActorUid_leavesCallerContextWithoutIdentity`() {
        interceptor.interceptCall(newCall(), metadata(uid = ACTOR_UID.toString()), IdentityCapturingHandler())

        assertNull(ActorIdentity.CONTEXT_KEY.get())
    }

    @Test
    fun `interceptCall_adminFlagTrueInAnyCase_marksIdentityAdmin`() {
        listOf("true", "TRUE", " True ").forEach { flag ->
            val handler = IdentityCapturingHandler()

            interceptor.interceptCall(newCall(), metadata(uid = ACTOR_UID.toString(), admin = flag), handler)

            assertEquals(true, handler.identityAtStart?.admin, "admin flag '$flag'")
        }
    }

    @Test
    fun `interceptCall_adminFlagNotTrue_keepsIdentityNonAdmin`() {
        listOf("false", "1", "yes", "").forEach { flag ->
            val handler = IdentityCapturingHandler()

            interceptor.interceptCall(newCall(), metadata(uid = ACTOR_UID.toString(), admin = flag), handler)

            assertEquals(false, handler.identityAtStart?.admin, "admin flag '$flag'")
        }
    }

    @Test
    fun `interceptCall_missingActorUid_closesUnauthenticatedWithoutStartingHandler`() {
        val call = newCall()
        val handler = IdentityCapturingHandler()

        val listener = interceptor.interceptCall(call, metadata(admin = "true"), handler)
        listener.onMessage("ignored")
        listener.onHalfClose()

        val status = assertNotNull(call.closedStatus)
        assertEquals(Status.Code.UNAUTHENTICATED, status.code)
        assertEquals("missing or invalid actor identity", status.description)
        assertTrue(assertNotNull(call.closedTrailers).keys().isEmpty())
        assertFalse(handler.started)
        assertNull(handler.identityAtHalfClose)
    }

    @Test
    fun `interceptCall_malformedOrBlankActorUid_closesUnauthenticatedWithoutStartingHandler`() {
        listOf("not-a-uuid", "   ", "").forEach { uid ->
            val call = newCall()
            val handler = IdentityCapturingHandler()

            interceptor.interceptCall(call, metadata(uid = uid), handler)

            assertEquals(Status.Code.UNAUTHENTICATED, call.closedStatus?.code, "uid '$uid'")
            assertFalse(handler.started, "uid '$uid'")
        }
    }

    @Test
    fun `fromMetadata_uidPaddedWithWhitespace_parsesTrimmedUid`() {
        val identity = ActorIdentity.fromMetadata(metadata(uid = "  $ACTOR_UID  "))

        assertEquals(ACTOR_UID, assertNotNull(identity).uid)
        assertEquals(ACTOR_UID.toString(), identity.customerUid)
        assertFalse(identity.admin)
    }

    @Test
    fun `require_noIdentityInContext_throwsUnauthenticated`() {
        val ex = assertFailsWith<StatusRuntimeException> { ActorIdentity.require() }

        assertEquals(Status.Code.UNAUTHENTICATED, ex.status.code)
        assertEquals("missing or invalid actor identity", ex.status.description)
    }

    @Test
    fun `require_insideInterceptedCall_yieldsIdentityFromMetadata`() {
        var required: ActorIdentity? = null
        val handler = ServerCallHandler<String, String> { _, _ ->
            required = ActorIdentity.require()
            object : ServerCall.Listener<String>() {}
        }

        interceptor.interceptCall(newCall(), metadata(uid = ACTOR_UID.toString(), admin = "true"), handler)

        assertEquals(ActorIdentity(ACTOR_UID, admin = true), required)
    }

    private companion object {
        val ACTOR_UID: UUID = UUID.fromString("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa")
    }
}
