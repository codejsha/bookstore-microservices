package com.codejsha.bookstore.payment.infrastructure.adapter.pgsql

import com.codejsha.bookstore.generated.infrastructure.adapter.jooq.tables.references.MANDATES
import com.codejsha.bookstore.payment.domain.model.command.MandateCreateCommand
import com.codejsha.bookstore.payment.domain.model.option.MandateQueryOption
import com.codejsha.bookstore.payment.infrastructure.support.AbstractPostgresTestSupport
import com.codejsha.bookstore.payment.infrastructure.support.PersistenceManager
import com.codejsha.bookstore.payment.support.PaymentTestFixtures
import org.jooq.DSLContext
import org.junit.jupiter.api.BeforeEach
import org.junit.jupiter.api.Test
import org.springframework.beans.factory.annotation.Autowired
import org.springframework.boot.jooq.test.autoconfigure.JooqTest
import org.springframework.context.annotation.Import
import org.springframework.dao.DataIntegrityViolationException
import org.springframework.dao.DuplicateKeyException
import org.springframework.data.domain.PageRequest
import org.springframework.data.domain.Sort
import org.springframework.transaction.annotation.Propagation
import org.springframework.transaction.annotation.Transactional
import java.time.LocalDateTime
import java.time.ZoneOffset
import java.util.UUID
import kotlin.test.assertEquals
import kotlin.test.assertFailsWith
import kotlin.test.assertIs
import kotlin.test.assertNotNull
import kotlin.test.assertNull
import kotlin.test.assertTrue

@JooqTest(
    properties = [
        "spring.cloud.config.enabled=false",
        "spring.cloud.config.import-check.enabled=false",
    ],
)
@Import(PersistenceManager::class, MandateRepoImpl::class)
@Transactional(propagation = Propagation.NOT_SUPPORTED)
class MandateRepoImplTest : AbstractPostgresTestSupport() {

    @Autowired
    private lateinit var dsl: DSLContext

    @Autowired
    private lateinit var mandateRepo: MandateRepoImpl

    private val ctx = PaymentTestFixtures.DEFAULT_CONTEXT

    @BeforeEach
    fun cleanTables() {
        dsl.deleteFrom(MANDATES).execute()
    }

    @Test
    fun `create_fullCommand_persistsEveryColumnAndStampsAcceptance`() {
        val created = mandateRepo.create(
            createCommand(mandateId = "man_full", metadata = mapOf("source" to "checkout")),
            ctx,
        )

        val found = mandateRepo.findOne(created.uid, ctx)
        assertTrue(found.id > 0L)
        assertEquals(created.uid, found.uid)
        assertEquals("man_full", found.mandateId)
        assertEquals("cus_001", found.customerId)
        assertEquals("pm_001", found.paymentMethodId)
        assertEquals("multi_use", found.mandateType)
        assertEquals("active", found.mandateStatus)
        assertEquals(100_000L, found.mandateAmount)
        assertEquals("KRW", found.mandateCurrency)
        assertEquals("off_session", found.setupFutureUsage)
        assertEquals("online", found.customerAcceptanceType)
        assertEquals("checkout", found.metadata?.get("source"))
        assertNotNull(found.customerAcceptedAt)
        assertNotNull(found.createdAt)
        assertNull(found.startDate)
        assertNull(found.endDate)
        assertNull(found.updatedAt)
    }

    @Test
    fun `create_nullOptionalFields_storesNulls`() {
        val created = mandateRepo.create(
            MandateCreateCommand(
                mandateId = "man_min",
                customerId = "cus_001",
                paymentMethodId = null,
                mandateType = "single_use",
                mandateStatus = "pending",
                mandateAmount = null,
                mandateCurrency = null,
                setupFutureUsage = null,
                customerAcceptanceType = null,
            ),
            ctx,
        )

        assertEquals("pending", created.mandateStatus)
        assertNull(created.paymentMethodId)
        assertNull(created.mandateAmount)
        assertNull(created.mandateCurrency)
        assertNull(created.setupFutureUsage)
        assertNull(created.customerAcceptanceType)
        assertNull(created.metadata)
    }

    @Test
    fun `create_duplicateMandateId_throwsDuplicateKeyException`() {
        mandateRepo.create(createCommand(mandateId = "man_dup"), ctx)

        val ex = assertFailsWith<DataIntegrityViolationException> {
            mandateRepo.create(createCommand(mandateId = "man_dup"), ctx)
        }

        assertIs<DuplicateKeyException>(ex)
        assertEquals(1, dsl.fetchCount(MANDATES))
    }

    @Test
    fun `findOne_softDeletedMandate_throwsNoSuchElementException`() {
        val created = mandateRepo.create(createCommand(mandateId = "man_gone"), ctx)
        softDelete(created.uid)

        assertFailsWith<NoSuchElementException> { mandateRepo.findOne(created.uid, ctx) }
    }

    @Test
    fun `findActiveByCustomerId_severalMandates_picksNewestActiveOfThatCustomer`() {
        val older = mandateRepo.create(createCommand(mandateId = "man_old"), ctx)
        val newer = mandateRepo.create(createCommand(mandateId = "man_new"), ctx)
        val pending = mandateRepo.create(createCommand(mandateId = "man_pending", mandateStatus = "pending"), ctx)
        mandateRepo.create(createCommand(mandateId = "man_other", customerId = "cus_002"), ctx)
        setCreatedAt(older.uid, BASE_TIME)
        setCreatedAt(newer.uid, BASE_TIME.plusHours(1))
        setCreatedAt(pending.uid, BASE_TIME.plusHours(2))

        val found = mandateRepo.findActiveByCustomerId("cus_001", ctx)

        assertEquals("man_new", assertNotNull(found).mandateId)
    }

