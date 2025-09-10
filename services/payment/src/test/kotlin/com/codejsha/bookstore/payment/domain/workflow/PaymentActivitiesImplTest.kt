package com.codejsha.bookstore.payment.domain.workflow

import com.codejsha.bookstore.payment.application.port.hyperswitch.HyperswitchClient
import com.codejsha.bookstore.payment.application.port.hyperswitch.HyperswitchClientException
import com.codejsha.bookstore.payment.application.port.repo.MandateRepo
import com.codejsha.bookstore.payment.application.port.repo.MandateResult
import com.codejsha.bookstore.payment.application.port.repo.PaymentAttemptRepo
import com.codejsha.bookstore.payment.application.port.repo.PaymentAttemptResult
import com.codejsha.bookstore.payment.application.port.repo.PaymentRepo
import com.codejsha.bookstore.payment.application.port.repo.PaymentResult
import com.codejsha.bookstore.payment.application.port.repo.RefundRepo
import com.codejsha.bookstore.payment.application.port.repo.RefundResult
import com.codejsha.bookstore.payment.domain.model.command.MandateCreateCommand
import com.codejsha.bookstore.payment.domain.model.command.PaymentAttemptCreateCommand
import com.codejsha.bookstore.payment.domain.model.command.PaymentCreateCommand
import com.codejsha.bookstore.payment.domain.model.command.PaymentUpdateCommand
import com.codejsha.bookstore.payment.domain.model.command.RefundCreateCommand
import com.codejsha.bookstore.payment.domain.model.external.HyperswitchPaymentCommand
import com.codejsha.bookstore.payment.domain.model.external.HyperswitchPaymentLookup
import com.codejsha.bookstore.payment.domain.model.external.HyperswitchPaymentResult
import com.codejsha.bookstore.payment.domain.model.external.HyperswitchRefundCommand
import com.codejsha.bookstore.payment.domain.model.external.HyperswitchRefundResult
import com.codejsha.bookstore.payment.domain.model.external.HyperswitchSetupMandateCommand
import com.codejsha.bookstore.payment.domain.model.external.HyperswitchSetupMandateResult
import com.codejsha.bookstore.payment.domain.model.external.HyperswitchWebhookEvent
import com.codejsha.bookstore.payment.domain.model.option.MandateQueryOption
import com.codejsha.bookstore.payment.domain.model.option.PaymentQueryOption
import com.codejsha.bookstore.payment.domain.model.option.RefundQueryOption
import com.codejsha.bookstore.payment.support.FakeDistributedLock
import com.codejsha.bookstore.payment.support.PaymentTestFixtures
import com.codejsha.platform.shared.data.ActorContext
import io.temporal.failure.ApplicationFailure
import org.junit.jupiter.api.Test
import org.springframework.data.domain.Page
import org.springframework.data.domain.PageImpl
import org.springframework.data.domain.Pageable
import java.math.BigDecimal
import java.util.UUID
import kotlin.test.assertEquals
import kotlin.test.assertFailsWith
import kotlin.test.assertFalse
import kotlin.test.assertNull
import kotlin.test.assertTrue

class PaymentActivitiesImplTest {

    private val orderUid = "11111111-1111-1111-1111-111111111111"

    // ─── processPayment ───────────────────────────────────────────────────────

    @Test
    fun `processPayment charges via gateway and persists the real status and gateway id`() {
        val client = StubHyperswitchClient(
            paymentResult = HyperswitchPaymentResult(
                gatewayPaymentId = "pay_hyper_123",
                status = "succeeded",
                connector = "stripe",
                amountCapturable = 0,
                amountReceived = 4_999,
                errorCode = null,
                errorMessage = null,
            ),
        )
        val paymentRepo = FakePaymentRepo()
        val attemptRepo = FakePaymentAttemptRepo()
        val lock = FakeDistributedLock()
        val activities = PaymentActivitiesImpl(paymentRepo, attemptRepo, FakeRefundRepo(), FakeMandateRepo(), client, lock)

        val result = activities.processPayment(
            ProcessPaymentRequest(orderUid = orderUid, userUid = "00000000-0000-0000-0000-000000000007", amount = BigDecimal("49.99"), currency = "USD"),
        )

        assertEquals("succeeded", result.status)
        assertEquals("pay_hyper_123", result.gatewayPaymentId)
        assertEquals(1, lock.invocationCount)
        assertEquals("payment:$orderUid", lock.lastKey)
        assertNull(lock.lastTtl, "activity relies on watchdog auto-renewal, not a fixed lease")

        val gwCmd = client.authorizeCommand!!
        assertEquals(4_999L, gwCmd.amountMinor)
        assertEquals(orderUid, gwCmd.idempotencyKey)
        assertEquals("USD", gwCmd.currency)

        val created = paymentRepo.createCommand!!
        assertEquals("succeeded", created.status)
        assertEquals("pay_hyper_123", created.paymentId)
        assertEquals(orderUid, created.idempotencyKey)
        assertEquals("stripe", created.connector)
        assertEquals(4_999L, created.amount)

        val attempt = attemptRepo.createCommand!!
        assertEquals("succeeded", attempt.status)
        assertEquals("pay_hyper_123", attempt.paymentId)
    }

