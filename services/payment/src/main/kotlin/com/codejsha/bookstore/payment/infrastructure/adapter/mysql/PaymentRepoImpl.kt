package com.codejsha.bookstore.payment.infrastructure.adapter.mysql

import com.codejsha.bookstore.payment.infrastructure.support.utils.bytesToUuid
import com.codejsha.bookstore.payment.infrastructure.support.utils.parseJson
import com.codejsha.bookstore.payment.infrastructure.support.utils.toJson
import com.codejsha.bookstore.payment.infrastructure.support.utils.uuidToBytes

import com.codejsha.bookstore.payment.application.port.repo.PaymentRepo
import com.codejsha.bookstore.payment.application.port.repo.PaymentResult
import com.codejsha.bookstore.payment.domain.constant.PaymentStatus
import com.codejsha.bookstore.payment.domain.model.command.PaymentCreateCommand
import com.codejsha.bookstore.payment.domain.model.command.PaymentUpdateCommand
import com.codejsha.bookstore.payment.domain.model.option.PaymentQueryOption
import com.codejsha.bookstore.generated.infrastructure.adapter.jooq.tables.references.PAYMENTS
import com.codejsha.platform.shared.data.ActorContext
import com.codejsha.platform.shared.jooq.buildOrderBy
import org.jooq.Condition
import org.jooq.DSLContext
import org.jooq.impl.DSL
import org.springframework.data.domain.Page
import org.springframework.data.domain.PageImpl
import org.springframework.data.domain.Pageable
import org.springframework.stereotype.Repository
import java.time.LocalDateTime
import java.time.ZoneOffset
import java.util.UUID
import kotlin.uuid.Uuid
import kotlin.uuid.toJavaUuid

