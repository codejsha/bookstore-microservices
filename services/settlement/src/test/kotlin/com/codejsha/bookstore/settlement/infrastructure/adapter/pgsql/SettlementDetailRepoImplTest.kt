package com.codejsha.bookstore.settlement.infrastructure.adapter.pgsql

import com.codejsha.bookstore.settlement.application.port.repo.SettlementDetailRepo
import com.codejsha.bookstore.settlement.application.port.repo.SettlementDetailResult
import com.codejsha.bookstore.settlement.support.SettlementPostgresIntegrationTest
import com.codejsha.platform.shared.data.ActorContext
import com.codejsha.platform.shared.data.ActorType
import org.junit.jupiter.api.Test
import org.springframework.beans.factory.annotation.Autowired
import java.time.LocalDate
import java.time.LocalDateTime
import java.util.UUID
import kotlin.test.assertEquals
import kotlin.test.assertNull
import kotlin.test.assertTrue

class SettlementDetailRepoImplTest : SettlementPostgresIntegrationTest() {

    @Autowired
    private lateinit var repo: SettlementDetailRepo

    private val context = ActorContext(actorId = 0L, actorType = ActorType.USER)

    private val date: LocalDate = LocalDate.of(2026, 9, 10)

    @Test
    fun `findAllByBucket_matchingLines_listsThemOldestFirst`() {
        val late = seed("PAYMENT", "pay_late", occurredAt = at(15, 0))
        val early = seed("PAYMENT", "pay_early", occurredAt = at(9, 0))
        val middle = seed("REFUND", "ref_middle", occurredAt = at(12, 0))

        val lines = repo.findAllByBucket(date, "KRW", "card", context)

        assertEquals(listOf(early, middle, late), lines.map { it.uid })
        assertEquals(listOf("PAYMENT", "REFUND", "PAYMENT"), lines.map { it.sourceType })
    }

    @Test
    fun `findAllByBucket_linesOfOtherBuckets_areLeftOut`() {
        val match = seed("PAYMENT", "pay_match")
        seed("PAYMENT", "pay_other_date", settlementDate = date.plusDays(1))
        seed("PAYMENT", "pay_other_currency", currency = "USD")
        seed("PAYMENT", "pay_other_method", method = "wallet")
        seed("PAYMENT", "pay_no_method", method = null)

        val lines = repo.findAllByBucket(date, "KRW", "card", context)

        assertEquals(listOf(match), lines.map { it.uid })
    }

    @Test
    fun `findAllByBucket_nullPaymentMethod_listsOnlyLinesWithoutAMethod`() {
        val noMethod = seed("PAYMENT", "pay_no_method", method = null)
        seed("PAYMENT", "pay_card", method = "card")

        val lines = repo.findAllByBucket(date, "KRW", null, context)

        assertEquals(listOf(noMethod), lines.map { it.uid })
        assertNull(lines.single().paymentMethod)
    }

    @Test
    fun `findAllByBucket_softDeletedLine_isLeftOut`() {
        val live = seed("PAYMENT", "pay_live")
        seed("PAYMENT", "pay_deleted", deletedAt = LocalDateTime.of(2026, 9, 12, 0, 0))

        val lines = repo.findAllByBucket(date, "KRW", "card", context)

        assertEquals(listOf(live), lines.map { it.uid })
    }

    @Test
    fun `findAllByBucket_noLinesForBucket_isEmpty`() {
        seed("PAYMENT", "pay_1")

        assertTrue(repo.findAllByBucket(date.minusDays(1), "KRW", "card", context).isEmpty())
    }

    @Test
    fun `findAllByBucket_storedLine_mapsEveryColumn`() {
        val occurredAt = LocalDateTime.of(2026, 9, 10, 3, 4, 5, 123_456_000)
        val createdAt = LocalDateTime.of(2026, 9, 11, 2, 0, 0, 1_000)
        val updatedAt = LocalDateTime.of(2026, 9, 11, 2, 1, 0, 999_999_000)
        val uid = seed(
            "REFUND", "ref_1", paymentId = "pay_1", amount = 3_000,
            occurredAt = occurredAt, createdAt = createdAt, updatedAt = updatedAt,
        )

        val line = repo.findAllByBucket(date, "KRW", "card", context).single()

        assertEquals(
            SettlementDetailResult(
                id = line.id,
                uid = uid,
                settlementDate = date,
                sourceType = "REFUND",
                sourceId = "ref_1",
                paymentId = "pay_1",
                amount = 3_000,
                currency = "KRW",
                paymentMethod = "card",
                occurredAt = occurredAt,
                createdAt = createdAt,
                updatedAt = updatedAt,
            ),
            line,
        )
        assertEquals(
            settlementJdbc.queryForObject("SELECT id FROM settlement_detail WHERE uid = ?", Long::class.java, uid),
            line.id,
        )
    }

    @Test
    fun `findAllByBucket_lineNeverUpdated_mapsNullUpdatedAt`() {
        seed("PAYMENT", "pay_1")

        assertNull(repo.findAllByBucket(date, "KRW", "card", context).single().updatedAt)
    }

    private fun at(hour: Int, minute: Int): LocalDateTime = date.atTime(hour, minute)

    private fun seed(
        sourceType: String,
        sourceId: String,
        paymentId: String = sourceId,
        amount: Long = 1_000,
        settlementDate: LocalDate = date,
        currency: String = "KRW",
        method: String? = "card",
        occurredAt: LocalDateTime = at(10, 0),
        createdAt: LocalDateTime = LocalDateTime.of(2026, 9, 11, 2, 0),
        updatedAt: LocalDateTime? = null,
        deletedAt: LocalDateTime? = null,
    ): UUID {
        val uid = UUID.randomUUID()
        settlementJdbc.update(
            """
            INSERT INTO settlement_detail
                (uid, settlement_date, source_type, source_id, payment_id,
                 amount, currency, payment_method, occurred_at, created_at, updated_at, deleted_at)
            VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
            """.trimIndent(),
            uid, settlementDate, sourceType, sourceId, paymentId,
            amount, currency, method, occurredAt, createdAt, updatedAt, deletedAt,
        )
        return uid
    }
}