    @Test
    fun `processPayment short-circuits a settled prior payment without re-charging`() {
        val client = StubHyperswitchClient()
        val paymentRepo = FakePaymentRepo(
            existingByKey = PaymentTestFixtures.paymentResult(status = "succeeded", paymentId = "pay_prev"),
        )
        val attemptRepo = FakePaymentAttemptRepo()
        val lock = FakeDistributedLock()
        val activities = PaymentActivitiesImpl(paymentRepo, attemptRepo, FakeRefundRepo(), FakeMandateRepo(), client, lock)

        val result = activities.processPayment(
            ProcessPaymentRequest(orderUid = orderUid, userUid = "00000000-0000-0000-0000-000000000007", amount = BigDecimal("10.00"), currency = "USD"),
        )

        assertEquals("succeeded", result.status)
        assertEquals("pay_prev", result.gatewayPaymentId)
        assertEquals(0, lock.invocationCount, "settled fast-path must not take the lock")
        assertNull(client.authorizeCommand, "gateway must not be called for a settled payment")
        assertNull(paymentRepo.createCommand, "no new payment row for a settled payment")
        assertNull(attemptRepo.createCommand)
    }

    @Test
    fun `processPayment re-attempts the gateway when a prior row is failed (non-settled)`() {
        val client = StubHyperswitchClient(
            paymentResult = HyperswitchPaymentResult("pay_retry", "succeeded", "stripe", 0, 1_000, null, null),
        )
        val paymentRepo = FakePaymentRepo(
            existingByKey = PaymentTestFixtures.paymentResult(status = "failed", paymentId = "pay_failed"),
        )
        val activities = PaymentActivitiesImpl(paymentRepo, FakePaymentAttemptRepo(), FakeRefundRepo(), FakeMandateRepo(), client, FakeDistributedLock())

        val result = activities.processPayment(
            ProcessPaymentRequest(orderUid = orderUid, userUid = "00000000-0000-0000-0000-000000000007", amount = BigDecimal("10.00"), currency = "USD"),
        )

        assertTrue(client.authorizeCommand != null, "a failed prior row must not short-circuit the retry")
        assertEquals("succeeded", result.status)
    }

    @Test
    fun `processPayment on gateway error records a failed attempt and rethrows`() {
        val client = StubHyperswitchClient(authorizeError = HyperswitchClientException("boom", errorCode = "CE_00"))
        val paymentRepo = FakePaymentRepo()
        val attemptRepo = FakePaymentAttemptRepo()
        val activities = PaymentActivitiesImpl(paymentRepo, attemptRepo, FakeRefundRepo(), FakeMandateRepo(), client, FakeDistributedLock())

        assertFailsWith<HyperswitchClientException> {
            activities.processPayment(
                ProcessPaymentRequest(orderUid = orderUid, userUid = "00000000-0000-0000-0000-000000000007", amount = BigDecimal("10.00"), currency = "USD"),
            )
        }

        val failed = paymentRepo.createCommand!!
        assertEquals("failed", failed.status)
        assertNull(failed.idempotencyKey, "failed reconciliation row must not carry the saga key")
        assertEquals("CE_00", failed.errorCode)
        assertEquals("failed", attemptRepo.createCommand!!.status)
    }

