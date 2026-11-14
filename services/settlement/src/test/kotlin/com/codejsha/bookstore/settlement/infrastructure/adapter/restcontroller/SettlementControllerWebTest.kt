package com.codejsha.bookstore.settlement.infrastructure.adapter.restcontroller

import com.codejsha.bookstore.settlement.application.SettlementRunConflictException
import com.codejsha.bookstore.settlement.application.SettlementRunValidationException
import com.codejsha.bookstore.settlement.application.usecase.SettlementUseCase
import com.codejsha.bookstore.settlement.application.usecase.TriggerSettlementRunUseCase
import com.codejsha.bookstore.settlement.domain.aggregate.DailySettlement
import com.codejsha.bookstore.settlement.domain.aggregate.SettlementDetailEntity
import com.codejsha.bookstore.settlement.domain.aggregate.SettlementRun
import com.codejsha.bookstore.settlement.domain.constant.SettlementSourceType
import com.codejsha.bookstore.settlement.domain.constant.SettlementStatus
import com.codejsha.bookstore.settlement.domain.model.command.TriggerSettlementRunCommand
import com.codejsha.bookstore.settlement.domain.model.option.SettlementQueryOption
import com.codejsha.bookstore.settlement.infrastructure.support.auth.HttpPrincipalResolver
import com.codejsha.platform.shared.data.ActorContext
import org.junit.jupiter.api.AfterEach
import org.junit.jupiter.api.Test
import org.springframework.data.domain.Page
import org.springframework.data.domain.PageImpl
import org.springframework.data.domain.Pageable
import org.springframework.data.web.PageableHandlerMethodArgumentResolver
import org.springframework.http.MediaType
import org.springframework.mock.web.MockHttpServletRequest
import org.springframework.test.web.servlet.MockMvc
import org.springframework.test.web.servlet.request.MockMvcRequestBuilders.get
import org.springframework.test.web.servlet.request.MockMvcRequestBuilders.post
import org.springframework.test.web.servlet.result.MockMvcResultMatchers.content
import org.springframework.test.web.servlet.result.MockMvcResultMatchers.jsonPath
import org.springframework.test.web.servlet.result.MockMvcResultMatchers.status
import org.springframework.test.web.servlet.setup.MockMvcBuilders
import org.springframework.web.context.request.RequestContextHolder
import org.springframework.web.context.request.ServletRequestAttributes
import tools.jackson.databind.ObjectMapper
import java.time.LocalDate
import java.time.LocalDateTime
import java.util.UUID
import kotlin.test.assertEquals
import kotlin.test.assertNull
import kotlin.test.assertTrue

class SettlementControllerWebTest {

    private class RecordingSettlementUseCase(
        private val rows: List<DailySettlement>,
    ) : SettlementUseCase {
        var option: SettlementQueryOption? = null
        var pageable: Pageable? = null
        var requestedUid: UUID? = null

        override fun findAllSettlements(
            option: SettlementQueryOption, pageable: Pageable, context: ActorContext,
        ): Page<DailySettlement> {
            this.option = option
            this.pageable = pageable
            return PageImpl(rows, pageable, 37L)
        }

        override fun findSettlement(uid: UUID, context: ActorContext): DailySettlement {
            requestedUid = uid
            return rows.firstOrNull { it.uid == uid }
                ?: throw NoSuchElementException("Settlement with uid $uid not found")
        }
    }

    private class RecordingTriggerUseCase(
        private val error: RuntimeException? = null,
    ) : TriggerSettlementRunUseCase {
        var command: TriggerSettlementRunCommand? = null

        override fun triggerSettlementRun(command: TriggerSettlementRunCommand, context: ActorContext): SettlementRun {
            error?.let { throw it }
            this.command = command
            return SettlementRun(uid = RUN_UID, targetDate = command.targetDate, status = "STARTING", startedAt = CREATED_AT)
        }
    }

    private val settlementUseCase = RecordingSettlementUseCase(listOf(settlement()))

    @AfterEach
    fun clearPrincipal() {
        RequestContextHolder.resetRequestAttributes()
    }

    @Test
    fun `settlementsGetAll_dateStatusAndPaging_passesTheFilterAndAnswersTheWireShape`() {
        mockMvc().perform(
            get("/api/v1/settlements")
                .param("date", "2026-07-01").param("status", "CONFIRMED")
                .param("page", "1").param("size", "5")
                .header("X-User-Id", ACTOR_UID).header("X-User-Roles", "STAFF"),
        )
            .andExpect(status().isOk)
            .andExpect(jsonPath("$.total").value(37))
            .andExpect(jsonPath("$.items.length()").value(1))
            .andExpect(jsonPath("$.items[0].uid").value(SETTLEMENT_UID.toString()))
            .andExpect(jsonPath("$.items[0].id").doesNotExist())
            .andExpect(jsonPath("$.items[0].settlement_date").value("2026-07-01"))
            .andExpect(jsonPath("$.items[0].payment_method").value("card"))
            .andExpect(jsonPath("$.items[0].gross_amount").value(10_000))
            .andExpect(jsonPath("$.items[0].refund_amount").value(3_000))
            .andExpect(jsonPath("$.items[0].fee_amount").value(200))
            .andExpect(jsonPath("$.items[0].net_amount").value(6_800))
            .andExpect(jsonPath("$.items[0].payment_count").value(4))
            .andExpect(jsonPath("$.items[0].refund_count").value(1))
            .andExpect(jsonPath("$.items[0].status").value("CONFIRMED"))

        assertEquals(SettlementQueryOption(settlementDate = DATE, status = "CONFIRMED"), settlementUseCase.option)
        assertEquals(1, settlementUseCase.pageable?.pageNumber)
        assertEquals(5, settlementUseCase.pageable?.pageSize)
    }

