package com.codejsha.bookstore.settlement.infrastructure.adapter.restcontroller

import com.codejsha.bookstore.generated.application.port.openapi.model.SettlementDetailResponse
import com.codejsha.bookstore.generated.application.port.openapi.model.SettlementFindResponse
import com.codejsha.bookstore.generated.application.port.openapi.model.SettlementRunResponse
import com.codejsha.bookstore.settlement.domain.aggregate.DailySettlement
import com.codejsha.bookstore.settlement.domain.aggregate.SettlementDetailEntity
import com.codejsha.bookstore.settlement.domain.aggregate.SettlementRun
import com.codejsha.bookstore.settlement.domain.constant.SettlementSourceType
import com.codejsha.bookstore.settlement.domain.constant.SettlementStatus
import org.junit.jupiter.api.Test
import java.time.LocalDate
import java.time.LocalDateTime
import java.time.OffsetDateTime
import java.time.ZoneOffset
import java.util.UUID
import kotlin.test.assertEquals
import kotlin.test.assertNull
import com.codejsha.bookstore.generated.application.port.openapi.model.SettlementSourceType as ApiSourceType
import com.codejsha.bookstore.generated.application.port.openapi.model.SettlementStatus as ApiStatus

class SettlementResponseMappingTest {

    @Test
    fun `toSettlementFindResponse_bucketWithDetails_mapsEveryField`() {
        val response = toSettlementFindResponse(
            settlement(details = listOf(detail(DETAIL_UID, SettlementSourceType.REFUND, paymentMethod = "card"))),
        )

        assertEquals(
            SettlementFindResponse(
                uid = SETTLEMENT_UID.toString(),
                settlementDate = DATE,
                currency = "KRW",
                paymentMethod = "card",
                grossAmount = 10_000,
                refundAmount = 3_000,
                feeAmount = 200,
                netAmount = 6_800,
                paymentCount = 4,
                refundCount = 1,
                status = ApiStatus.CONFIRMED,
                details = listOf(
                    SettlementDetailResponse(
                        uid = DETAIL_UID.toString(),
                        sourceType = ApiSourceType.REFUND,
                        sourceId = "ref_1",
                        paymentId = "pay_1",
                        amount = 3_000,
                        currency = "KRW",
                        paymentMethod = "card",
                        occurredAt = OffsetDateTime.of(2026, 7, 1, 3, 4, 5, 123_456_000, ZoneOffset.UTC),
                    ),
                ),
                createdAt = OffsetDateTime.of(2026, 7, 2, 2, 0, 0, 0, ZoneOffset.UTC),
                updatedAt = OffsetDateTime.of(2026, 7, 2, 2, 5, 0, 0, ZoneOffset.UTC),
            ),
            response,
        )
    }

    @Test
    fun `toSettlementFindResponse_nullableFieldsUnset_mapsNulls`() {
        val response = toSettlementFindResponse(settlement(paymentMethod = null, updatedAt = null))

        assertNull(response.paymentMethod)
        assertNull(response.updatedAt)
    }

    @Test
    fun `toSettlementFindResponse_bucketWithoutDetails_omitsTheDetailList`() {
        assertNull(toSettlementFindResponse(settlement(details = emptyList())).details)
    }

    @Test
    fun `toSettlementFindResponse_everyDomainStatus_mapsToTheSameApiStatus`() {
        SettlementStatus.entries.forEach { status ->
            val response = toSettlementFindResponse(settlement(status = status))

            assertEquals(status.value, response.status.value)
        }
        assertEquals(SettlementStatus.entries.size, ApiStatus.entries.size)
    }

    @Test
    fun `toSettlementFindResponse_localTimestamps_areStampedAsUtcWithoutShifting`() {
        val response = toSettlementFindResponse(settlement())

        assertEquals(ZoneOffset.UTC, response.createdAt.offset)
        assertEquals(CREATED_AT, response.createdAt.toLocalDateTime())
        assertEquals(UPDATED_AT, response.updatedAt?.toLocalDateTime())
    }

    @Test
    fun `toSettlementDetailResponse_everySourceType_mapsToTheSameApiSourceType`() {
        SettlementSourceType.entries.forEach { sourceType ->
            val response = toSettlementDetailResponse(detail(DETAIL_UID, sourceType, paymentMethod = "card"))

            assertEquals(sourceType.value, response.sourceType.value)
        }
        assertEquals(SettlementSourceType.entries.size, ApiSourceType.entries.size)
    }

    @Test
    fun `toSettlementDetailResponse_lineWithoutPaymentMethod_mapsNull`() {
        val response = toSettlementDetailResponse(detail(DETAIL_UID, SettlementSourceType.PAYMENT, paymentMethod = null))

        assertNull(response.paymentMethod)
        assertEquals(ApiSourceType.PAYMENT, response.sourceType)
        assertEquals(OCCURRED_AT, response.occurredAt.toLocalDateTime())
        assertEquals(ZoneOffset.UTC, response.occurredAt.offset)
    }

    @Test
    fun `toSettlementRunResponse_launchedRun_mapsEveryField`() {
        val runUid = UUID.fromString("33333333-3333-3333-3333-333333333333")

        val response = toSettlementRunResponse(
            SettlementRun(uid = runUid, targetDate = DATE, status = "STARTING", startedAt = CREATED_AT),
        )

        assertEquals(
            SettlementRunResponse(
                uid = runUid.toString(),
                targetDate = DATE,
                status = "STARTING",
                startedAt = OffsetDateTime.of(2026, 7, 2, 2, 0, 0, 0, ZoneOffset.UTC),
            ),
            response,
        )
    }

    private fun settlement(
        paymentMethod: String? = "card",
        status: SettlementStatus = SettlementStatus.CONFIRMED,
        details: List<SettlementDetailEntity> = emptyList(),
        updatedAt: LocalDateTime? = UPDATED_AT,
    ) = DailySettlement(
        id = 42L,
        uid = SETTLEMENT_UID,
        settlementDate = DATE,
        currency = "KRW",
        paymentMethod = paymentMethod,
        grossAmount = 10_000,
        refundAmount = 3_000,
        feeAmount = 200,
        netAmount = 6_800,
        paymentCount = 4,
        refundCount = 1,
        status = status,
        details = details,
        createdAt = CREATED_AT,
        updatedAt = updatedAt,
    )

    private fun detail(uid: UUID, sourceType: SettlementSourceType, paymentMethod: String?) = SettlementDetailEntity(
        id = 7L,
        uid = uid,
        settlementDate = DATE,
        sourceType = sourceType,
        sourceId = "ref_1",
        paymentId = "pay_1",
        amount = 3_000,
        currency = "KRW",
        paymentMethod = paymentMethod,
        occurredAt = OCCURRED_AT,
        createdAt = CREATED_AT,
        updatedAt = null,
    )

    private companion object {
        val SETTLEMENT_UID: UUID = UUID.fromString("22222222-2222-2222-2222-222222222222")
        val DETAIL_UID: UUID = UUID.fromString("44444444-4444-4444-4444-444444444444")
        val DATE: LocalDate = LocalDate.of(2026, 7, 1)
        val OCCURRED_AT: LocalDateTime = LocalDateTime.of(2026, 7, 1, 3, 4, 5, 123_456_000)
        val CREATED_AT: LocalDateTime = LocalDateTime.of(2026, 7, 2, 2, 0, 0)
        val UPDATED_AT: LocalDateTime = LocalDateTime.of(2026, 7, 2, 2, 5, 0)
    }
}