    @Test
    fun `processPayment retries while the gateway still reports the fresh charge pending`() {
        val client = StubHyperswitchClient(
            paymentResult = HyperswitchPaymentResult("pay_pending", "processing", "stripe", 1_000, 0, null, null),
        )
        val activities = PaymentActivitiesImpl(FakePaymentRepo(), FakePaymentAttemptRepo(), FakeRefundRepo(), FakeMandateRepo(), client, FakeDistributedLock())

        val failure = assertFailsWith<ApplicationFailure> {
            activities.processPayment(
                ProcessPaymentRequest(orderUid = orderUid, userUid = "00000000-0000-0000-0000-000000000007", amount = BigDecimal("10.00"), currency = "USD"),
            )
        }

        assertEquals("GatewayPaymentPending", failure.type)
        assertFalse(failure.isNonRetryable)
    }

    @Test
    fun `processPayment resolves a pending prior row from the gateway instead of re-charging`() {
        val client = StubHyperswitchClient(
            lookupById = HyperswitchPaymentLookup("pay_pending", "succeeded", "USD", 1_000, 0, 1_000, "stripe", null, null),
        )
        val paymentRepo = FakePaymentRepo(
            existingByKey = PaymentTestFixtures.paymentResult(status = "processing", paymentId = "pay_pending"),
        )
        val activities = PaymentActivitiesImpl(paymentRepo, FakePaymentAttemptRepo(), FakeRefundRepo(), FakeMandateRepo(), client, FakeDistributedLock())

        val result = activities.processPayment(
            ProcessPaymentRequest(orderUid = orderUid, userUid = "00000000-0000-0000-0000-000000000007", amount = BigDecimal("10.00"), currency = "USD"),
        )

        assertNull(client.authorizeCommand, "a pending prior row must be resolved, not re-charged")
        assertEquals("pay_pending", client.lookedUpPaymentId)
        assertEquals("succeeded", paymentRepo.syncedStatus)
        assertEquals("succeeded", result.status)
    }

    @Test
    fun `processPayment recovers a gateway-held payment when authorization errors`() {
        val client = StubHyperswitchClient(
            authorizeError = HyperswitchClientException("duplicate payment", errorCode = "HE_01"),
            lookupByKey = HyperswitchPaymentLookup("pay_held", "succeeded", "USD", 1_000, 0, 1_000, "stripe", null, null),
        )
        val paymentRepo = FakePaymentRepo()
        val attemptRepo = FakePaymentAttemptRepo()
        val activities = PaymentActivitiesImpl(paymentRepo, attemptRepo, FakeRefundRepo(), FakeMandateRepo(), client, FakeDistributedLock())

        val result = activities.processPayment(
            ProcessPaymentRequest(orderUid = orderUid, userUid = "00000000-0000-0000-0000-000000000007", amount = BigDecimal("10.00"), currency = "USD"),
        )

        assertEquals(orderUid, client.lookedUpKey)
        assertEquals("succeeded", result.status)
        assertEquals("pay_held", result.gatewayPaymentId)
        assertEquals(orderUid, paymentRepo.createCommand!!.idempotencyKey)
        assertEquals("succeeded", attemptRepo.createCommand!!.status)
    }

    @Test
    fun `processPayment converts the amount using the currency's ISO 4217 minor unit`() {
        val client = StubHyperswitchClient()
        val activities = PaymentActivitiesImpl(FakePaymentRepo(), FakePaymentAttemptRepo(), FakeRefundRepo(), FakeMandateRepo(), client, FakeDistributedLock())

        activities.processPayment(
            ProcessPaymentRequest(orderUid = orderUid, userUid = "00000000-0000-0000-0000-000000000007", amount = BigDecimal("15000"), currency = "KRW"),
        )

        assertEquals(15_000L, client.authorizeCommand!!.amountMinor, "KRW has no minor unit and must not be scaled by 100")
    }

    @Test
    fun `processPayment fails non-retryably for an unknown currency code`() {
        val client = StubHyperswitchClient()
        val activities = PaymentActivitiesImpl(FakePaymentRepo(), FakePaymentAttemptRepo(), FakeRefundRepo(), FakeMandateRepo(), client, FakeDistributedLock())

        val failure = assertFailsWith<ApplicationFailure> {
            activities.processPayment(
                ProcessPaymentRequest(orderUid = orderUid, userUid = "00000000-0000-0000-0000-000000000007", amount = BigDecimal("10.00"), currency = "XXX"),
            )
        }

        assertTrue(failure.isNonRetryable)
        assertEquals("UnsupportedCurrency", failure.type)
        assertNull(client.authorizeCommand)
    }

