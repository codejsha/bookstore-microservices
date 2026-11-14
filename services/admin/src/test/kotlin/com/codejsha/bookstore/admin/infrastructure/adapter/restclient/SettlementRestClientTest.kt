package com.codejsha.bookstore.admin.infrastructure.adapter.restclient

import com.codejsha.bookstore.admin.domain.model.UnsupportedValueException
import com.codejsha.bookstore.admin.domain.model.command.TriggerSettlementRunCommand
import com.codejsha.bookstore.admin.domain.model.external.SettlementBucket
import com.codejsha.bookstore.admin.domain.model.external.SettlementDetailLine
import com.codejsha.bookstore.admin.domain.model.external.SettlementRunAck
import com.codejsha.bookstore.admin.domain.model.option.SettlementQueryOption
import org.junit.jupiter.api.AfterAll
import org.junit.jupiter.api.AfterEach
import org.junit.jupiter.api.Test
import org.junit.jupiter.api.TestInstance
import org.springframework.data.domain.Pageable
import org.springframework.http.HttpStatus
import org.springframework.web.client.HttpClientErrorException
import org.springframework.web.client.HttpServerErrorException
import org.springframework.web.context.request.RequestContextHolder
import java.time.LocalDate
import kotlin.test.assertEquals
import kotlin.test.assertFailsWith

@TestInstance(TestInstance.Lifecycle.PER_CLASS)
class SettlementRestClientTest {

    private val downstream = DownstreamStub()

    private val client = SettlementRestClient(downstream.config.settlementSettlementApi())

    @AfterEach
    fun clearRequest() {
        RequestContextHolder.resetRequestAttributes()
    }

    @AfterAll
    fun stopDownstream() {
        downstream.stop()
    }

    @Test
    fun `findAllSettlements_dateAndStatusWithSortedPage_forwardsIsoDateAndMapsBuckets`() {
        downstream.respond("""{"total":96,"items":[$FULL_BUCKET_JSON,$BARE_BUCKET_JSON]}""")

        val page = client.findAllSettlements(
            SettlementQueryOption(settlementDate = LocalDate.of(2026, 10, 7), status = "DISCREPANCY"),
            SORTED_PAGE,
        )

        val request = downstream.onlyRequest()
        assertEquals("GET", request.method)
        assertEquals("/api/v1/settlements", request.path)
        assertEquals(mapOf("date" to "2026-10-07", "status" to "DISCREPANCY") + SORTED_PAGE_QUERY, request.query)
        assertEquals(96, page.totalElements)
        assertEquals(SORTED_PAGE, page.pageable)
        assertEquals(listOf(FULL_BUCKET, BARE_BUCKET), page.content)
    }

    @Test
    fun `findAllSettlements_noFilterUnpaged_sendsNoQueryParameters`() {
        downstream.respond("""{"total":0,"items":[]}""")

        val page = client.findAllSettlements(SettlementQueryOption(), Pageable.unpaged())

        assertEquals(emptyMap(), downstream.onlyRequest().query)
        assertEquals(emptyList(), page.content)
    }

    @Test
    fun `findAllSettlements_unknownStatus_failsWithoutCallingDownstream`() {
        downstream.respond("""{"total":0,"items":[]}""")

        val failure = assertFailsWith<UnsupportedValueException> {
            client.findAllSettlements(SettlementQueryOption(status = "CLOSED"), SORTED_PAGE)
        }

        assertEquals("unknown settlement status: CLOSED", failure.message)
        assertEquals(0, downstream.requestCount())
    }

    @Test
    fun `findSettlement_bucketWithDetails_mapsBucketAndDetailLines`() {
        downstream.respond(FULL_BUCKET_JSON)

        val bucket = client.findSettlement("settlement-1")

        val request = downstream.onlyRequest()
        assertEquals("GET", request.method)
        assertEquals("/api/v1/settlements/settlement-1", request.path)
        assertEquals(FULL_BUCKET, bucket)
    }

    @Test
    fun `findSettlement_detailsAndOptionalFieldsAbsent_mapsEmptyDetailsAndNulls`() {
        downstream.respond(BARE_BUCKET_JSON)

        assertEquals(BARE_BUCKET, client.findSettlement("settlement-2"))
    }

    @Test
    fun `findSettlement_unknownUid_propagatesNotFound`() {
        downstream.respond(PROBLEM_NOT_FOUND, status = 404)

        val failure = assertFailsWith<HttpClientErrorException.NotFound> { client.findSettlement("missing") }

        assertEquals(PROBLEM_NOT_FOUND, failure.responseBodyAsString)
    }

    @Test
    fun `findSettlement_inboundBearerToken_replaysItDownstream`() {
        bindInboundAuthorization(CALLER_AUTHORIZATION)
        downstream.respond(FULL_BUCKET_JSON)

        client.findSettlement("settlement-1")

        assertEquals(CALLER_AUTHORIZATION, downstream.onlyRequest().authorization)
    }

