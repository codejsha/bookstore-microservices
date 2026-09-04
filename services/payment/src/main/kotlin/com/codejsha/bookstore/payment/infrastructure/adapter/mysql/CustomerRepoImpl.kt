package com.codejsha.bookstore.payment.infrastructure.adapter.mysql

import com.codejsha.bookstore.payment.infrastructure.support.utils.bytesToUuid
import com.codejsha.bookstore.payment.infrastructure.support.utils.parseJson
import com.codejsha.bookstore.payment.infrastructure.support.utils.toJson
import com.codejsha.bookstore.payment.infrastructure.support.utils.uuidToBytes

import com.codejsha.bookstore.payment.application.port.repo.CustomerRepo
import com.codejsha.bookstore.payment.application.port.repo.CustomerResult
import com.codejsha.bookstore.payment.domain.model.command.CustomerCreateCommand
import com.codejsha.bookstore.payment.domain.model.command.CustomerUpdateCommand
import com.codejsha.bookstore.payment.domain.model.option.CustomerQueryOption
import com.codejsha.bookstore.generated.infrastructure.adapter.jooq.tables.references.CUSTOMERS
import com.codejsha.platform.shared.data.ActorContext
import com.codejsha.platform.shared.jooq.buildOrderBy
import org.jooq.Condition
import org.jooq.DSLContext
import org.springframework.data.domain.Page
import org.springframework.data.domain.PageImpl
import org.springframework.data.domain.Pageable
import org.springframework.stereotype.Repository
import java.time.LocalDateTime
import java.util.UUID
import kotlin.uuid.Uuid
import kotlin.uuid.toJavaUuid

@Repository
class CustomerRepoImpl(
    private val dslContext: DSLContext
) : CustomerRepo {
    private val c = CUSTOMERS.`as`("c")

    override fun findAll(
        option: CustomerQueryOption, pageable: Pageable, context: ActorContext
    ): Page<CustomerResult> {
        var condition: Condition = c.DELETED_AT.isNull
        option.customerId?.let { condition = condition.and(c.CUSTOMER_ID.eq(it)) }
        option.email?.let { condition = condition.and(c.EMAIL.eq(it)) }

        val content = dslContext
            .select(c.asterisk())
            .from(c)
            .where(condition)
            .orderBy(buildOrderBy(pageable))
            .limit(pageable.pageSize)
            .offset(pageable.offset)
            .fetch { toCustomerResult(it) }

        val total = dslContext
            .selectCount()
            .from(c)
            .where(condition)
            .fetchOneInto(Int::class.java) ?: 0

        return PageImpl(content, pageable, total.toLong())
    }

    override fun findOne(uid: UUID, context: ActorContext): CustomerResult {
        return dslContext
            .select(c.asterisk())
            .from(c)
            .where(c.UID.eq(uuidToBytes(uid)).and(c.DELETED_AT.isNull))
            .fetchOne { toCustomerResult(it) }
            ?: throw NoSuchElementException("Customer with uid $uid not found")
    }

    override fun findByCustomerId(customerId: String, context: ActorContext): CustomerResult? {
        return dslContext
            .select(c.asterisk())
            .from(c)
            .where(c.CUSTOMER_ID.eq(customerId).and(c.DELETED_AT.isNull))
            .fetchOne { toCustomerResult(it) }
    }

    override fun create(command: CustomerCreateCommand, context: ActorContext): CustomerResult {
        val now = LocalDateTime.now()
        val uid = Uuid.generateV7().toJavaUuid()

        dslContext
            .insertInto(c)
            .set(c.UID, uuidToBytes(uid))
            .set(c.CUSTOMER_ID, command.customerId)
            .set(c.NAME, command.name)
            .set(c.EMAIL, command.email)
            .set(c.PHONE, command.phone)
            .set(c.PHONE_COUNTRY_CODE, command.phoneCountryCode)
            .set(c.DESCRIPTION, command.description)
            .set(c.METADATA, toJson(command.metadata))
            .set(c.DEFAULT_BILLING_ADDRESS, toJson(command.defaultBillingAddress))
            .set(c.DEFAULT_SHIPPING_ADDRESS, toJson(command.defaultShippingAddress))
            .set(c.CREATED_AT, now)
            .execute()

        return findOne(uid, context)
    }

    override fun update(uid: UUID, command: CustomerUpdateCommand, context: ActorContext): CustomerResult {
        val updated = dslContext
            .update(c)
            .apply {
                command.name?.let { set(c.NAME, it) }
                command.email?.let { set(c.EMAIL, it) }
                command.phone?.let { set(c.PHONE, it) }
                command.phoneCountryCode?.let { set(c.PHONE_COUNTRY_CODE, it) }
                command.description?.let { set(c.DESCRIPTION, it) }
                command.metadata?.let { set(c.METADATA, toJson(it)!!) }
                command.defaultBillingAddress?.let { set(c.DEFAULT_BILLING_ADDRESS, toJson(it)!!) }
                command.defaultShippingAddress?.let { set(c.DEFAULT_SHIPPING_ADDRESS, toJson(it)!!) }
            }
            .set(c.UPDATED_AT, LocalDateTime.now())
            .where(c.UID.eq(uuidToBytes(uid)).and(c.DELETED_AT.isNull))
            .execute()

        check(updated > 0) { "Customer with uid $uid not found" }
        return findOne(uid, context)
    }

    override fun delete(uid: UUID, context: ActorContext) {
        dslContext
            .update(c)
            .set(c.DELETED_AT, LocalDateTime.now())
            .where(c.UID.eq(uuidToBytes(uid)).and(c.DELETED_AT.isNull))
            .execute()
    }

    private fun toCustomerResult(record: org.jooq.Record): CustomerResult {
        return CustomerResult(
            id = record.get(c.ID)!!.toLong(),
            uid = bytesToUuid(record.get(c.UID)!!),
            customerId = record.get(c.CUSTOMER_ID)!!,
            name = record.get(c.NAME),
            email = record.get(c.EMAIL),
            phone = record.get(c.PHONE),
            phoneCountryCode = record.get(c.PHONE_COUNTRY_CODE),
            description = record.get(c.DESCRIPTION),
            metadata = record.get(c.METADATA)?.let { parseJson(it) },
            defaultBillingAddress = record.get(c.DEFAULT_BILLING_ADDRESS)?.let { parseJson(it) },
            defaultShippingAddress = record.get(c.DEFAULT_SHIPPING_ADDRESS)?.let { parseJson(it) },
            createdAt = record.get(c.CREATED_AT)!!,
            updatedAt = record.get(c.UPDATED_AT),
        )
    }
}