    @Test
    fun `refundPayment treats an unrecognized payment status as non-refundable instead of crashing`() {
        val paymentUid = UUID.randomUUID()
        val client = StubHyperswitchClient()
        val paymentRepo = FakePaymentRepo(findOne = PaymentTestFixtures.paymentResult(status = "partially_captured_and_capturable"))
        val activities = PaymentActivitiesImpl(paymentRepo, FakePaymentAttemptRepo(), FakeRefundRepo(), FakeMandateRepo(), client, FakeDistributedLock())

        assertFalse(activities.refundPayment(paymentUid.toString()))
        assertNull(client.refundCommand)
    }

    // ─── refundPayment ─────────────────────────────────────────────────────────

    @Test
    fun `refundPayment refunds via gateway and persists the real refund status and id`() {
        val paymentUid = UUID.randomUUID()
        val client = StubHyperswitchClient(
            refundResult = HyperswitchRefundResult("ref_hyper_1", "succeeded", "stripe", null, null),
        )
        val paymentRepo = FakePaymentRepo(
            findOne = PaymentTestFixtures.paymentResult(status = "succeeded", paymentId = "pay_ok"),
        )
        val refundRepo = FakeRefundRepo()
        val lock = FakeDistributedLock()
        val activities = PaymentActivitiesImpl(paymentRepo, FakePaymentAttemptRepo(), refundRepo, FakeMandateRepo(), client, lock)

        val refunded = activities.refundPayment(paymentUid.toString())

        assertTrue(refunded)
        assertEquals(1, lock.invocationCount)
        assertEquals("refund:$paymentUid", lock.lastKey)
        val gwCmd = client.refundCommand!!
        assertEquals("pay_ok", gwCmd.gatewayPaymentId)
        assertEquals("refund:$paymentUid", gwCmd.idempotencyKey)

        val created = refundRepo.createCommand!!
        assertEquals("ref_hyper_1", created.refundId)
        assertEquals("succeeded", created.status)
        assertEquals("refund:$paymentUid", created.idempotencyKey)
    }

    @Test
    fun `refundPayment is idempotent when a refund already exists for the key`() {
        val paymentUid = UUID.randomUUID()
        val client = StubHyperswitchClient()
        val refundRepo = FakeRefundRepo(existingByKey = PaymentTestFixtures.refundResult())
        val lock = FakeDistributedLock()
        val activities = PaymentActivitiesImpl(FakePaymentRepo(), FakePaymentAttemptRepo(), refundRepo, FakeMandateRepo(), client, lock)

        val refunded = activities.refundPayment(paymentUid.toString())

        assertTrue(refunded)
        assertEquals(0, lock.invocationCount, "existing refund fast-path must not take the lock")
        assertNull(client.refundCommand, "gateway must not be called when the refund already exists")
        assertNull(refundRepo.createCommand)
    }

    @Test
    fun `refundPayment reports false when the gateway refund fails`() {
        val paymentUid = UUID.randomUUID()
        val client = StubHyperswitchClient(
            refundResult = HyperswitchRefundResult("ref_failed", "failed", "stripe", "RE_01", "insufficient funds"),
        )
        val paymentRepo = FakePaymentRepo(findOne = PaymentTestFixtures.paymentResult(status = "succeeded"))
        val refundRepo = FakeRefundRepo()
        val activities = PaymentActivitiesImpl(paymentRepo, FakePaymentAttemptRepo(), refundRepo, FakeMandateRepo(), client, FakeDistributedLock())

        val refunded = activities.refundPayment(paymentUid.toString())

        assertFalse(refunded)
        assertEquals("failed", refundRepo.createCommand!!.status, "the failed refund must still be persisted for reconciliation")
    }

    @Test
    fun `refundPayment reports false when the existing refund for the key had failed`() {
        val paymentUid = UUID.randomUUID()
        val client = StubHyperswitchClient()
        val refundRepo = FakeRefundRepo(existingByKey = PaymentTestFixtures.refundResult(status = "failed"))
        val activities = PaymentActivitiesImpl(FakePaymentRepo(), FakePaymentAttemptRepo(), refundRepo, FakeMandateRepo(), client, FakeDistributedLock())

        assertFalse(activities.refundPayment(paymentUid.toString()))
        assertNull(client.refundCommand)
    }

