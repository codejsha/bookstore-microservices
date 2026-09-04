package com.codejsha.bookstore.payment.infrastructure.adapter.protosvc

import com.codejsha.bookstore.generated.application.port.pb.paymentpb.FindPaymentRequest
import com.codejsha.bookstore.generated.application.port.pb.paymentpb.FindPaymentResponse
import com.codejsha.bookstore.generated.application.port.pb.paymentpb.ListPaymentsRequest
import com.codejsha.bookstore.generated.application.port.pb.paymentpb.ListPaymentsResponse
import com.codejsha.bookstore.payment.application.usecase.PaymentUseCase
import com.codejsha.bookstore.payment.application.usecase.RefundUseCase
import com.codejsha.bookstore.payment.domain.aggregate.PaymentAggregate
import com.codejsha.bookstore.payment.domain.aggregate.PaymentAttemptEntity
import com.codejsha.bookstore.payment.domain.aggregate.RefundAggregate
import com.codejsha.bookstore.payment.domain.constant.PaymentStatus
import com.codejsha.bookstore.payment.domain.model.command.PaymentCreateCommand
import com.codejsha.bookstore.payment.domain.model.command.PaymentUpdateCommand
import com.codejsha.bookstore.payment.domain.model.option.PaymentQueryOption
import com.codejsha.bookstore.payment.domain.model.option.RefundQueryOption
import com.codejsha.bookstore.payment.infrastructure.support.auth.ActorIdentity
import com.codejsha.bookstore.payment.support.PaymentTestFixtures
import com.codejsha.platform.shared.data.ActorContext
import io.grpc.Context
import io.grpc.Status
import io.grpc.stub.StreamObserver
import org.junit.jupiter.api.Test
import org.springframework.data.domain.Page
import org.springframework.data.domain.PageImpl
import org.springframework.data.domain.Pageable
import java.util.UUID
import kotlin.test.assertEquals
import com.codejsha.bookstore.generated.application.port.pb.paymentpb.PaymentStatus as ProtoPaymentStatus

class PaymentGrpcServerPaymentTest {

    private companion object {
        val ACTOR_UID: UUID = UUID.fromString("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa")
        val OTHER_UID: UUID = UUID.fromString("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb")
        val ADMIN_UID: UUID = UUID.fromString("cccccccc-cccc-cccc-cccc-cccccccccccc")
    }

    private class FakePaymentUseCase(
        var findResult: PaymentAggregate? = null,
        var findError: Throwable? = null,
        var listError: Throwable? = null,
    ) : PaymentUseCase {
        var capturedListOption: PaymentQueryOption? = null

        override suspend fun findAllPayments(option: PaymentQueryOption, pageable: Pageable, context: ActorContext): Page<PaymentAggregate> {
            capturedListOption = option
            listError?.let { throw it }
            return PageImpl(listOfNotNull(findResult))
        }

        override suspend fun findPayment(uid: UUID, context: ActorContext): PaymentAggregate {
            findError?.let { throw it }
            return findResult ?: error("findResult not seeded")
        }

        override suspend fun findPaymentByPaymentId(paymentId: String, context: ActorContext): PaymentAggregate? = error("not used")
        override suspend fun createPayment(command: PaymentCreateCommand, context: ActorContext): PaymentAggregate = error("not used")
        override suspend fun updatePayment(uid: UUID, command: PaymentUpdateCommand, context: ActorContext): PaymentAggregate = error("not used")
        override suspend fun findAllPaymentAttempts(paymentUid: UUID, pageable: Pageable, context: ActorContext): Page<PaymentAttemptEntity> = error("not used")
        override suspend fun findPaymentAttempt(paymentUid: UUID, uid: UUID, context: ActorContext): PaymentAttemptEntity = error("not used")
    }

    private class UnusedRefundUseCase : RefundUseCase {
        override suspend fun findAllRefunds(option: RefundQueryOption, pageable: Pageable, context: ActorContext): Page<RefundAggregate> = error("not used")
        override suspend fun findRefund(uid: UUID, context: ActorContext): RefundAggregate = error("not used")
        override suspend fun createRefund(command: com.codejsha.bookstore.payment.domain.model.command.RefundCreateCommand, context: ActorContext): RefundAggregate = error("not used")
    }

