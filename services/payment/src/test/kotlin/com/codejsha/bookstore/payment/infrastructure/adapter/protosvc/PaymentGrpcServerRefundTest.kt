package com.codejsha.bookstore.payment.infrastructure.adapter.protosvc

import com.codejsha.bookstore.generated.application.port.pb.paymentpb.FindRefundRequest
import com.codejsha.bookstore.generated.application.port.pb.paymentpb.FindRefundResponse
import com.codejsha.bookstore.generated.application.port.pb.paymentpb.ListRefundsRequest
import com.codejsha.bookstore.generated.application.port.pb.paymentpb.ListRefundsResponse
import com.codejsha.bookstore.payment.application.usecase.PaymentUseCase
import com.codejsha.bookstore.payment.application.usecase.RefundUseCase
import com.codejsha.bookstore.payment.domain.aggregate.PaymentAggregate
import com.codejsha.bookstore.payment.domain.aggregate.PaymentAttemptEntity
import com.codejsha.bookstore.payment.domain.aggregate.RefundAggregate
import com.codejsha.bookstore.payment.domain.constant.RefundStatus
import com.codejsha.bookstore.payment.domain.model.command.PaymentCreateCommand
import com.codejsha.bookstore.payment.domain.model.command.PaymentUpdateCommand
import com.codejsha.bookstore.payment.domain.model.command.RefundCreateCommand
import com.codejsha.bookstore.payment.domain.model.option.PaymentQueryOption
import com.codejsha.bookstore.payment.domain.model.option.RefundQueryOption
import com.codejsha.bookstore.payment.infrastructure.support.auth.ActorIdentity
import com.codejsha.bookstore.payment.support.PaymentTestFixtures
import com.codejsha.platform.shared.data.ActorContext
import io.grpc.Context
import io.grpc.Status
import io.grpc.StatusRuntimeException
import org.junit.jupiter.api.Test
import org.springframework.data.domain.Page
import org.springframework.data.domain.PageImpl
import org.springframework.data.domain.PageRequest
import org.springframework.data.domain.Pageable
import java.util.UUID
import kotlin.test.assertEquals
import kotlin.test.assertNull
import kotlin.test.assertTrue
import com.codejsha.bookstore.generated.application.port.pb.paymentpb.RefundStatus as ProtoRefundStatus

class PaymentGrpcServerRefundTest {

    private companion object {
        val ACTOR_UID: UUID = UUID.fromString("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa")
        val OTHER_UID: UUID = UUID.fromString("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb")
        val ADMIN_UID: UUID = UUID.fromString("cccccccc-cccc-cccc-cccc-cccccccccccc")
    }

    // ─── Hand-rolled fakes ──────────────────────────────────────────────────

    private class FakeRefundUseCase : RefundUseCase {
        var findRefundResult: RefundAggregate? = null
        var findAllResult: Page<RefundAggregate> = PageImpl(emptyList())

        var capturedFindUid: UUID? = null
        var capturedOption: RefundQueryOption? = null
        var capturedPageable: Pageable? = null
        var findRefundError: Throwable? = null

        override suspend fun findAllRefunds(
            option: RefundQueryOption,
            pageable: Pageable,
            context: ActorContext,
        ): Page<RefundAggregate> {
            capturedOption = option
            capturedPageable = pageable
            return findAllResult
        }

        override suspend fun findRefund(uid: UUID, context: ActorContext): RefundAggregate {
            capturedFindUid = uid
            findRefundError?.let { throw it }
            return findRefundResult ?: error("findRefundResult not seeded")
        }

        override suspend fun createRefund(
            command: RefundCreateCommand,
            context: ActorContext,
        ): RefundAggregate = error("not used in these tests")
    }

    private class FakePaymentUseCase : PaymentUseCase {
        var ownerResult: PaymentAggregate? = null
        var capturedPaymentId: String? = null

        override suspend fun findPaymentByPaymentId(paymentId: String, context: ActorContext): PaymentAggregate? {
            capturedPaymentId = paymentId
            return ownerResult
        }