    @Test
    fun `refundPayment refunds the captured amount, not the authorized amount, for a partial capture`() {
        val paymentUid = UUID.randomUUID()
        val client = StubHyperswitchClient(
            refundResult = HyperswitchRefundResult("ref_hyper_2", "succeeded", "stripe", null, null),
        )
        val paymentRepo = FakePaymentRepo(
            findOne = PaymentTestFixtures.paymentResult(status = "partially_captured", amount = 4999L)
                .copy(amountCaptured = 2000L),
        )
        val refundRepo = FakeRefundRepo()
        val activities = PaymentActivitiesImpl(paymentRepo, FakePaymentAttemptRepo(), refundRepo, FakeMandateRepo(), client, FakeDistributedLock())

        activities.refundPayment(paymentUid.toString())

        assertEquals(2000L, client.refundCommand!!.amountMinor)
        assertEquals(2000L, refundRepo.createCommand!!.amount)
    }

    @Test
    fun `refundPayment skips a non-refundable payment without calling the gateway`() {
        val paymentUid = UUID.randomUUID()
        val client = StubHyperswitchClient()
        val paymentRepo = FakePaymentRepo(findOne = PaymentTestFixtures.paymentResult(status = "requires_payment_method"))
        val refundRepo = FakeRefundRepo()
        val activities = PaymentActivitiesImpl(paymentRepo, FakePaymentAttemptRepo(), refundRepo, FakeMandateRepo(), client, FakeDistributedLock())

        val refunded = activities.refundPayment(paymentUid.toString())

        assertFalse(refunded)
        assertNull(client.refundCommand, "gateway must not be called for a non-refundable payment")
        assertNull(refundRepo.createCommand)
    }

    // ─── refundPaymentByOrder ──────────────────────────────────────────────────

    @Test
    fun `refundPaymentByOrder refunds the payment recorded under the order's idempotency key`() {
        val client = StubHyperswitchClient(
            refundResult = HyperswitchRefundResult("ref_hyper_3", "succeeded", "stripe", null, null),
        )
        val payment = PaymentTestFixtures.paymentResult(status = "succeeded", paymentId = "pay_by_order")
        val paymentRepo = FakePaymentRepo(existingByKey = payment)
        val refundRepo = FakeRefundRepo()
        val activities = PaymentActivitiesImpl(paymentRepo, FakePaymentAttemptRepo(), refundRepo, FakeMandateRepo(), client, FakeDistributedLock())

        val refunded = activities.refundPaymentByOrder(orderUid)

        assertTrue(refunded)
        val gwCmd = client.refundCommand!!
        assertEquals("pay_by_order", gwCmd.gatewayPaymentId)
        assertEquals("refund:${payment.uid}", gwCmd.idempotencyKey)
        assertEquals("refund:${payment.uid}", refundRepo.createCommand!!.idempotencyKey)
    }

    @Test
    fun `refundPaymentByOrder returns false when no payment exists for the order`() {
        val client = StubHyperswitchClient()
        val refundRepo = FakeRefundRepo()
        val activities = PaymentActivitiesImpl(FakePaymentRepo(), FakePaymentAttemptRepo(), refundRepo, FakeMandateRepo(), client, FakeDistributedLock())

        val refunded = activities.refundPaymentByOrder(orderUid)

        assertFalse(refunded)
        assertNull(client.refundCommand)
        assertNull(refundRepo.createCommand)
    }

    @Test
    fun `refundPaymentByOrder is idempotent when the payment's refund already exists`() {
        val client = StubHyperswitchClient()
        val payment = PaymentTestFixtures.paymentResult(status = "succeeded")
        val paymentRepo = FakePaymentRepo(existingByKey = payment)
        val refundRepo = FakeRefundRepo(existingByKey = PaymentTestFixtures.refundResult())
        val activities = PaymentActivitiesImpl(paymentRepo, FakePaymentAttemptRepo(), refundRepo, FakeMandateRepo(), client, FakeDistributedLock())

        val refunded = activities.refundPaymentByOrder(orderUid)

        assertTrue(refunded)
        assertNull(client.refundCommand, "gateway must not be called when the refund already exists")
        assertNull(refundRepo.createCommand)
    }

