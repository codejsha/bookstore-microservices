package com.codejsha.bookstore.settlement.infrastructure.adapter.pgsql

import com.codejsha.bookstore.settlement.application.port.repo.DailySettlementRepo
import com.codejsha.bookstore.settlement.application.port.repo.DailySettlementResult
import com.codejsha.bookstore.settlement.domain.model.option.SettlementQueryOption
import com.codejsha.bookstore.settlement.support.SettlementPostgresIntegrationTest
import com.codejsha.platform.shared.data.ActorContext
import com.codejsha.platform.shared.data.ActorType
import org.junit.jupiter.api.Test
import org.springframework.beans.factory.annotation.Autowired
import org.springframework.data.domain.PageRequest
import org.springframework.data.domain.Pageable
import org.springframework.data.domain.Sort
import java.time.LocalDate
import java.time.LocalDateTime
import java.util.UUID
import kotlin.test.assertEquals
import kotlin.test.assertFailsWith
import kotlin.test.assertNull
import kotlin.test.assertTrue

class DailySettlementRepoImplTest : SettlementPostgresIntegrationTest() {

    @Autowired
    private lateinit var repo: DailySettlementRepo

    private val context = ActorContext(actorId = 0L, actorType = ActorType.USER)

    private val day1: LocalDate = LocalDate.of(2026, 9, 1)
    private val day2: LocalDate = LocalDate.of(2026, 9, 2)
    private val day3: LocalDate = LocalDate.of(2026, 9, 3)

    @Test
    fun `findAll_noFilterUnpaged_listsEveryLiveBucket`() {
        val krw = seed(day1, "KRW", "card", "CONFIRMED")
        val usd = seed(day1, "USD", "wallet", "OPEN")
        val next = seed(day2, "KRW", "card", "DISCREPANCY")
        seed(day3, "KRW", "card", "CONFIRMED", deletedAt = LocalDateTime.of(2026, 9, 4, 0, 0))

        val page = repo.findAll(SettlementQueryOption(), Pageable.unpaged(), context)

        assertEquals(3, page.totalElements)
        assertEquals(setOf(krw, usd, next), page.content.map { it.uid }.toSet())
    }

    @Test
    fun `findAll_dateFilter_listsOnlyThatDate`() {
        val krw = seed(day1, "KRW", "card", "CONFIRMED")
        val usd = seed(day1, "USD", "wallet", "OPEN")
        seed(day2, "KRW", "card", "CONFIRMED")

        val page = repo.findAll(SettlementQueryOption(settlementDate = day1), Pageable.unpaged(), context)

        assertEquals(2, page.totalElements)
        assertEquals(setOf(krw, usd), page.content.map { it.uid }.toSet())
        assertTrue(page.content.all { it.settlementDate == day1 })
    }

    @Test
    fun `findAll_statusFilter_listsOnlyThatStatus`() {
        seed(day1, "KRW", "card", "CONFIRMED")
        val open = seed(day1, "USD", "wallet", "OPEN")
        val discrepancy = seed(day2, "KRW", "card", "DISCREPANCY")

        val openPage = repo.findAll(SettlementQueryOption(status = "OPEN"), Pageable.unpaged(), context)
        val discrepancyPage = repo.findAll(SettlementQueryOption(status = "DISCREPANCY"), Pageable.unpaged(), context)

        assertEquals(listOf(open), openPage.content.map { it.uid })
        assertEquals(1, openPage.totalElements)
        assertEquals(listOf(discrepancy), discrepancyPage.content.map { it.uid })
    }

    @Test
    fun `findAll_dateAndStatusFilter_listsTheIntersection`() {
        val match = seed(day1, "KRW", "card", "CONFIRMED")
        seed(day1, "USD", "wallet", "OPEN")
        seed(day2, "KRW", "card", "CONFIRMED")

        val page = repo.findAll(
            SettlementQueryOption(settlementDate = day1, status = "CONFIRMED"), Pageable.unpaged(), context,
        )

        assertEquals(listOf(match), page.content.map { it.uid })
        assertEquals(1, page.totalElements)
    }

    @Test
    fun `findAll_filterMatchesNothing_isEmpty`() {
        seed(day1, "KRW", "card", "CONFIRMED")

        val page = repo.findAll(SettlementQueryOption(settlementDate = day3), Pageable.unpaged(), context)

        assertTrue(page.content.isEmpty())
        assertEquals(0, page.totalElements)
    }