    private class CapturingObserver<T> : StreamObserver<T> {
        val values = mutableListOf<T>()
        var completed = false
        var error: Throwable? = null
        override fun onNext(value: T) { values.add(value) }
        override fun onError(t: Throwable) { error = t }
        override fun onCompleted() { completed = true }
    }

    private inline fun withActor(uid: UUID, admin: Boolean, block: () -> Unit) {
        val ctx = Context.current().withValue(ActorIdentity.CONTEXT_KEY, ActorIdentity(uid, admin))
        val previous = ctx.attach()
        try {
            block()
        } finally {
            ctx.detach(previous)
        }
    }

    // ─── actor guard ──────────────────────────────────────────────────────────

    @Test
    fun `findPayment_whenActorIdentityMissing_returnsUnauthenticated`() {
        val server = PaymentGrpcServer(FakePaymentUseCase(findResult = PaymentTestFixtures.paymentAggregate()), UnusedRefundUseCase())

        val observer = CapturingObserver<FindPaymentResponse>()
        server.findPayment(FindPaymentRequest.newBuilder().setUid(PaymentTestFixtures.PAYMENT_UID.toString()).build(), observer)

        val status = Status.fromThrowable(observer.error!!)
        assertEquals(Status.Code.UNAUTHENTICATED, status.code)
        assertEquals("missing or invalid actor identity", status.description)
    }

    // ─── findPayment ownership ──────────────────────────────────────────────────

    @Test
    fun `findPayment_whenPaymentOwnedByAnotherCustomerAndActorNotAdmin_returnsNotFound`() {
        val uid = PaymentTestFixtures.PAYMENT_UID
        val useCase = FakePaymentUseCase(findResult = PaymentTestFixtures.paymentAggregate(uid = uid, customerId = OTHER_UID.toString()))
        val server = PaymentGrpcServer(useCase, UnusedRefundUseCase())

        val observer = CapturingObserver<FindPaymentResponse>()
        withActor(ACTOR_UID, admin = false) {
            server.findPayment(FindPaymentRequest.newBuilder().setUid(uid.toString()).build(), observer)
        }

        val status = Status.fromThrowable(observer.error!!)
        assertEquals(Status.Code.NOT_FOUND, status.code)
        assertEquals("Payment with uid $uid not found", status.description)
    }

    @Test
    fun `findPayment_whenActorOwnsPayment_returnsPayment`() {
        val useCase = FakePaymentUseCase(findResult = PaymentTestFixtures.paymentAggregate(customerId = ACTOR_UID.toString()))
        val server = PaymentGrpcServer(useCase, UnusedRefundUseCase())

        val observer = CapturingObserver<FindPaymentResponse>()
        withActor(ACTOR_UID, admin = false) {
            server.findPayment(FindPaymentRequest.newBuilder().setUid(PaymentTestFixtures.PAYMENT_UID.toString()).build(), observer)
        }

        assertEquals(1, observer.values.size)
    }

    @Test
    fun `findPayment_whenActorIsAdmin_returnsAnotherCustomersPayment`() {
        val useCase = FakePaymentUseCase(findResult = PaymentTestFixtures.paymentAggregate(customerId = OTHER_UID.toString()))
        val server = PaymentGrpcServer(useCase, UnusedRefundUseCase())

        val observer = CapturingObserver<FindPaymentResponse>()
        withActor(ADMIN_UID, admin = true) {
            server.findPayment(FindPaymentRequest.newBuilder().setUid(PaymentTestFixtures.PAYMENT_UID.toString()).build(), observer)
        }

        assertEquals(1, observer.values.size)
        assertEquals(OTHER_UID.toString(), observer.values.single().payment.customerUid)
    }

    // ─── findPayment error mapping ──────────────────────────────────────────────

    @Test
    fun `findPayment_whenUsecaseThrowsNoSuchElement_returnsNotFound`() {
        val useCase = FakePaymentUseCase(findError = NoSuchElementException("Payment with uid x not found"))
        val server = PaymentGrpcServer(useCase, UnusedRefundUseCase())

        val observer = CapturingObserver<FindPaymentResponse>()
        withActor(ADMIN_UID, admin = true) {
            server.findPayment(FindPaymentRequest.newBuilder().setUid(UUID.randomUUID().toString()).build(), observer)
        }

        assertEquals(Status.Code.NOT_FOUND, Status.fromThrowable(observer.error!!).code)
    }