    @Test
    fun `refundPaymentByOrder restores a gateway-only payment locally and refunds it`() {
        val client = StubHyperswitchClient(
            lookupByKey = HyperswitchPaymentLookup(
                gatewayPaymentId = "pay_gw_only",
                status = "succeeded",
                currency = "USD",
                amountMinor = 4_999,
                amountCapturable = 0,
                amountReceived = 4_999,
                connector = "stripe",
                errorCode = null,
                errorMessage = null,
            ),
            refundResult = HyperswitchRefundResult("ref_heal", "succeeded", "stripe", null, null),
        )
        val paymentRepo = FakePaymentRepo()
        val refundRepo = FakeRefundRepo()
        val activities = PaymentActivitiesImpl(paymentRepo, FakePaymentAttemptRepo(), refundRepo, FakeMandateRepo(), client, FakeDistributedLock())

        val refunded = activities.refundPaymentByOrder(orderUid)

        assertTrue(refunded)
        assertEquals(orderUid, client.lookedUpKey)
        val healed = paymentRepo.createCommand!!
        assertEquals(orderUid, healed.idempotencyKey, "restored row must carry the saga key so the heal is idempotent")
        assertEquals("pay_gw_only", healed.paymentId)
        assertEquals("succeeded", healed.status)
        assertEquals(4_999L, healed.amountCaptured)
        assertEquals("pay_gw_only", client.refundCommand!!.gatewayPaymentId)
        assertEquals("ref_heal", refundRepo.createCommand!!.refundId)
    }

    @Test
    fun `refundPayment retries while the gateway still reports the payment pending`() {
        val paymentUid = UUID.randomUUID()
        val client = StubHyperswitchClient(
            lookupById = HyperswitchPaymentLookup(
                "pay_pending", "processing", "USD", 1_000, 1_000, null, "stripe", null, null,
            ),
        )
        val paymentRepo = FakePaymentRepo(
            findOne = PaymentTestFixtures.paymentResult(status = "processing", paymentId = "pay_pending"),
        )
        val refundRepo = FakeRefundRepo()
        val activities = PaymentActivitiesImpl(paymentRepo, FakePaymentAttemptRepo(), refundRepo, FakeMandateRepo(), client, FakeDistributedLock())

        val failure = assertFailsWith<ApplicationFailure> {
            activities.refundPayment(paymentUid.toString())
        }

        assertEquals("GatewayPaymentPending", failure.type)
        assertFalse(failure.isNonRetryable, "a pending gateway payment must stay retryable")
        assertEquals("pay_pending", client.lookedUpPaymentId)
        assertNull(client.refundCommand)
        assertNull(refundRepo.createCommand)
    }

    @Test
    fun `refundPayment refunds once the gateway reports a pending payment settled`() {
        val paymentUid = UUID.randomUUID()
        val client = StubHyperswitchClient(
            lookupById = HyperswitchPaymentLookup(
                "pay_late", "succeeded", "USD", 2_000, 0, 2_000, "stripe", null, null,
            ),
            refundResult = HyperswitchRefundResult("ref_late", "succeeded", "stripe", null, null),
        )
        val paymentRepo = FakePaymentRepo(
            findOne = PaymentTestFixtures.paymentResult(status = "processing", paymentId = "pay_late"),
        )
        val refundRepo = FakeRefundRepo()
        val activities = PaymentActivitiesImpl(paymentRepo, FakePaymentAttemptRepo(), refundRepo, FakeMandateRepo(), client, FakeDistributedLock())

        val refunded = activities.refundPayment(paymentUid.toString())

        assertTrue(refunded)
        assertEquals("succeeded", paymentRepo.syncedStatus, "the settled status must be persisted locally")
        assertEquals(2_000L, client.refundCommand!!.amountMinor, "refund must use the live captured amount")
    }

    @Test
    fun `refundPaymentByOrder returns false for a non-refundable payment`() {
        val client = StubHyperswitchClient()
        val payment = PaymentTestFixtures.paymentResult(status = "requires_payment_method")
        val paymentRepo = FakePaymentRepo(existingByKey = payment)
        val refundRepo = FakeRefundRepo()
        val activities = PaymentActivitiesImpl(paymentRepo, FakePaymentAttemptRepo(), refundRepo, FakeMandateRepo(), client, FakeDistributedLock())

        val refunded = activities.refundPaymentByOrder(orderUid)

        assertFalse(refunded)
        assertNull(client.refundCommand)
        assertNull(refundRepo.createCommand)
    }

    // ─── processPayment: the mandate is the instrument ────────────────────────