    @Test
    fun `settlementsGetAll_noFilter_queriesEveryBucket`() {
        mockMvc().perform(
            get("/api/v1/settlements").header("X-User-Id", ACTOR_UID).header("X-User-Roles", "SYSTEM"),
        ).andExpect(status().isOk)

        assertEquals(SettlementQueryOption(), settlementUseCase.option)
    }

    @Test
    fun `settlementsGetAll_noPageable_queriesUnpaged`() {
        bindPrincipal(roles = "STAFF")
        val controller = SettlementController(settlementUseCase, RecordingTriggerUseCase(), resolver)

        val response = controller.settlementsGetAll(null, null, null)

        assertEquals(200, response.statusCode.value())
        assertTrue(settlementUseCase.pageable!!.isUnpaged)
        assertEquals(37L, response.body?.total)
        assertEquals(listOf(SETTLEMENT_UID.toString()), response.body?.items?.map { it.uid })
    }

    @Test
    fun `settlementsGetAll_callerWithoutStaffRole_answersForbidden`() {
        mockMvc().perform(
            get("/api/v1/settlements").header("X-User-Id", ACTOR_UID).header("X-User-Roles", "USER"),
        ).andExpect(status().isForbidden)

        assertNull(settlementUseCase.option)
    }

    @Test
    fun `settlementsGetAll_anonymousCaller_answersUnauthorized`() {
        mockMvc().perform(get("/api/v1/settlements")).andExpect(status().isUnauthorized)

        assertNull(settlementUseCase.option)
    }

    @Test
    fun `settlementsRead_existingUid_answersTheBucketWithItsDetailLines`() {
        mockMvc().perform(
            get("/api/v1/settlements/{uid}", SETTLEMENT_UID)
                .header("X-User-Id", ACTOR_UID).header("X-User-Roles", "STAFF"),
        )
            .andExpect(status().isOk)
            .andExpect(jsonPath("$.uid").value(SETTLEMENT_UID.toString()))
            .andExpect(jsonPath("$.id").doesNotExist())
            .andExpect(jsonPath("$.details.length()").value(1))
            .andExpect(jsonPath("$.details[0].uid").value(DETAIL_UID.toString()))
            .andExpect(jsonPath("$.details[0].id").doesNotExist())
            .andExpect(jsonPath("$.details[0].source_type").value("PAYMENT"))
            .andExpect(jsonPath("$.details[0].source_id").value("pay_1"))
            .andExpect(jsonPath("$.details[0].payment_id").value("pay_1"))
            .andExpect(jsonPath("$.details[0].amount").value(10_000))

        assertEquals(SETTLEMENT_UID, settlementUseCase.requestedUid)
    }

    @Test
    fun `settlementsRead_unknownUid_answersNotFoundProblem`() {
        val unknown = UUID.fromString("99999999-9999-9999-9999-999999999999")

        mockMvc().perform(
            get("/api/v1/settlements/{uid}", unknown).header("X-User-Id", ACTOR_UID).header("X-User-Roles", "STAFF"),
        )
            .andExpect(status().isNotFound)
            .andExpect(content().contentTypeCompatibleWith(MediaType.APPLICATION_PROBLEM_JSON))
            .andExpect(jsonPath("$.status").value(404))
            .andExpect(jsonPath("$.title").value("Not Found"))
            .andExpect(jsonPath("$.detail").value("Settlement with uid $unknown not found"))
    }

    @Test
    fun `settlementsRead_malformedUid_answersBadRequestProblem`() {
        mockMvc().perform(
            get("/api/v1/settlements/{uid}", "not-a-uuid")
                .header("X-User-Id", ACTOR_UID).header("X-User-Roles", "STAFF"),
        )
            .andExpect(status().isBadRequest)
            .andExpect(content().contentTypeCompatibleWith(MediaType.APPLICATION_PROBLEM_JSON))
            .andExpect(jsonPath("$.status").value(400))
            .andExpect(jsonPath("$.title").value("Bad Request"))

        assertNull(settlementUseCase.requestedUid)
    }