    @Test
    fun `triggerSettlementRun_rerunForDate_postsIsoTargetDateAndMapsAck`() {
        downstream.respond(
            """{"uid":"run-1","target_date":"2026-10-07","status":"STARTED","started_at":"2026-10-08T02:00:00Z"}""",
            status = 202,
        )

        val ack = client.triggerSettlementRun(TriggerSettlementRunCommand(LocalDate.of(2026, 10, 7), rerun = true))

        val request = downstream.onlyRequest()
        assertEquals("POST", request.method)
        assertEquals("/api/v1/settlement-runs", request.path)
        assertEquals(mapOf("target_date" to "2026-10-07", "rerun" to true), request.json())
        assertEquals(
            SettlementRunAck(
                uid = "run-1",
                targetDate = LocalDate.of(2026, 10, 7),
                status = "STARTED",
                startedAt = utc(2026, 10, 8, 2, 0),
            ),
            ack,
        )
    }

    @Test
    fun `triggerSettlementRun_firstRun_postsRerunFalse`() {
        downstream.respond(
            """{"uid":"run-2","target_date":"2026-10-08","status":"STARTED","started_at":"2026-10-09T02:00:00Z"}""",
            status = 202,
        )

        client.triggerSettlementRun(TriggerSettlementRunCommand(LocalDate.of(2026, 10, 8), rerun = false))

        assertEquals(mapOf("target_date" to "2026-10-08", "rerun" to false), downstream.onlyRequest().json())
    }

    @Test
    fun `triggerSettlementRun_runAlreadyInProgress_propagatesConflict`() {
        downstream.respond("""{"title":"Conflict","status":409}""", status = 409)

        assertFailsWith<HttpClientErrorException.Conflict> {
            client.triggerSettlementRun(TriggerSettlementRunCommand(LocalDate.of(2026, 10, 7), rerun = false))
        }
    }

    @Test
    fun `triggerSettlementRun_downstreamFailure_propagatesServerError`() {
        downstream.respond(status = 500)

        val failure = assertFailsWith<HttpServerErrorException> {
            client.triggerSettlementRun(TriggerSettlementRunCommand(LocalDate.of(2026, 10, 7), rerun = true))
        }

        assertEquals(HttpStatus.INTERNAL_SERVER_ERROR, failure.statusCode)
    }

    private companion object {
        private const val FULL_BUCKET_JSON =
            """{
                "uid":"settlement-1",
                "settlement_date":"2026-10-07",
                "currency":"KRW",
                "payment_method":"card",
                "gross_amount":100000,
                "refund_amount":5000,
                "fee_amount":2900,
                "net_amount":92100,
                "payment_count":3,
                "refund_count":1,
                "status":"DISCREPANCY",
                "details":[
                    {
                        "uid":"detail-1",
                        "source_type":"PAYMENT",
                        "source_id":"payment-1",
                        "payment_id":"pay_1",
                        "amount":100000,
                        "currency":"KRW",
                        "payment_method":"card",
                        "occurred_at":"2026-10-07T01:02:00Z"
                    },
                    {
                        "uid":"detail-2",
                        "source_type":"REFUND",
                        "source_id":"refund-1",
                        "payment_id":"pay_1",
                        "amount":5000,
                        "currency":"KRW",
                        "occurred_at":"2026-10-07T03:04:00Z"
                    }
                ],
                "created_at":"2026-10-08T02:00:00Z",
                "updated_at":"2026-10-08T02:01:00Z"
            }"""

        private const val BARE_BUCKET_JSON =
            """{
                "uid":"settlement-2",
                "settlement_date":"2026-10-06",
                "currency":"USD",
                "gross_amount":0,
                "refund_amount":0,
                "fee_amount":0,
                "net_amount":0,
                "payment_count":0,
                "refund_count":0,
                "status":"OPEN",
                "created_at":"2026-10-07T02:00:00Z"
            }"""

        private val FULL_BUCKET = SettlementBucket(
            uid = "settlement-1",
            settlementDate = LocalDate.of(2026, 10, 7),
            currency = "KRW",
            paymentMethod = "card",
            grossAmount = 100000,
            refundAmount = 5000,
            feeAmount = 2900,
            netAmount = 92100,
            paymentCount = 3,
            refundCount = 1,
            status = "DISCREPANCY",
            details = listOf(
                SettlementDetailLine(
                    uid = "detail-1",
                    sourceType = "PAYMENT",
                    sourceId = "payment-1",
                    paymentId = "pay_1",
                    amount = 100000,
                    currency = "KRW",
                    paymentMethod = "card",
                    occurredAt = utc(2026, 10, 7, 1, 2),
                ),
                SettlementDetailLine(
                    uid = "detail-2",
                    sourceType = "REFUND",
                    sourceId = "refund-1",
                    paymentId = "pay_1",
                    amount = 5000,
                    currency = "KRW",
                    paymentMethod = null,
                    occurredAt = utc(2026, 10, 7, 3, 4),
                ),
            ),
            createdAt = utc(2026, 10, 8, 2, 0),
            updatedAt = utc(2026, 10, 8, 2, 1),
        )

        private val BARE_BUCKET = SettlementBucket(
            uid = "settlement-2",
            settlementDate = LocalDate.of(2026, 10, 6),
            currency = "USD",
            paymentMethod = null,
            grossAmount = 0,
            refundAmount = 0,
            feeAmount = 0,
            netAmount = 0,
            paymentCount = 0,
            refundCount = 0,
            status = "OPEN",
            details = emptyList(),
            createdAt = utc(2026, 10, 7, 2, 0),
            updatedAt = null,
        )
    }
}