    @Test
    fun `findActiveByCustomerId_noActiveMandate_yieldsNull`() {
        mandateRepo.create(createCommand(mandateId = "man_pending", mandateStatus = "pending"), ctx)
        val deleted = mandateRepo.create(createCommand(mandateId = "man_deleted"), ctx)
        softDelete(deleted.uid)

        assertNull(mandateRepo.findActiveByCustomerId("cus_001", ctx))
        assertNull(mandateRepo.findActiveByCustomerId("cus_unknown", ctx))
    }

    @Test
    fun `findAll_customerAndStatusFilters_matchOnlyIntersection`() {
        mandateRepo.create(createCommand(mandateId = "man_a"), ctx)
        mandateRepo.create(createCommand(mandateId = "man_b", mandateStatus = "pending"), ctx)
        mandateRepo.create(createCommand(mandateId = "man_c", customerId = "cus_002"), ctx)

        val byCustomer = mandateRepo.findAll(MandateQueryOption(customerId = "cus_001"), PageRequest.of(0, 10), ctx)
        val byBoth = mandateRepo.findAll(
            MandateQueryOption(customerId = "cus_001", mandateStatus = "active"),
            PageRequest.of(0, 10),
            ctx,
        )

        assertEquals(setOf("man_a", "man_b"), byCustomer.content.map { it.mandateId }.toSet())
        assertEquals(2L, byCustomer.totalElements)
        assertEquals(listOf("man_a"), byBoth.content.map { it.mandateId })
        assertEquals(1L, byBoth.totalElements)
    }

    @Test
    fun `findAll_pagedSortedByMandateAmountAsc_slicesInOrderAndSkipsDeleted`() {
        mandateRepo.create(createCommand(mandateId = "man_30", mandateAmount = 30_000L), ctx)
        mandateRepo.create(createCommand(mandateId = "man_10", mandateAmount = 10_000L), ctx)
        mandateRepo.create(createCommand(mandateId = "man_20", mandateAmount = 20_000L), ctx)
        val deleted = mandateRepo.create(createCommand(mandateId = "man_05", mandateAmount = 5_000L), ctx)
        softDelete(deleted.uid)
        val sort = Sort.by(Sort.Direction.ASC, "mandate_amount")

        val first = mandateRepo.findAll(MandateQueryOption(), PageRequest.of(0, 2, sort), ctx)
        val second = mandateRepo.findAll(MandateQueryOption(), PageRequest.of(1, 2, sort), ctx)

        assertEquals(3L, first.totalElements)
        assertEquals(listOf("man_10", "man_20"), first.content.map { it.mandateId })
        assertEquals(listOf("man_30"), second.content.map { it.mandateId })
    }

    @Test
    fun `revoke_activeMandate_marksRevokedAndStampsUpdatedAt`() {
        val created = mandateRepo.create(createCommand(mandateId = "man_rev"), ctx)

        val revoked = mandateRepo.revoke(created.uid, ctx)

        assertEquals("revoked", revoked.mandateStatus)
        assertNotNull(revoked.updatedAt)
        assertNull(mandateRepo.findActiveByCustomerId("cus_001", ctx))
    }

    @Test
    fun `revoke_alreadyRevokedMandate_throwsIllegalStateException`() {
        val created = mandateRepo.create(createCommand(mandateId = "man_rev_twice"), ctx)
        mandateRepo.revoke(created.uid, ctx)

        assertFailsWith<IllegalStateException> { mandateRepo.revoke(created.uid, ctx) }
    }

    @Test
    fun `revoke_pendingMandate_throwsIllegalStateExceptionAndKeepsStatus`() {
        val created = mandateRepo.create(createCommand(mandateId = "man_pending", mandateStatus = "pending"), ctx)

        assertFailsWith<IllegalStateException> { mandateRepo.revoke(created.uid, ctx) }

        assertEquals("pending", mandateRepo.findOne(created.uid, ctx).mandateStatus)
    }

    private fun softDelete(uid: UUID) {
        dsl.update(MANDATES)
            .set(MANDATES.DELETED_AT, LocalDateTime.now(ZoneOffset.UTC))
            .where(MANDATES.UID.eq(uid))
            .execute()
    }

    private fun setCreatedAt(uid: UUID, createdAt: LocalDateTime) {
        dsl.update(MANDATES).set(MANDATES.CREATED_AT, createdAt).where(MANDATES.UID.eq(uid)).execute()
    }

    private fun createCommand(
        mandateId: String,
        customerId: String = "cus_001",
        mandateStatus: String = "active",
        mandateAmount: Long? = 100_000L,
        metadata: Map<String, Any>? = null,
    ) = MandateCreateCommand(
        mandateId = mandateId,
        customerId = customerId,
        paymentMethodId = "pm_001",
        mandateType = "multi_use",
        mandateStatus = mandateStatus,
        mandateAmount = mandateAmount,
        mandateCurrency = "KRW",
        setupFutureUsage = "off_session",
        customerAcceptanceType = "online",
        metadata = metadata,
    )

    private companion object {
        val BASE_TIME: LocalDateTime = LocalDateTime.of(2026, 1, 1, 0, 0, 0)
    }
}