    @Test
    fun `settlementRunsCreate_managerRequest_answersAcceptedWithTheRun`() {
        val trigger = RecordingTriggerUseCase()

        mockMvc(trigger).perform(
            post("/api/v1/settlement-runs")
                .contentType(MediaType.APPLICATION_JSON)
                .content("""{"target_date":"2026-07-01","rerun":true}""")
                .header("X-User-Id", ACTOR_UID).header("X-User-Roles", "MANAGER,STAFF"),
        )
            .andExpect(status().isAccepted)
            .andExpect(jsonPath("$.uid").value(RUN_UID.toString()))
            .andExpect(jsonPath("$.target_date").value("2026-07-01"))
            .andExpect(jsonPath("$.status").value("STARTING"))

        assertEquals(TriggerSettlementRunCommand(targetDate = DATE, rerun = true), trigger.command)
    }

    @Test
    fun `settlementRunsCreate_rerunOmitted_defaultsToNoRerun`() {
        val trigger = RecordingTriggerUseCase()

        mockMvc(trigger).perform(
            post("/api/v1/settlement-runs")
                .contentType(MediaType.APPLICATION_JSON)
                .content("""{"target_date":"2026-07-01"}""")
                .header("X-User-Id", ACTOR_UID).header("X-User-Roles", "MANAGER"),
        ).andExpect(status().isAccepted)

        assertEquals(TriggerSettlementRunCommand(targetDate = DATE, rerun = false), trigger.command)
    }

    @Test
    fun `settlementRunsCreate_runConflict_answersBadRequestProblem`() {
        val trigger = RecordingTriggerUseCase(SettlementRunConflictException("A settlement run for 2026-07-01 is already running"))

        mockMvc(trigger).perform(
            post("/api/v1/settlement-runs")
                .contentType(MediaType.APPLICATION_JSON)
                .content("""{"target_date":"2026-07-01"}""")
                .header("X-User-Id", ACTOR_UID).header("X-User-Roles", "MANAGER"),
        )
            .andExpect(status().isBadRequest)
            .andExpect(content().contentTypeCompatibleWith(MediaType.APPLICATION_PROBLEM_JSON))
            .andExpect(jsonPath("$.status").value(400))
            .andExpect(jsonPath("$.detail").value("A settlement run for 2026-07-01 is already running"))
    }

    @Test
    fun `settlementRunsCreate_dateNotClosed_answersBadRequestProblem`() {
        val trigger = RecordingTriggerUseCase(SettlementRunValidationException("targetDate must be a closed past day"))

        mockMvc(trigger).perform(
            post("/api/v1/settlement-runs")
                .contentType(MediaType.APPLICATION_JSON)
                .content("""{"target_date":"2026-07-01"}""")
                .header("X-User-Id", ACTOR_UID).header("X-User-Roles", "SYSTEM"),
        )
            .andExpect(status().isBadRequest)
            .andExpect(content().contentTypeCompatibleWith(MediaType.APPLICATION_PROBLEM_JSON))
            .andExpect(jsonPath("$.detail").value("targetDate must be a closed past day"))
    }

    private val resolver = HttpPrincipalResolver(ObjectMapper())

    private fun mockMvc(trigger: TriggerSettlementRunUseCase = RecordingTriggerUseCase()): MockMvc =
        MockMvcBuilders
            .standaloneSetup(SettlementController(settlementUseCase, trigger, resolver))
            .setControllerAdvice(SettlementExceptionHandler())
            .setCustomArgumentResolvers(PageableHandlerMethodArgumentResolver())
            .build()

    private fun bindPrincipal(roles: String) {
        val request = MockHttpServletRequest()
        request.addHeader("X-User-Id", ACTOR_UID)
        request.addHeader("X-User-Roles", roles)
        RequestContextHolder.setRequestAttributes(ServletRequestAttributes(request))
    }

    private fun settlement() = DailySettlement(
        id = 42L,
        uid = SETTLEMENT_UID,
        settlementDate = DATE,
        currency = "KRW",
        paymentMethod = "card",
        grossAmount = 10_000,
        refundAmount = 3_000,
        feeAmount = 200,
        netAmount = 6_800,
        paymentCount = 4,
        refundCount = 1,
        status = SettlementStatus.CONFIRMED,
        details = listOf(
            SettlementDetailEntity(
                id = 7L,
                uid = DETAIL_UID,
                settlementDate = DATE,
                sourceType = SettlementSourceType.PAYMENT,
                sourceId = "pay_1",
                paymentId = "pay_1",
                amount = 10_000,
                currency = "KRW",
                paymentMethod = "card",
                occurredAt = LocalDateTime.of(2026, 7, 1, 3, 0),
                createdAt = CREATED_AT,
                updatedAt = null,
            ),
        ),
        createdAt = CREATED_AT,
        updatedAt = null,
    )

    private companion object {
        const val ACTOR_UID = "11111111-1111-1111-1111-111111111111"
        val SETTLEMENT_UID: UUID = UUID.fromString("22222222-2222-2222-2222-222222222222")
        val RUN_UID: UUID = UUID.fromString("33333333-3333-3333-3333-333333333333")
        val DETAIL_UID: UUID = UUID.fromString("44444444-4444-4444-4444-444444444444")
        val DATE: LocalDate = LocalDate.of(2026, 7, 1)
        val CREATED_AT: LocalDateTime = LocalDateTime.of(2026, 7, 2, 2, 0, 0)
    }
}