    @Test
    fun `findAll_pagedAndSortedByDateDescending_slicesInOrderAndKeepsTheTotal`() {
        val first = seed(day1, "KRW", "card", "CONFIRMED")
        val second = seed(day2, "KRW", "card", "CONFIRMED")
        val third = seed(day3, "KRW", "card", "CONFIRMED")
        val sort = Sort.by(Sort.Order.desc("settlement_date"))

        val pageOne = repo.findAll(SettlementQueryOption(), PageRequest.of(0, 2, sort), context)
        val pageTwo = repo.findAll(SettlementQueryOption(), PageRequest.of(1, 2, sort), context)
        val beyond = repo.findAll(SettlementQueryOption(), PageRequest.of(5, 2, sort), context)

        assertEquals(listOf(third, second), pageOne.content.map { it.uid })
        assertEquals(listOf(first), pageTwo.content.map { it.uid })
        assertEquals(3, pageOne.totalElements)
        assertEquals(3, pageTwo.totalElements)
        assertEquals(2, pageOne.totalPages)
        assertTrue(beyond.content.isEmpty())
        assertEquals(3, beyond.totalElements)
    }

    @Test
    fun `findAll_sortedByGrossAmountAscending_ordersByThatColumn`() {
        val large = seed(day1, "KRW", "card", "CONFIRMED", gross = 9_000)
        val small = seed(day2, "KRW", "card", "CONFIRMED", gross = 1_000)
        val medium = seed(day3, "KRW", "card", "CONFIRMED", gross = 5_000)

        val page = repo.findAll(
            SettlementQueryOption(), PageRequest.of(0, 10, Sort.by(Sort.Order.asc("gross_amount"))), context,
        )

        assertEquals(listOf(small, medium, large), page.content.map { it.uid })
    }

    @Test
    fun `findOne_existingUid_mapsEveryColumn`() {
        val createdAt = LocalDateTime.of(2026, 9, 2, 2, 0, 0, 123_456_000)
        val updatedAt = LocalDateTime.of(2026, 9, 2, 2, 5, 0, 654_321_000)
        val uid = seed(
            day1, "KRW", "card", "CONFIRMED",
            gross = 10_000, refund = 3_000, fee = 200, net = 6_800, paymentCount = 4, refundCount = 1,
            createdAt = createdAt, updatedAt = updatedAt,
        )

        val result = repo.findOne(uid, context)

        assertEquals(
            DailySettlementResult(
                id = result.id,
                uid = uid,
                settlementDate = day1,
                currency = "KRW",
                paymentMethod = "card",
                grossAmount = 10_000,
                refundAmount = 3_000,
                feeAmount = 200,
                netAmount = 6_800,
                paymentCount = 4,
                refundCount = 1,
                status = "CONFIRMED",
                createdAt = createdAt,
                updatedAt = updatedAt,
            ),
            result,
        )
        assertEquals(idOf(uid), result.id)
    }

    @Test
    fun `findOne_nullableColumnsUnset_mapsNulls`() {
        val uid = seed(day1, "USD", null, "OPEN")

        val result = repo.findOne(uid, context)

        assertNull(result.paymentMethod)
        assertNull(result.updatedAt)
        assertEquals("USD", result.currency)
        assertEquals("OPEN", result.status)
    }

    @Test
    fun `findOne_unknownUid_throwsNoSuchElement`() {
        seed(day1, "KRW", "card", "CONFIRMED")
        val unknown = UUID.randomUUID()

        val error = assertFailsWith<NoSuchElementException> { repo.findOne(unknown, context) }

        assertTrue(error.message!!.contains(unknown.toString()))
    }

    @Test
    fun `findOne_softDeletedBucket_throwsNoSuchElement`() {
        val uid = seed(day1, "KRW", "card", "CONFIRMED", deletedAt = LocalDateTime.of(2026, 9, 4, 0, 0))

        assertFailsWith<NoSuchElementException> { repo.findOne(uid, context) }
    }

    private fun idOf(uid: UUID): Long =
        settlementJdbc.queryForObject("SELECT id FROM daily_settlement WHERE uid = ?", Long::class.java, uid)!!

    private fun seed(
        date: LocalDate,
        currency: String,
        method: String?,
        status: String,
        gross: Long = 1_000,
        refund: Long = 0,
        fee: Long = 0,
        net: Long = gross - refund - fee,
        paymentCount: Int = 1,
        refundCount: Int = 0,
        createdAt: LocalDateTime = LocalDateTime.of(2026, 9, 4, 2, 0),
        updatedAt: LocalDateTime? = null,
        deletedAt: LocalDateTime? = null,
    ): UUID {
        val uid = UUID.randomUUID()
        settlementJdbc.update(
            """
            INSERT INTO daily_settlement
                (uid, settlement_date, currency, payment_method, gross_amount, refund_amount,
                 fee_amount, net_amount, payment_count, refund_count, status, created_at, updated_at, deleted_at)
            VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
            """.trimIndent(),
            uid, date, currency, method, gross, refund, fee, net, paymentCount, refundCount, status,
            createdAt, updatedAt, deletedAt,
        )
        return uid
    }
}