        override suspend fun findAllPayments(option: PaymentQueryOption, pageable: Pageable, context: ActorContext): Page<PaymentAggregate> =
            error("not used")
        override suspend fun findPayment(uid: UUID, context: ActorContext): PaymentAggregate = error("not used")
        override suspend fun createPayment(command: PaymentCreateCommand, context: ActorContext): PaymentAggregate = error("not used")
        override suspend fun updatePayment(uid: UUID, command: PaymentUpdateCommand, context: ActorContext): PaymentAggregate = error("not used")
        override suspend fun findAllPaymentAttempts(paymentUid: UUID, pageable: Pageable, context: ActorContext): Page<PaymentAttemptEntity> =
            error("not used")
        override suspend fun findPaymentAttempt(paymentUid: UUID, uid: UUID, context: ActorContext): PaymentAttemptEntity = error("not used")
    }

    private class CapturingObserver<T> : io.grpc.stub.StreamObserver<T> {
        val values = mutableListOf<T>()
        var completed = false
        var error: Throwable? = null

        override fun onNext(value: T) { values.add(value) }
        override fun onError(t: Throwable) { error = t }
        override fun onCompleted() { completed = true }
    }

    private val refundUseCase = FakeRefundUseCase()
    private val paymentUseCase = FakePaymentUseCase()
    private val server = PaymentGrpcServer(paymentUseCase, refundUseCase)

    private inline fun withActor(uid: UUID, admin: Boolean, block: () -> Unit) {
        val ctx = Context.current().withValue(ActorIdentity.CONTEXT_KEY, ActorIdentity(uid, admin))
        val previous = ctx.attach()
        try {
            block()
        } finally {
            ctx.detach(previous)
        }
    }

    // ─── actor guard ────────────────────────────────────────────────────────

    @Test
    fun `findRefund without an actor identity fails UNAUTHENTICATED`() {
        val observer = CapturingObserver<FindRefundResponse>()
        server.findRefund(
            FindRefundRequest.newBuilder().setUid(PaymentTestFixtures.REFUND_UID.toString()).build(),
            observer,
        )

        val status = Status.fromThrowable(observer.error!!)
        assertEquals(Status.Code.UNAUTHENTICATED, status.code)
        assertEquals("missing or invalid actor identity", status.description)
    }

    // ─── findRefund parent-payment ownership ──────────────────────────────────

    @Test
    fun `findRefund gates a non-admin on the parent payment's owner`() {
        val uid = PaymentTestFixtures.REFUND_UID
        refundUseCase.findRefundResult = PaymentTestFixtures.refundAggregate(uid = uid, paymentId = "pay_001")
        paymentUseCase.ownerResult = PaymentTestFixtures.paymentAggregate(customerId = OTHER_UID.toString())

        val observer = CapturingObserver<FindRefundResponse>()
        withActor(ACTOR_UID, admin = false) {
            server.findRefund(FindRefundRequest.newBuilder().setUid(uid.toString()).build(), observer)
        }

        val status = Status.fromThrowable(observer.error!!)
        assertEquals(Status.Code.NOT_FOUND, status.code)
        assertEquals("Refund with uid $uid not found", status.description)
        assertEquals("pay_001", paymentUseCase.capturedPaymentId)
    }

    @Test
    fun `findRefund returns NOT_FOUND for a non-admin when the parent payment cannot be resolved`() {
        refundUseCase.findRefundResult = PaymentTestFixtures.refundAggregate()
        paymentUseCase.ownerResult = null

        val observer = CapturingObserver<FindRefundResponse>()
        withActor(ACTOR_UID, admin = false) {
            server.findRefund(FindRefundRequest.newBuilder().setUid(PaymentTestFixtures.REFUND_UID.toString()).build(), observer)
        }

        assertEquals(Status.Code.NOT_FOUND, Status.fromThrowable(observer.error!!).code)
    }

    @Test
    fun `findRefund lets a non-admin read a refund on their own payment`() {
        refundUseCase.findRefundResult = PaymentTestFixtures.refundAggregate()
        paymentUseCase.ownerResult = PaymentTestFixtures.paymentAggregate(customerId = ACTOR_UID.toString())

        val observer = CapturingObserver<FindRefundResponse>()
        withActor(ACTOR_UID, admin = false) {
            server.findRefund(FindRefundRequest.newBuilder().setUid(PaymentTestFixtures.REFUND_UID.toString()).build(), observer)
        }

        assertNull(observer.error)
        assertEquals(1, observer.values.size)
    }

    // ─── findRefund mapping (admin bypasses ownership) ────────────────────────