    @Test
    fun `processPayment charges against the customer's active mandate off-session`() {
        val client = StubHyperswitchClient()
        val userUid = "00000000-0000-0000-0000-000000000007"
        val activities = PaymentActivitiesImpl(
            FakePaymentRepo(),
            FakePaymentAttemptRepo(),
            FakeRefundRepo(),
            FakeMandateRepo(PaymentTestFixtures.mandateResult(mandateId = "mand_active", customerId = userUid)),
            client,
            FakeDistributedLock(),
        )

        activities.processPayment(
            ProcessPaymentRequest(orderUid = orderUid, userUid = userUid, amount = BigDecimal("49.99"), currency = "USD"),
        )

        val gwCmd = client.authorizeCommand!!
        assertEquals("mand_active", gwCmd.mandateId)
        assertEquals(userUid, gwCmd.customerId)
    }

    @Test
    fun `processPayment persists the customer id so the payment is visible to its owner`() {
        val userUid = "00000000-0000-0000-0000-000000000007"
        val paymentRepo = FakePaymentRepo()
        val activities = PaymentActivitiesImpl(
            paymentRepo,
            FakePaymentAttemptRepo(),
            FakeRefundRepo(),
            FakeMandateRepo(),
            StubHyperswitchClient(),
            FakeDistributedLock(),
        )

        activities.processPayment(
            ProcessPaymentRequest(orderUid = orderUid, userUid = userUid, amount = BigDecimal("10.00"), currency = "USD"),
        )

        assertEquals(userUid, paymentRepo.createCommand!!.customerId)
    }

    @Test
    fun `processPayment fails non-retryably when the customer has no active mandate`() {
        val client = StubHyperswitchClient()
        val paymentRepo = FakePaymentRepo()
        val activities = PaymentActivitiesImpl(
            paymentRepo,
            FakePaymentAttemptRepo(),
            FakeRefundRepo(),
            FakeMandateRepo(active = null),
            client,
            FakeDistributedLock(),
        )

        val failure = assertFailsWith<ApplicationFailure> {
            activities.processPayment(
                ProcessPaymentRequest(
                    orderUid = orderUid,
                    userUid = "00000000-0000-0000-0000-000000000007",
                    amount = BigDecimal("49.99"),
                    currency = "USD",
                ),
            )
        }

        assertTrue(failure.isNonRetryable)
        assertEquals("NoActivePaymentMandate", failure.type)
        assertNull(client.authorizeCommand)
        assertNull(paymentRepo.createCommand)
    }

    // ─── Stubs / fakes ─────────────────────────────────────────────────────────

    private class StubHyperswitchClient(
        private val paymentResult: HyperswitchPaymentResult = HyperswitchPaymentResult(
            "pay_default", "succeeded", "stripe", 0, 0, null, null,
        ),
        private val refundResult: HyperswitchRefundResult = HyperswitchRefundResult("ref_default", "succeeded", "stripe", null, null),
        private val authorizeError: HyperswitchClientException? = null,
        private val refundError: HyperswitchClientException? = null,
        private val lookupByKey: HyperswitchPaymentLookup? = null,
        private val lookupById: HyperswitchPaymentLookup? = null,
    ) : HyperswitchClient {
        var authorizeCommand: HyperswitchPaymentCommand? = null
        var refundCommand: HyperswitchRefundCommand? = null
        var lookedUpKey: String? = null
        var lookedUpPaymentId: String? = null

        override fun authorizePayment(command: HyperswitchPaymentCommand): HyperswitchPaymentResult {
            authorizeCommand = command
            authorizeError?.let { throw it }
            return paymentResult
        }

        override fun refundPayment(command: HyperswitchRefundCommand): HyperswitchRefundResult {
            refundCommand = command
            refundError?.let { throw it }
            return refundResult
        }

        override fun findPaymentByIdempotencyKey(idempotencyKey: String): HyperswitchPaymentLookup? {
            lookedUpKey = idempotencyKey
            return lookupByKey
        }

        override fun findPaymentById(gatewayPaymentId: String): HyperswitchPaymentLookup? {
            lookedUpPaymentId = gatewayPaymentId
            return lookupById
        }

        override fun setupMandate(command: HyperswitchSetupMandateCommand): HyperswitchSetupMandateResult =
            HyperswitchSetupMandateResult("man_stub", "pm_stub", "active", null, null)

        override fun verifyWebhookSignature(payload: ByteArray, signature: String?): Boolean = error("not used")

        override fun parseWebhookEvent(payload: ByteArray): HyperswitchWebhookEvent = error("not used")
    }