@Repository
class PaymentRepoImpl(
    private val dslContext: DSLContext
) : PaymentRepo {
    private val p = PAYMENTS.`as`("p")

    private val idempotencyKeyAliased = DSL.field(DSL.name("p", "idempotency_key"), String::class.java)
    private val idempotencyKeyColumn = DSL.field(DSL.name("idempotency_key"), String::class.java)

    override fun findAll(
        option: PaymentQueryOption, pageable: Pageable, context: ActorContext
    ): Page<PaymentResult> {
        var condition: Condition = p.DELETED_AT.isNull
        option.customerId?.let { condition = condition.and(p.CUSTOMER_ID.eq(it)) }
        option.status?.let { condition = condition.and(p.STATUS.eq(it)) }
        option.connector?.let { condition = condition.and(p.CONNECTOR.eq(it)) }

        val content = dslContext
            .select(p.asterisk())
            .from(p)
            .where(condition)
            .orderBy(buildOrderBy(pageable))
            .limit(pageable.pageSize)
            .offset(pageable.offset)
            .fetch { toPaymentResult(it) }

        val total = dslContext
            .selectCount()
            .from(p)
            .where(condition)
            .fetchOneInto(Int::class.java) ?: 0

        return PageImpl(content, pageable, total.toLong())
    }

    override fun findOne(uid: UUID, context: ActorContext): PaymentResult {
        return dslContext
            .select(p.asterisk())
            .from(p)
            .where(p.UID.eq(uuidToBytes(uid)).and(p.DELETED_AT.isNull))
            .fetchOne { toPaymentResult(it) }
            ?: throw NoSuchElementException("Payment with uid $uid not found")
    }

    override fun findByPaymentId(paymentId: String, context: ActorContext): PaymentResult? {
        return dslContext
            .select(p.asterisk())
            .from(p)
            .where(p.PAYMENT_ID.eq(paymentId).and(p.DELETED_AT.isNull))
            .fetchOne { toPaymentResult(it) }
    }

    override fun findByIdempotencyKey(idempotencyKey: String, context: ActorContext): PaymentResult? {
        return dslContext
            .select(p.asterisk())
            .from(p)
            .where(idempotencyKeyAliased.eq(idempotencyKey).and(p.DELETED_AT.isNull))
            .fetchOne { toPaymentResult(it) }
    }

    override fun create(command: PaymentCreateCommand, context: ActorContext): PaymentResult {
        val uid = Uuid.generateV7().toJavaUuid()
        val paymentId = command.paymentId
            ?: "pay_${UUID.randomUUID().toString().replace("-", "").take(24)}"
        val now = LocalDateTime.now(ZoneOffset.UTC)
        val captured = command.status == PaymentStatus.SUCCEEDED.value ||
            command.status == PaymentStatus.PARTIALLY_CAPTURED.value

        dslContext
            .insertInto(p)
            .set(p.UID, uuidToBytes(uid))
            .set(p.PAYMENT_ID, paymentId)
            .set(idempotencyKeyColumn, command.idempotencyKey)
            .set(p.CUSTOMER_ID, command.customerId)
            .set(p.CONNECTOR, command.connector)
            .set(p.AMOUNT, command.amount)
            .set(p.AMOUNT_CAPTURABLE, command.amountCapturable)
            .set(p.AMOUNT_CAPTURED, command.amountCaptured)
            .set(p.CURRENCY, command.currency)
            .set(p.STATUS, command.status ?: "requires_payment_method")
            .set(p.CAPTURE_METHOD, "automatic")
            .set(p.OFF_SESSION, 0.toByte())
            .set(p.PAYMENT_METHOD, command.paymentMethod)
            .set(p.PAYMENT_METHOD_TYPE, command.paymentMethodType)
            .set(p.AUTHENTICATION_TYPE, command.authenticationType)
            .set(p.SETUP_FUTURE_USAGE, command.setupFutureUsage)
            .set(p.DESCRIPTION, command.description)
            .set(p.RETURN_URL, command.returnUrl)
            .set(p.BILLING_ADDRESS, toJson(command.billingAddress))
            .set(p.SHIPPING_ADDRESS, toJson(command.shippingAddress))
            .set(p.METADATA, toJson(command.metadata))
            .set(p.ERROR_CODE, command.errorCode)
            .set(p.ERROR_MESSAGE, command.errorMessage)
            .apply { if (captured) set(p.CAPTURED_AT, now) }
            .set(p.CREATED_AT, now)
            .execute()

        return findOne(uid, context)
    }

    override fun update(uid: UUID, command: PaymentUpdateCommand, context: ActorContext): PaymentResult {
        val updated = dslContext
            .update(p)
            .apply {
                command.amount?.let { set(p.AMOUNT, it) }
                command.currency?.let { set(p.CURRENCY, it) }
                command.customerId?.let { set(p.CUSTOMER_ID, it) }
                command.paymentMethod?.let { set(p.PAYMENT_METHOD, it) }
                command.paymentMethodType?.let { set(p.PAYMENT_METHOD_TYPE, it) }
                command.authenticationType?.let { set(p.AUTHENTICATION_TYPE, it) }
                command.setupFutureUsage?.let { set(p.SETUP_FUTURE_USAGE, it) }
                command.description?.let { set(p.DESCRIPTION, it) }
                command.returnUrl?.let { set(p.RETURN_URL, it) }
                command.billingAddress?.let { set(p.BILLING_ADDRESS, toJson(it)!!) }
                command.shippingAddress?.let { set(p.SHIPPING_ADDRESS, toJson(it)!!) }
                command.metadata?.let { set(p.METADATA, toJson(it)!!) }
            }
            .set(p.UPDATED_AT, LocalDateTime.now(ZoneOffset.UTC))
            .where(p.UID.eq(uuidToBytes(uid)).and(p.DELETED_AT.isNull))
            .execute()

        check(updated > 0) { "Payment with uid $uid not found" }
        return findOne(uid, context)
    }

    override fun syncGatewayStatus(
        uid: UUID,
        status: String,
        amountCapturable: Long?,
        amountCaptured: Long?,
        context: ActorContext,
    ): PaymentResult {
        val captured = status == PaymentStatus.SUCCEEDED.value ||
            status == PaymentStatus.PARTIALLY_CAPTURED.value

        val updated = dslContext
            .update(p)
            .set(p.STATUS, status)
            .apply {
                amountCapturable?.let { set(p.AMOUNT_CAPTURABLE, it) }
                amountCaptured?.let { set(p.AMOUNT_CAPTURED, it) }
                if (captured) set(p.CAPTURED_AT, LocalDateTime.now(ZoneOffset.UTC))
            }
            .set(p.UPDATED_AT, LocalDateTime.now(ZoneOffset.UTC))
            .where(p.UID.eq(uuidToBytes(uid)).and(p.DELETED_AT.isNull))
            .execute()

        check(updated > 0) { "Payment with uid $uid not found" }
        return findOne(uid, context)
    }

    private fun toPaymentResult(record: org.jooq.Record): PaymentResult {
        return PaymentResult(
            id = record.get(p.ID)!!.toLong(),
            uid = bytesToUuid(record.get(p.UID)!!),
            paymentId = record.get(p.PAYMENT_ID)!!,
            merchantId = record.get(p.MERCHANT_ID),
            profileId = record.get(p.PROFILE_ID),
            customerId = record.get(p.CUSTOMER_ID),
            paymentMethodId = record.get(p.PAYMENT_METHOD_ID),
            mandateId = record.get(p.MANDATE_ID),
            connector = record.get(p.CONNECTOR),
            connectorTransactionId = record.get(p.CONNECTOR_TRANSACTION_ID),
            amount = record.get(p.AMOUNT)!!,
            amountCapturable = record.get(p.AMOUNT_CAPTURABLE),
            amountCaptured = record.get(p.AMOUNT_CAPTURED),
            surchargeAmount = record.get(p.SURCHARGE_AMOUNT),
            taxAmount = record.get(p.TAX_AMOUNT),
            currency = record.get(p.CURRENCY)!!,
            status = record.get(p.STATUS)!!,
            captureMethod = record.get(p.CAPTURE_METHOD)!!,
            authenticationType = record.get(p.AUTHENTICATION_TYPE),
            paymentMethod = record.get(p.PAYMENT_METHOD),
            paymentMethodType = record.get(p.PAYMENT_METHOD_TYPE),
            clientSecret = record.get(p.CLIENT_SECRET),
            setupFutureUsage = record.get(p.SETUP_FUTURE_USAGE),
            offSession = record.get(p.OFF_SESSION)?.toInt() == 1,
            description = record.get(p.DESCRIPTION),
            returnUrl = record.get(p.RETURN_URL),
            statementDescriptor = record.get(p.STATEMENT_DESCRIPTOR),
            billingAddress = record.get(p.BILLING_ADDRESS)?.let { parseJson(it) },
            shippingAddress = record.get(p.SHIPPING_ADDRESS)?.let { parseJson(it) },
            metadata = record.get(p.METADATA)?.let { parseJson(it) },
            errorCode = record.get(p.ERROR_CODE),
            errorMessage = record.get(p.ERROR_MESSAGE),
            confirmedAt = record.get(p.CONFIRMED_AT),
            capturedAt = record.get(p.CAPTURED_AT),
            cancelledAt = record.get(p.CANCELLED_AT),
            createdAt = record.get(p.CREATED_AT)!!,
            updatedAt = record.get(p.UPDATED_AT),
        )
    }
}
