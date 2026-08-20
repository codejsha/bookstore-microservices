package com.codejsha.bookstore.settlement.domain.service

import com.codejsha.bookstore.settlement.application.port.repo.DailySettlementResult
import com.codejsha.bookstore.settlement.application.port.repo.SettlementDetailResult
import com.codejsha.bookstore.settlement.domain.constant.SettlementSourceType
import com.codejsha.bookstore.settlement.domain.constant.SettlementStatus
import com.codejsha.bookstore.settlement.domain.model.option.SettlementQueryOption
import com.codejsha.bookstore.settlement.support.FakeDailySettlementRepo
import com.codejsha.bookstore.settlement.support.FakeSettlementDetailRepo
import com.codejsha.bookstore.settlement.support.FakeTransactionRunner
import com.codejsha.platform.shared.data.ActorContext
import com.codejsha.platform.shared.data.ActorType
import kotlinx.coroutines.runBlocking
import org.springframework.data.domain.Pageable
import java.time.LocalDate
import java.time.LocalDateTime
import java.util.UUID
import kotlin.test.Test
import kotlin.test.assertEquals

class SettlementServiceTest {

    private val context = ActorContext(actorId = 0L, actorType = ActorType.USER)
    private val date = LocalDate.of(2026, 7, 14)

    private fun settlementRow(uid: UUID) = DailySettlementResult(
        id = 1L,
        uid = uid,
        settlementDate = date,
        currency = "KRW",
        paymentMethod = "card",
        grossAmount = 10_000L,
        refundAmount = 1_000L,
        feeAmount = 300L,
        netAmount = 8_700L,
        paymentCount = 3,
        refundCount = 1,
        status = SettlementStatus.OPEN.value,
        createdAt = LocalDateTime.of(2026, 7, 15, 1, 0),
        updatedAt = null,
    )

    private fun detailRow(uid: UUID) = SettlementDetailResult(
        id = 1L,
        uid = uid,
        settlementDate = date,
        sourceType = SettlementSourceType.PAYMENT.value,
        sourceId = "pay_abc",
        paymentId = "pay_abc",
        amount = 10_000L,
        currency = "KRW",
        paymentMethod = "card",
        occurredAt = LocalDateTime.of(2026, 7, 14, 12, 0),
        createdAt = LocalDateTime.of(2026, 7, 15, 1, 0),
        updatedAt = null,
    )

    @Test
    fun `findSettlement composes matching bucket details onto the aggregate`(): Unit = runBlocking {
        val settlementUid = UUID.randomUUID()
        val service = SettlementService(
            dailySettlementRepo = FakeDailySettlementRepo(listOf(settlementRow(settlementUid))),
            settlementDetailRepo = FakeSettlementDetailRepo(listOf(detailRow(UUID.randomUUID()))),
            txRunner = FakeTransactionRunner(),
        )

        val result = service.findSettlement(settlementUid, context)

        assertEquals(settlementUid, result.uid)
        assertEquals(SettlementStatus.OPEN, result.status)
        assertEquals(1, result.details.size)
        assertEquals(SettlementSourceType.PAYMENT, result.details.first().sourceType)
    }

    @Test
    fun `findAllSettlements maps rows to aggregates and honours the status filter`(): Unit = runBlocking {
        val service = SettlementService(
            dailySettlementRepo = FakeDailySettlementRepo(listOf(settlementRow(UUID.randomUUID()))),
            settlementDetailRepo = FakeSettlementDetailRepo(emptyList()),
            txRunner = FakeTransactionRunner(),
        )

        val open = service.findAllSettlements(
            SettlementQueryOption(status = SettlementStatus.OPEN.value),
            Pageable.unpaged(),
            context,
        )
        val confirmed = service.findAllSettlements(
            SettlementQueryOption(status = SettlementStatus.CONFIRMED.value),
            Pageable.unpaged(),
            context,
        )

        assertEquals(1, open.totalElements)
        assertEquals(0, confirmed.totalElements)
    }
}
