package com.codejsha.bookstore.payment.infrastructure.adapter.mysql

import com.codejsha.bookstore.generated.infrastructure.adapter.jooq.tables.references.WEBHOOK_EVENTS
import com.codejsha.bookstore.payment.application.port.repo.WebhookEventRepo
import com.codejsha.bookstore.payment.application.port.repo.WebhookEventResult
import com.codejsha.bookstore.payment.domain.model.command.WebhookEventCreateCommand
import com.codejsha.bookstore.payment.infrastructure.support.utils.bytesToUuid
import com.codejsha.bookstore.payment.infrastructure.support.utils.uuidToBytes
import com.codejsha.platform.shared.data.ActorContext
import org.jooq.DSLContext
import org.jooq.JSON
import org.springframework.stereotype.Repository
import java.time.LocalDateTime
import java.time.ZoneOffset
import java.util.UUID
import kotlin.uuid.Uuid
import kotlin.uuid.toJavaUuid

@Repository
class WebhookEventRepoImpl(
    private val dslContext: DSLContext
) : WebhookEventRepo {
    private val w = WEBHOOK_EVENTS.`as`("w")

    override fun create(command: WebhookEventCreateCommand, context: ActorContext): WebhookEventResult {
        val uid = Uuid.generateV7().toJavaUuid()

        dslContext
            .insertInto(w)
            .set(w.UID, uuidToBytes(uid))
            .set(w.EVENT_ID, command.eventId)
            .set(w.EVENT_TYPE, command.eventType)
            .set(w.OBJECT_TYPE, command.objectType)
            .set(w.OBJECT_ID, command.objectId)
            .set(w.PAYLOAD, JSON.valueOf(command.payload))
            .set(w.SIGNATURE, command.signature)
            .set(w.PROCESSED, 0.toByte())
            .set(w.CREATED_AT, LocalDateTime.now(ZoneOffset.UTC))
            .execute()

        return findOne(uid)
    }

    override fun markProcessed(uid: UUID, context: ActorContext): WebhookEventResult {
        dslContext
            .update(w)
            .set(w.PROCESSED, 1.toByte())
            .set(w.UPDATED_AT, LocalDateTime.now(ZoneOffset.UTC))
            .where(w.UID.eq(uuidToBytes(uid)).and(w.DELETED_AT.isNull))
            .execute()

        return findOne(uid)
    }

    private fun findOne(uid: UUID): WebhookEventResult {
        return dslContext
            .select(w.asterisk())
            .from(w)
            .where(w.UID.eq(uuidToBytes(uid)).and(w.DELETED_AT.isNull))
            .fetchOne { toWebhookEventResult(it) }
            ?: throw NoSuchElementException("Webhook event with uid $uid not found")
    }

    private fun toWebhookEventResult(record: org.jooq.Record): WebhookEventResult {
        return WebhookEventResult(
            id = record.get(w.ID)!!.toLong(),
            uid = bytesToUuid(record.get(w.UID)!!),
            eventId = record.get(w.EVENT_ID)!!,
            eventType = record.get(w.EVENT_TYPE)!!,
            objectType = record.get(w.OBJECT_TYPE)!!,
            objectId = record.get(w.OBJECT_ID)!!,
            processed = record.get(w.PROCESSED)!!.toInt() != 0,
            createdAt = record.get(w.CREATED_AT)!!,
            updatedAt = record.get(w.UPDATED_AT),
        )
    }
}