    private class FakeMandateRepo(
        private val active: MandateResult? = PaymentTestFixtures.mandateResult(mandateStatus = "active"),
    ) : MandateRepo {
        override fun findAll(option: MandateQueryOption, pageable: Pageable, context: ActorContext): Page<MandateResult> =
            PageImpl(emptyList())

        override fun findOne(uid: UUID, context: ActorContext): MandateResult = error("findOne not seeded")

        override fun findActiveByCustomerId(customerId: String, context: ActorContext): MandateResult? = active

        override fun create(command: MandateCreateCommand, context: ActorContext): MandateResult =
            error("create not seeded")

        override fun revoke(uid: UUID, context: ActorContext): MandateResult = error("revoke not seeded")
    }

    private class FakePaymentRepo(
        private val existingByKey: PaymentResult? = null,
        private val findOne: PaymentResult? = null,
    ) : PaymentRepo {
        var createCommand: PaymentCreateCommand? = null
        var syncedStatus: String? = null

        override fun findAll(option: PaymentQueryOption, pageable: Pageable, context: ActorContext): Page<PaymentResult> =
            PageImpl(emptyList())

        override fun findOne(uid: UUID, context: ActorContext): PaymentResult =
            findOne ?: error("findOne not seeded")

        override fun findByPaymentId(paymentId: String, context: ActorContext): PaymentResult? = null

        override fun findByIdempotencyKey(idempotencyKey: String, context: ActorContext): PaymentResult? = existingByKey

        override fun create(command: PaymentCreateCommand, context: ActorContext): PaymentResult {
            createCommand = command
            return PaymentTestFixtures.paymentResult(
                status = command.status ?: "requires_payment_method",
                paymentId = command.paymentId ?: "pay_generated",
                amount = command.amount,
                currency = command.currency,
            )
        }

        override fun update(uid: UUID, command: PaymentUpdateCommand, context: ActorContext): PaymentResult =
            error("not used")

        override fun syncGatewayStatus(
            uid: UUID,
            status: String,
            amountCapturable: Long?,
            amountCaptured: Long?,
            context: ActorContext,
        ): PaymentResult {
            syncedStatus = status
            val base = findOne ?: existingByKey ?: error("no payment seeded for syncGatewayStatus")
            return base.copy(status = status, amountCapturable = amountCapturable, amountCaptured = amountCaptured)
        }
    }

    private class FakePaymentAttemptRepo : PaymentAttemptRepo {
        var createCommand: PaymentAttemptCreateCommand? = null

        override fun findAllByPayment(paymentUid: UUID, pageable: Pageable, context: ActorContext): Page<PaymentAttemptResult> =
            PageImpl(emptyList())

        override fun findOne(paymentUid: UUID, uid: UUID, context: ActorContext): PaymentAttemptResult =
            error("not used")

        override fun create(command: PaymentAttemptCreateCommand, context: ActorContext): PaymentAttemptResult {
            createCommand = command
            return PaymentTestFixtures.paymentAttemptResult(status = command.status, paymentId = command.paymentId)
        }
    }

    private class FakeRefundRepo(
        private val existingByKey: RefundResult? = null,
    ) : RefundRepo {
        var createCommand: RefundCreateCommand? = null

        override fun findAll(option: RefundQueryOption, pageable: Pageable, context: ActorContext): Page<RefundResult> =
            PageImpl(emptyList())

        override fun findOne(uid: UUID, context: ActorContext): RefundResult = error("not used")

        override fun findByIdempotencyKey(idempotencyKey: String, context: ActorContext): RefundResult? = existingByKey

        override fun findByRefundId(refundId: String, context: ActorContext): RefundResult? = error("not used")

        override fun create(command: RefundCreateCommand, context: ActorContext): RefundResult {
            createCommand = command
            return PaymentTestFixtures.refundResult(
                refundId = command.refundId ?: "ref_generated",
                status = command.status ?: "pending",
            )
        }

        override fun syncGatewayStatus(
            uid: UUID, status: String, errorCode: String?, errorMessage: String?, context: ActorContext,
        ): RefundResult = error("not used")
    }
}