    @Test
    fun `findRefund maps aggregate fields onto the Refund proto`() {
        val aggregate = PaymentTestFixtures.refundAggregate(
            amount = 5_000L,
            currency = "KRW",
            status = RefundStatus.SUCCEEDED,
            reason = "requested_by_customer",
        )
        refundUseCase.findRefundResult = aggregate

        val observer = CapturingObserver<FindRefundResponse>()
        withActor(ADMIN_UID, admin = true) {
            server.findRefund(FindRefundRequest.newBuilder().setUid(aggregate.uid.toString()).build(), observer)
        }

        assertEquals(1, observer.values.size, "onNext should be called exactly once")
        assertTrue(observer.completed, "onCompleted should be called")
        assertNull(observer.error, "onError should not be called")

        assertEquals(aggregate.uid, refundUseCase.capturedFindUid)

        val refund = observer.values.single().refund
        assertEquals(aggregate.uid.toString(), refund.uid)
        assertEquals("pay_001", refund.paymentId)
        assertEquals(5_000L, refund.amount)
        assertEquals("KRW", refund.currency)
        assertEquals(ProtoRefundStatus.REFUND_STATUS_SUCCEEDED, refund.status)
        assertEquals("requested_by_customer", refund.reason)
        assertEquals("instant", refund.refundType)
        assertEquals(aggregate.createdAt.toString(), refund.createdAt)
    }

    @Test
    fun `findRefund maps null errorCode and errorMessage to empty strings`() {
        refundUseCase.findRefundResult = PaymentTestFixtures.refundAggregate(
            errorCode = null,
            errorMessage = null,
        )

        val observer = CapturingObserver<FindRefundResponse>()
        withActor(ADMIN_UID, admin = true) {
            server.findRefund(FindRefundRequest.newBuilder().setUid(PaymentTestFixtures.REFUND_UID.toString()).build(), observer)
        }

        val refund = observer.values.single().refund
        assertEquals("", refund.errorCode)
        assertEquals("", refund.errorMessage)
    }

    @Test
    fun `findRefund surfaces usecase failure via onError without onCompleted`() {
        refundUseCase.findRefundError = IllegalStateException("boom")

        val observer = CapturingObserver<FindRefundResponse>()
        withActor(ADMIN_UID, admin = true) {
            server.findRefund(FindRefundRequest.newBuilder().setUid(PaymentTestFixtures.REFUND_UID.toString()).build(), observer)
        }

        assertTrue(observer.values.isEmpty(), "no value should be emitted on failure")
        assertTrue(observer.completed.not(), "onCompleted should not be called on failure")
        assertTrue(observer.error is StatusRuntimeException, "error should be a StatusRuntimeException")
    }

    // ─── listRefunds scoping ──────────────────────────────────────────────────

    @Test
    fun `listRefunds refuses an unscoped non-admin listing with PERMISSION_DENIED`() {
        val observer = CapturingObserver<ListRefundsResponse>()
        withActor(ACTOR_UID, admin = false) {
            server.listRefunds(ListRefundsRequest.newBuilder().build(), observer)
        }

        assertEquals(Status.Code.PERMISSION_DENIED, Status.fromThrowable(observer.error!!).code)
    }

    @Test
    fun `listRefunds lets a non-admin list refunds of their own payment`() {
        refundUseCase.findAllResult = PageImpl(emptyList())
        paymentUseCase.ownerResult = PaymentTestFixtures.paymentAggregate(customerId = ACTOR_UID.toString())

        val observer = CapturingObserver<ListRefundsResponse>()
        withActor(ACTOR_UID, admin = false) {
            server.listRefunds(ListRefundsRequest.newBuilder().setPaymentId("pay_001").build(), observer)
        }

        assertNull(observer.error)
        assertEquals("pay_001", refundUseCase.capturedOption?.paymentId)
        assertEquals("pay_001", paymentUseCase.capturedPaymentId)
    }

    @Test
    fun `listRefunds hides another customer's payment refunds from a non-admin as NOT_FOUND`() {
        paymentUseCase.ownerResult = PaymentTestFixtures.paymentAggregate(customerId = OTHER_UID.toString())

        val observer = CapturingObserver<ListRefundsResponse>()
        withActor(ACTOR_UID, admin = false) {
            server.listRefunds(ListRefundsRequest.newBuilder().setPaymentId("pay_001").build(), observer)
        }

        assertEquals(Status.Code.NOT_FOUND, Status.fromThrowable(observer.error!!).code)
    }