    @Test
    fun `findPayment_whenUidMalformed_returnsInvalidArgument`() {
        val server = PaymentGrpcServer(FakePaymentUseCase(), UnusedRefundUseCase())

        val observer = CapturingObserver<FindPaymentResponse>()
        withActor(ACTOR_UID, admin = false) {
            server.findPayment(FindPaymentRequest.newBuilder().setUid("not-a-uuid").build(), observer)
        }

        val status = Status.fromThrowable(observer.error!!)
        assertEquals(Status.Code.INVALID_ARGUMENT, status.code)
        assertEquals("invalid argument", status.description)
    }

    // ─── listPayments scoping ───────────────────────────────────────────────────

    @Test
    fun `listPayments_whenActorNotAdmin_pinsFilterToActorIgnoringRequestField`() {
        val useCase = FakePaymentUseCase()
        val server = PaymentGrpcServer(useCase, UnusedRefundUseCase())

        val observer = CapturingObserver<ListPaymentsResponse>()
        withActor(ACTOR_UID, admin = false) {
            server.listPayments(
                ListPaymentsRequest.newBuilder().setCustomerUid(OTHER_UID.toString()).build(),
                observer,
            )
        }

        assertEquals(ACTOR_UID.toString(), useCase.capturedListOption?.customerId)
    }

    @Test
    fun `listPayments_whenActorIsAdmin_honorsRequestedCustomerUid`() {
        val useCase = FakePaymentUseCase()
        val server = PaymentGrpcServer(useCase, UnusedRefundUseCase())

        val observer = CapturingObserver<ListPaymentsResponse>()
        withActor(ADMIN_UID, admin = true) {
            server.listPayments(
                ListPaymentsRequest.newBuilder().setCustomerUid(OTHER_UID.toString()).build(),
                observer,
            )
        }

        assertEquals(OTHER_UID.toString(), useCase.capturedListOption?.customerId)
    }

    @Test
    fun `listPayments_whenActorIsAdminAndCustomerUidEmpty_returnsEveryCustomersPayments`() {
        val useCase = FakePaymentUseCase()
        val server = PaymentGrpcServer(useCase, UnusedRefundUseCase())

        val observer = CapturingObserver<ListPaymentsResponse>()
        withActor(ADMIN_UID, admin = true) {
            server.listPayments(ListPaymentsRequest.newBuilder().build(), observer)
        }

        assertEquals(null, useCase.capturedListOption?.customerId)
    }

    @Test
    fun `listPayments_whenUsecaseFailsUnexpectedly_returnsGenericInternal`() {
        val useCase = FakePaymentUseCase(listError = IllegalStateException("db down: secret host"))
        val server = PaymentGrpcServer(useCase, UnusedRefundUseCase())

        val observer = CapturingObserver<ListPaymentsResponse>()
        withActor(ADMIN_UID, admin = true) {
            server.listPayments(
                ListPaymentsRequest.newBuilder().setCustomerUid(UUID.randomUUID().toString()).build(),
                observer,
            )
        }

        val status = Status.fromThrowable(observer.error!!)
        assertEquals(Status.Code.INTERNAL, status.code)
        assertEquals("internal error", status.description)
    }

    // ─── enum bridge ────────────────────────────────────────────────────────────

    @Test
    fun `findPayment_whenStatusPartiallyCaptured_mapsOntoDedicatedProtoStatus`() {
        val useCase = FakePaymentUseCase(
            findResult = PaymentTestFixtures.paymentAggregate(status = PaymentStatus.PARTIALLY_CAPTURED),
        )
        val server = PaymentGrpcServer(useCase, UnusedRefundUseCase())

        val observer = CapturingObserver<FindPaymentResponse>()
        withActor(ADMIN_UID, admin = true) {
            server.findPayment(FindPaymentRequest.newBuilder().setUid(PaymentTestFixtures.PAYMENT_UID.toString()).build(), observer)
        }

        assertEquals(ProtoPaymentStatus.PAYMENT_STATUS_PARTIALLY_CAPTURED, observer.values.single().payment.status)
    }
}
