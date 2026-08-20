package com.codejsha.bookstore.payment.infrastructure.adapter.mysql

import com.codejsha.bookstore.payment.infrastructure.support.utils.bytesToUuid
import com.codejsha.bookstore.payment.infrastructure.support.utils.parseJson
import com.codejsha.bookstore.payment.infrastructure.support.utils.toJson
import com.codejsha.bookstore.payment.infrastructure.support.utils.uuidToBytes

import com.codejsha.bookstore.payment.application.port.repo.PaymentMethodRepo
import com.codejsha.bookstore.payment.application.port.repo.PaymentMethodResult
import com.codejsha.bookstore.payment.domain.model.command.PaymentMethodCreateCommand
import com.codejsha.bookstore.payment.domain.model.command.PaymentMethodUpdateCommand
import com.codejsha.bookstore.payment.domain.model.option.PaymentMethodQueryOption
import com.codejsha.bookstore.generated.infrastructure.adapter.jooq.tables.references.CUSTOMERS
import com.codejsha.bookstore.generated.infrastructure.adapter.jooq.tables.references.PAYMENT_METHODS
import com.codejsha.platform.shared.data.ActorContext
import com.codejsha.platform.shared.jooq.buildOrderBy
import org.jooq.Condition
import org.jooq.DSLContext
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
class PaymentMethodRepoImpl(
    private val dslContext: DSLContext
) : PaymentMethodRepo {
    private val pm = PAYMENT_METHODS.`as`("pm")
    private val c = CUSTOMERS.`as`("c")

    override fun findAllByCustomer(
        customerUid: UUID, option: PaymentMethodQueryOption, pageable: Pageable, context: ActorContext
    ): Page<PaymentMethodResult> {
        val customerId = resolveCustomerId(customerUid)

        var condition: Condition = pm.CUSTOMER_ID.eq(customerId).and(pm.DELETED_AT.isNull)
        option.paymentMethod?.let { condition = condition.and(pm.PAYMENT_METHOD.eq(it)) }

        val content = dslContext
            .select(pm.asterisk())
            .from(pm)
            .where(condition)
            .orderBy(buildOrderBy(pageable))
            .limit(pageable.pageSize)
            .offset(pageable.offset)
            .fetch { toPaymentMethodResult(it) }

        val total = dslContext
            .selectCount()
            .from(pm)
            .where(condition)
            .fetchOneInto(Int::class.java) ?: 0

        return PageImpl(content, pageable, total.toLong())
    }

    override fun findOne(customerUid: UUID, uid: UUID, context: ActorContext): PaymentMethodResult {
        val customerId = resolveCustomerId(customerUid)
        return dslContext
            .select(pm.asterisk())
            .from(pm)
            .where(
                pm.UID.eq(uuidToBytes(uid))
                    .and(pm.CUSTOMER_ID.eq(customerId))
                    .and(pm.DELETED_AT.isNull)
            )
            .fetchOne { toPaymentMethodResult(it) }
            ?: throw NoSuchElementException("PaymentMethod with uid $uid not found")
    }

    override fun create(
        customerUid: UUID, command: PaymentMethodCreateCommand, context: ActorContext
    ): PaymentMethodResult {
        val customerId = resolveCustomerId(customerUid)
        val uid = Uuid.generateV7().toJavaUuid()
        val paymentMethodId = "pm_${UUID.randomUUID().toString().replace("-", "").take(24)}"

        dslContext
            .insertInto(pm)
            .set(pm.UID, uuidToBytes(uid))
            .set(pm.PAYMENT_METHOD_ID, paymentMethodId)
            .set(pm.CUSTOMER_ID, customerId)
            .set(pm.PAYMENT_METHOD, command.paymentMethod)
            .set(pm.PAYMENT_METHOD_TYPE, command.paymentMethodType)
            .set(pm.PAYMENT_METHOD_ISSUER, command.paymentMethodIssuer)
            .set(pm.CARD_NETWORK, command.cardNetwork)
            .set(pm.CARD_LAST4, command.cardLast4)
            .set(pm.CARD_EXP_MONTH, command.cardExpMonth?.let { org.jooq.types.UByte.valueOf(it) })
            .set(pm.CARD_EXP_YEAR, command.cardExpYear?.let { org.jooq.types.UShort.valueOf(it) })
            .set(pm.CARD_HOLDER_NAME, command.cardHolderName)
            .set(pm.IS_DEFAULT, if (command.isDefault == true) 1.toByte() else 0.toByte())
            .set(pm.METADATA, toJson(command.metadata))
            .set(pm.CREATED_AT, LocalDateTime.now(ZoneOffset.UTC))
            .execute()

        return findOne(customerUid, uid, context)
    }

    override fun update(
        customerUid: UUID, uid: UUID, command: PaymentMethodUpdateCommand, context: ActorContext
    ): PaymentMethodResult {
        val customerId = resolveCustomerId(customerUid)

        val updated = dslContext
            .update(pm)
            .apply {
                command.paymentMethod?.let { set(pm.PAYMENT_METHOD, it) }
                command.paymentMethodType?.let { set(pm.PAYMENT_METHOD_TYPE, it) }
                command.paymentMethodIssuer?.let { set(pm.PAYMENT_METHOD_ISSUER, it) }
                command.cardNetwork?.let { set(pm.CARD_NETWORK, it) }
                command.cardLast4?.let { set(pm.CARD_LAST4, it) }
                command.cardExpMonth?.let { set(pm.CARD_EXP_MONTH, org.jooq.types.UByte.valueOf(it)) }
                command.cardExpYear?.let { set(pm.CARD_EXP_YEAR, org.jooq.types.UShort.valueOf(it)) }
                command.cardHolderName?.let { set(pm.CARD_HOLDER_NAME, it) }
                command.isDefault?.let { set(pm.IS_DEFAULT, if (it) 1.toByte() else 0.toByte()) }
                command.metadata?.let { set(pm.METADATA, toJson(it)!!) }
            }
            .set(pm.UPDATED_AT, LocalDateTime.now(ZoneOffset.UTC))
            .where(
                pm.UID.eq(uuidToBytes(uid))
                    .and(pm.CUSTOMER_ID.eq(customerId))
                    .and(pm.DELETED_AT.isNull)
            )
            .execute()

        check(updated > 0) { "PaymentMethod with uid $uid not found" }
        return findOne(customerUid, uid, context)
    }

    override fun delete(customerUid: UUID, uid: UUID, context: ActorContext) {
        val customerId = resolveCustomerId(customerUid)
        dslContext
            .update(pm)
            .set(pm.DELETED_AT, LocalDateTime.now(ZoneOffset.UTC))
            .where(
                pm.UID.eq(uuidToBytes(uid))
                    .and(pm.CUSTOMER_ID.eq(customerId))
                    .and(pm.DELETED_AT.isNull)
            )
            .execute()
    }

    private fun resolveCustomerId(customerUid: UUID): String {
        return dslContext
            .select(c.CUSTOMER_ID)
            .from(c)
            .where(c.UID.eq(uuidToBytes(customerUid)).and(c.DELETED_AT.isNull))
            .fetchOneInto(String::class.java)
            ?: throw NoSuchElementException("Customer with uid $customerUid not found")
    }

    private fun toPaymentMethodResult(record: org.jooq.Record): PaymentMethodResult {
        return PaymentMethodResult(
            id = record.get(pm.ID)!!.toLong(),
            uid = bytesToUuid(record.get(pm.UID)!!),
            paymentMethodId = record.get(pm.PAYMENT_METHOD_ID)!!,
            customerId = record.get(pm.CUSTOMER_ID)!!,
            paymentMethod = record.get(pm.PAYMENT_METHOD)!!,
            paymentMethodType = record.get(pm.PAYMENT_METHOD_TYPE),
            paymentMethodIssuer = record.get(pm.PAYMENT_METHOD_ISSUER),
            cardNetwork = record.get(pm.CARD_NETWORK),
            cardLast4 = record.get(pm.CARD_LAST4),
            cardExpMonth = record.get(pm.CARD_EXP_MONTH)?.toInt(),
            cardExpYear = record.get(pm.CARD_EXP_YEAR)?.toInt(),
            cardHolderName = record.get(pm.CARD_HOLDER_NAME),
            isDefault = record.get(pm.IS_DEFAULT)?.toInt() == 1,
            metadata = record.get(pm.METADATA)?.let { parseJson(it) },
            createdAt = record.get(pm.CREATED_AT)!!,
            updatedAt = record.get(pm.UPDATED_AT),
        )
    }
}