    // ─── listRefunds mapping (admin) ──────────────────────────────────────────

    @Test
    fun `listRefunds maps page content and totalSize`() {
        refundUseCase.findAllResult = PageImpl(
            listOf(
                PaymentTestFixtures.refundAggregate(uid = PaymentTestFixtures.REFUND_UID, amount = 1_000L, status = RefundStatus.PENDING),
                PaymentTestFixtures.refundAggregate(id = 2L, uid = UUID.fromString("99999999-9999-9999-9999-999999999999"), amount = 2_000L, status = RefundStatus.SUCCEEDED),
            ),
            PageRequest.of(0, 10),
            2L,
        )

        val observer = CapturingObserver<ListRefundsResponse>()
        withActor(ADMIN_UID, admin = true) {
            server.listRefunds(ListRefundsRequest.newBuilder().build(), observer)
        }

        assertTrue(observer.completed)
        assertNull(observer.error)
        val response = observer.values.single()
        assertEquals(2, response.refundsList.size)
        assertEquals(2, response.totalSize)
        assertEquals(1_000L, response.refundsList[0].amount)
        assertEquals(ProtoRefundStatus.REFUND_STATUS_PENDING, response.refundsList[0].status)
        assertEquals(ProtoRefundStatus.REFUND_STATUS_SUCCEEDED, response.refundsList[1].status)
    }

    @Test
    fun `listRefunds applies paymentId filter and a concrete status filter`() {
        refundUseCase.findAllResult = PageImpl(emptyList())

        val observer = CapturingObserver<ListRefundsResponse>()
        withActor(ADMIN_UID, admin = true) {
            server.listRefunds(
                ListRefundsRequest.newBuilder()
                    .setPaymentId("pay_001")
                    .setStatus(ProtoRefundStatus.REFUND_STATUS_FAILED)
                    .build(),
                observer,
            )
        }

        val option = refundUseCase.capturedOption!!
        assertEquals("pay_001", option.paymentId)
        assertEquals(RefundStatus.FAILED.value, option.status)
    }

    @Test
    fun `listRefunds does not apply an UNSPECIFIED status filter`() {
        refundUseCase.findAllResult = PageImpl(emptyList())

        val observer = CapturingObserver<ListRefundsResponse>()
        withActor(ADMIN_UID, admin = true) {
            server.listRefunds(ListRefundsRequest.newBuilder().build(), observer)
        }

        val option = refundUseCase.capturedOption!!
        assertNull(option.status, "UNSPECIFIED status must not be translated into a filter")
        assertNull(option.paymentId, "empty paymentId must not be translated into a filter")
    }

    // ─── status enum round-trip ─────────────────────────────────────────────

    @Test
    fun `refund status enum round-trips domain to proto and back through the server`() {
        val cases = listOf(
            Triple(RefundStatus.PENDING, ProtoRefundStatus.REFUND_STATUS_PENDING, "pending"),
            Triple(RefundStatus.SUCCEEDED, ProtoRefundStatus.REFUND_STATUS_SUCCEEDED, "succeeded"),
            Triple(RefundStatus.FAILED, ProtoRefundStatus.REFUND_STATUS_FAILED, "failed"),
            Triple(RefundStatus.REVIEW, ProtoRefundStatus.REFUND_STATUS_REVIEW, "review"),
        )

        cases.forEach { (domain, proto, value) ->
            refundUseCase.findRefundResult = PaymentTestFixtures.refundAggregate(status = domain)
            val findObserver = CapturingObserver<FindRefundResponse>()
            withActor(ADMIN_UID, admin = true) {
                server.findRefund(FindRefundRequest.newBuilder().setUid(PaymentTestFixtures.REFUND_UID.toString()).build(), findObserver)
            }
            assertEquals(proto, findObserver.values.single().refund.status, "domain $domain -> proto")

            refundUseCase.findAllResult = PageImpl(emptyList())
            val listObserver = CapturingObserver<ListRefundsResponse>()
            withActor(ADMIN_UID, admin = true) {
                server.listRefunds(ListRefundsRequest.newBuilder().setStatus(proto).build(), listObserver)
            }
            assertEquals(value, refundUseCase.capturedOption?.status, "proto $proto -> domain value")
        }
    }
}
