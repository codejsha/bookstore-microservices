package com.codejsha.bookstore.settlement.infrastructure.batch

import com.codejsha.bookstore.settlement.infrastructure.support.utils.uuidToBytes
import org.flywaydb.core.Flyway
import org.junit.jupiter.api.Assertions.assertEquals
import org.junit.jupiter.api.BeforeEach
import org.junit.jupiter.api.Test
import org.springframework.batch.core.BatchStatus
import org.springframework.batch.core.job.parameters.JobParametersBuilder
import org.springframework.batch.core.repository.JobRepository
import org.springframework.batch.infrastructure.item.ExecutionContext
import org.springframework.batch.test.JobLauncherTestUtils
import org.springframework.batch.test.context.SpringBatchTest
import org.springframework.beans.factory.annotation.Autowired
import org.springframework.boot.test.context.SpringBootTest
import org.springframework.jdbc.core.JdbcTemplate
import org.springframework.jdbc.datasource.DriverManagerDataSource
import org.springframework.test.context.DynamicPropertyRegistry
import org.springframework.test.context.DynamicPropertySource
import org.testcontainers.containers.MySQLContainer
import org.testcontainers.junit.jupiter.Testcontainers
import java.time.LocalDate
import java.time.LocalDateTime
import java.time.LocalTime
import java.time.ZoneId
import java.time.ZoneOffset
import java.time.ZonedDateTime
import java.util.UUID
import javax.sql.DataSource

@SpringBootTest(
    properties = [
        "spring.cloud.config.enabled=false",
        "spring.cloud.config.import-check.enabled=false",
        "spring.flyway.enabled=false",
        "spring.batch.job.enabled=false",
        "management.otlp.metrics.export.enabled=false",
    ],
)
@SpringBatchTest
@Testcontainers(disabledWithoutDocker = true)
class SettlementJobIntegrationTest {

    @Autowired
    private lateinit var jobLauncherTestUtils: JobLauncherTestUtils

    @Autowired
    private lateinit var jobRepository: JobRepository

    @Autowired
    private lateinit var staleExecutionRecovery: StaleExecutionRecovery

    @Autowired
    private lateinit var dataSource: DataSource

    private val settlementJdbc: JdbcTemplate by lazy { JdbcTemplate(dataSource) }

    private val kst = ZoneId.of("Asia/Seoul")

    @BeforeEach
    fun cleanTables() {
        listOf("settlement_detail", "daily_settlement", "settlement_job_run")
            .forEach { settlementJdbc.execute("DELETE FROM $it") }
        listOf("payments", "refunds").forEach { paymentJdbc.execute("DELETE FROM $it") }
    }

    // ─── Tests ─────────────────────────────────────────────────────────────────

    @Test
    fun `happy path settles payments and refunds then confirms buckets`() {
        val date = LocalDate.of(2026, 7, 10)
        insertPayment("pay_1", amount = 10_000, captured = 10_000, "KRW", "card", "succeeded", kstAt(date, 12, 0))
        insertPayment("pay_2", amount = 5_000, captured = null, "KRW", "card", "partially_captured", kstAt(date, 13, 0))
        insertPayment("pay_3", amount = 2_000, captured = 2_000, "USD", "wallet", "succeeded", kstAt(date, 14, 0))
        insertPayment("pay_old", amount = 9_999, captured = 9_999, "KRW", "card", "succeeded", kstAt(date.minusDays(1), 12, 0))
        insertRefund("ref_1", "pay_1", amount = 3_000, "KRW", "succeeded", kstAt(date, 15, 0))
        insertRefund("ref_pending", "pay_1", amount = 1_000, "KRW", "pending", kstAt(date, 16, 0))

        assertEquals(BatchStatus.COMPLETED, runJob(date))
        assertEquals(4, detailCount(date))

        val krw = bucket(date, "KRW", "card")
        assertEquals(10_000L, krw.gross)
        assertEquals(3_000L, krw.refund)
        assertEquals(0L, krw.fee)
        assertEquals(7_000L, krw.net)
        assertEquals(2, krw.paymentCount)
        assertEquals(1, krw.refundCount)
        assertEquals("CONFIRMED", krw.status)

        val usd = bucket(date, "USD", "wallet")
        assertEquals(2_000L, usd.gross)
        assertEquals(2_000L, usd.net)
        assertEquals("CONFIRMED", usd.status)

        assertEquals("COMPLETED", jobRunStatus(date))
    }

    @Test
    fun `re-running the same date is idempotent (no duplicate rows)`() {
        val date = LocalDate.of(2026, 7, 11)
        insertPayment("pay_a", amount = 4_000, captured = 4_000, "KRW", "card", "succeeded", kstAt(date, 10, 0))
        insertRefund("ref_a", "pay_a", amount = 1_000, "KRW", "succeeded", kstAt(date, 11, 0))

        assertEquals(BatchStatus.COMPLETED, runJob(date))
        assertEquals(2, detailCount(date))
        assertEquals(1, bucketCount(date))

        assertEquals(BatchStatus.COMPLETED, runJob(date, rerunId = "2"))
        assertEquals(2, detailCount(date))
        assertEquals(1, bucketCount(date))
        assertEquals(3_000L, bucket(date, "KRW", "card").net)
    }

    @Test
    fun `refunds can push a bucket net negative`() {
        val date = LocalDate.of(2026, 7, 12)
        insertPayment("pay_n", amount = 5_000, captured = 5_000, "KRW", "card", "succeeded", kstAt(date, 9, 0))
        insertRefund("ref_n", "pay_n", amount = 8_000, "KRW", "succeeded", kstAt(date, 10, 0))

        assertEquals(BatchStatus.COMPLETED, runJob(date))
        val b = bucket(date, "KRW", "card")
        assertEquals(5_000L, b.gross)
        assertEquals(8_000L, b.refund)
        assertEquals(-3_000L, b.net)
        assertEquals("CONFIRMED", b.status)
    }

    @Test
    fun `internal reconcile mismatch marks discrepancy and fails`() {
        val date = LocalDate.of(2026, 7, 13)
        insertPayment("pay_m", amount = 10_000, captured = 10_000, "KRW", "card", "succeeded", kstAt(date, 9, 0))
        seedDailySettlement(date, "KRW", "card", gross = 99_999, refund = 0, paymentCount = 1, refundCount = 0)

        val execution = jobLauncherTestUtils.launchStep("reconcileStep", jobParameters(date))

        assertEquals(BatchStatus.FAILED, execution.status)
        assertEquals("DISCREPANCY", bucket(date, "KRW", "card").status)
    }

    @Test
    fun `internal reconcile mismatch marks only the mismatched currency`() {
        val date = LocalDate.of(2026, 7, 16)
        seedDetail(date, "PAYMENT", "det_krw", amount = 7_000, "KRW", "card", kstAt(date, 9, 0))
        seedDetail(date, "PAYMENT", "det_usd", amount = 2_000, "USD", "card", kstAt(date, 9, 30))
        seedDailySettlement(date, "KRW", "card", gross = 99_999, refund = 0, paymentCount = 1, refundCount = 0)
        seedDailySettlement(date, "USD", "card", gross = 2_000, refund = 0, paymentCount = 1, refundCount = 0)

        val execution = jobLauncherTestUtils.launchStep("reconcileStep", jobParameters(date))

        assertEquals(BatchStatus.FAILED, execution.status)
        assertEquals("DISCREPANCY", bucket(date, "KRW", "card").status)
        assertEquals("OPEN", bucket(date, "USD", "card").status)
    }

    @Test
    fun `stale STARTED execution is abandoned and no longer counts as live`() {
        val date = LocalDate.of(2026, 7, 17)
        val jobName = "dailySettlementJob"
        createStartedExecution(jobName, date, lastUpdated = LocalDateTime.now().minusMinutes(10))

        assertEquals(true, date in jobRepositoryRunningDates(jobName))
        assertEquals(false, staleExecutionRecovery.hasLiveExecution(jobName, date))

        val abandoned = staleExecutionRecovery.abandonStaleFor(jobName, date)

        assertEquals(1, abandoned)
        assertEquals(false, date in jobRepositoryRunningDates(jobName))
        insertPayment("pay_z", amount = 1_000, captured = 1_000, "KRW", "card", "succeeded", kstAt(date, 10, 0))
        assertEquals(BatchStatus.COMPLETED, runJob(date, rerunId = "recovered"))
    }

    @Test
    fun `fresh STARTED execution stays live and is not abandoned`() {
        val date = LocalDate.of(2026, 7, 18)
        val jobName = "dailySettlementJob"
        val executionId = createStartedExecution(jobName, date, lastUpdated = LocalDateTime.now())

        assertEquals(true, staleExecutionRecovery.hasLiveExecution(jobName, date))
        assertEquals(0, staleExecutionRecovery.abandonStaleFor(jobName, date))

        finishExecution(executionId)
    }

    @Test
    fun `reconcile resets a prior-run DISCREPANCY bucket to OPEN then confirms it`() {
        val date = LocalDate.of(2026, 7, 15)
        seedDetail(date, "PAYMENT", "det_pay", amount = 7_000, "KRW", "card", kstAt(date, 9, 0))
        seedDailySettlement(date, "KRW", "card", gross = 7_000, refund = 0, paymentCount = 1, refundCount = 0)
        settlementJdbc.update(
            "UPDATE daily_settlement SET status = 'DISCREPANCY' WHERE settlement_date = ?", date,
        )

        val execution = jobLauncherTestUtils.launchStep("reconcileStep", jobParameters(date))

        assertEquals(BatchStatus.COMPLETED, execution.status)
        assertEquals("CONFIRMED", bucket(date, "KRW", "card").status)
    }

    @Test
    fun `boundary payment at 23_59_59_999999 KST is included and next-day 00_00_00 is excluded`() {
        val date = LocalDate.of(2026, 7, 14)
        val lastMicroOfDay = ZonedDateTime.of(date, LocalTime.of(23, 59, 59, 999_999_000), kst)
            .toInstant().atZone(ZoneOffset.UTC).toLocalDateTime()
        val firstOfNextDay = ZonedDateTime.of(date.plusDays(1), LocalTime.MIDNIGHT, kst)
            .toInstant().atZone(ZoneOffset.UTC).toLocalDateTime()

        insertPayment("pay_in", amount = 1_000, captured = 1_000, "KRW", "card", "succeeded", lastMicroOfDay)
        insertPayment("pay_out", amount = 2_000, captured = 2_000, "KRW", "card", "succeeded", firstOfNextDay)

        assertEquals(BatchStatus.COMPLETED, runJob(date))
        assertEquals(1, detailCount(date))
        assertEquals(1_000L, bucket(date, "KRW", "card").gross)
    }

    // ─── Helpers ─────────────────────────────────────────────────────────────

    private fun runJob(date: LocalDate, rerunId: String? = null): BatchStatus =
        jobLauncherTestUtils.launchJob(jobParameters(date, rerunId)).status

    private fun jobParameters(date: LocalDate, rerunId: String? = null) =
        JobParametersBuilder()
            .addLocalDate("targetDate", date, true)
            .also { if (rerunId != null) it.addString("rerun.id", rerunId, true) }
            .toJobParameters()

    private fun kstAt(date: LocalDate, hour: Int, minute: Int): LocalDateTime =
        ZonedDateTime.of(date, LocalTime.of(hour, minute), kst)
            .toInstant().atZone(ZoneOffset.UTC).toLocalDateTime()

    private fun jobRepositoryRunningDates(jobName: String): Set<LocalDate> =
        jobRepository.findRunningJobExecutions(jobName)
            .mapNotNull { it.jobParameters.getLocalDate("targetDate") }
            .toSet()

    private fun createStartedExecution(jobName: String, date: LocalDate, lastUpdated: LocalDateTime): Long {
        val parameters = jobParameters(date)
        val instance = jobRepository.createJobInstance(jobName, parameters)
        val execution = jobRepository.createJobExecution(instance, parameters, ExecutionContext())
        settlementJdbc.update(
            "UPDATE BATCH_JOB_EXECUTION SET STATUS = 'STARTED', START_TIME = ?, LAST_UPDATED = ? WHERE JOB_EXECUTION_ID = ?",
            lastUpdated, lastUpdated, execution.id,
        )
        return execution.id
    }

    private fun finishExecution(executionId: Long) {
        val execution = jobRepository.getJobExecution(executionId)!!
        execution.status = BatchStatus.FAILED
        execution.setEndTime(LocalDateTime.now())
        jobRepository.update(execution)
    }

    private fun insertPayment(
        paymentId: String,
        amount: Long,
        captured: Long?,
        currency: String,
        method: String?,
        status: String,
        capturedAtUtc: LocalDateTime,
    ) {
        paymentJdbc.update(
            """
            INSERT INTO payments
                (uid, payment_id, amount, amount_captured, currency, status, payment_method, captured_at, created_at)
            VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
            """.trimIndent(),
            uuidToBytes(UUID.randomUUID()), paymentId, amount, captured, currency, status, method,
            capturedAtUtc, LocalDateTime.now(ZoneOffset.UTC),
        )
    }

    private fun insertRefund(
        refundId: String,
        paymentId: String,
        amount: Long,
        currency: String,
        status: String,
        createdAtUtc: LocalDateTime,
    ) {
        paymentJdbc.update(
            """
            INSERT INTO refunds
                (uid, refund_id, payment_id, amount, currency, status, created_at)
            VALUES (?, ?, ?, ?, ?, ?, ?)
            """.trimIndent(),
            uuidToBytes(UUID.randomUUID()), refundId, paymentId, amount, currency, status,
            createdAtUtc,
        )
    }

    private fun seedDailySettlement(
        date: LocalDate,
        currency: String,
        method: String?,
        gross: Long,
        refund: Long,
        paymentCount: Int,
        refundCount: Int,
    ) {
        settlementJdbc.update(
            """
            INSERT INTO daily_settlement
                (uid, settlement_date, currency, payment_method, gross_amount, refund_amount,
                 fee_amount, net_amount, payment_count, refund_count, status, created_at)
            VALUES (?, ?, ?, ?, ?, ?, 0, ?, ?, ?, 'OPEN', ?)
            """.trimIndent(),
            uuidToBytes(UUID.randomUUID()), date, currency, method, gross, refund,
            gross - refund, paymentCount, refundCount, LocalDateTime.now(ZoneOffset.UTC),
        )
    }

    private fun seedDetail(
        date: LocalDate,
        sourceType: String,
        sourceId: String,
        amount: Long,
        currency: String,
        method: String?,
        occurredAtUtc: LocalDateTime,
    ) {
        settlementJdbc.update(
            """
            INSERT INTO settlement_detail
                (uid, settlement_date, source_type, source_id, payment_id,
                 amount, currency, payment_method, occurred_at, created_at)
            VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
            """.trimIndent(),
            uuidToBytes(UUID.randomUUID()), date, sourceType, sourceId, sourceId,
            amount, currency, method, occurredAtUtc, LocalDateTime.now(ZoneOffset.UTC),
        )
    }

    private fun detailCount(date: LocalDate): Int =
        settlementJdbc.queryForObject(
            "SELECT COUNT(*) FROM settlement_detail WHERE settlement_date = ?", Int::class.java, date,
        )!!

    private fun bucketCount(date: LocalDate): Int =
        settlementJdbc.queryForObject(
            "SELECT COUNT(*) FROM daily_settlement WHERE settlement_date = ?", Int::class.java, date,
        )!!

    private fun jobRunStatus(date: LocalDate): String? =
        settlementJdbc.query(
            "SELECT status FROM settlement_job_run WHERE target_date = ? ORDER BY id DESC LIMIT 1",
            { rs, _ -> rs.getString("status") }, date,
        ).firstOrNull()

    private fun bucket(date: LocalDate, currency: String, method: String): Bucket =
        settlementJdbc.queryForObject(
            """
            SELECT gross_amount, refund_amount, fee_amount, net_amount,
                   payment_count, refund_count, status
              FROM daily_settlement
             WHERE settlement_date = ? AND currency = ? AND payment_method = ?
            """.trimIndent(),
            { rs, _ ->
                Bucket(
                    gross = rs.getLong("gross_amount"),
                    refund = rs.getLong("refund_amount"),
                    fee = rs.getLong("fee_amount"),
                    net = rs.getLong("net_amount"),
                    paymentCount = rs.getInt("payment_count"),
                    refundCount = rs.getInt("refund_count"),
                    status = rs.getString("status"),
                )
            },
            date, currency, method,
        )!!

    private data class Bucket(
        val gross: Long,
        val refund: Long,
        val fee: Long,
        val net: Long,
        val paymentCount: Int,
        val refundCount: Int,
        val status: String,
    )

    companion object {
        private const val PAYMENT_DB = "payment_db"

        @JvmStatic
        private val mysql: MySQLContainer<*> = MySQLContainer("mysql:8.4")
            .withDatabaseName("settlement_db")
            .withUsername("root")
            .withPassword("test")
            .also { it.start() }

        @JvmStatic
        private val paymentJdbc: JdbcTemplate = run {
            val base = mysql.jdbcUrl.substringBefore("/settlement_db")
            JdbcTemplate(
                DriverManagerDataSource(
                    "$base/$PAYMENT_DB?allowPublicKeyRetrieval=true&useSSL=false",
                    "root",
                    "test",
                ).apply { setDriverClassName("com.mysql.cj.jdbc.Driver") },
            )
        }

        @JvmStatic
        private var migrated = false

        private fun migrateOnce() {
            if (migrated) return
            JdbcTemplate(
                DriverManagerDataSource(mysql.jdbcUrl, "root", "test")
                    .apply { setDriverClassName("com.mysql.cj.jdbc.Driver") },
            ).execute("CREATE DATABASE IF NOT EXISTS $PAYMENT_DB")

            Flyway.configure()
                .dataSource(mysql.jdbcUrl, "root", "test")
                .locations("filesystem:db/migrations")
                .load()
                .migrate()

            val base = mysql.jdbcUrl.substringBefore("/settlement_db")
            Flyway.configure()
                .dataSource("$base/$PAYMENT_DB?allowPublicKeyRetrieval=true&useSSL=false", "root", "test")
                .locations("classpath:payment-fixture")
                .load()
                .migrate()

            migrated = true
        }

        @JvmStatic
        @DynamicPropertySource
        fun properties(registry: DynamicPropertyRegistry) {
            migrateOnce()

            registry.add("spring.datasource.url") { mysql.jdbcUrl }
            registry.add("spring.datasource.username") { "root" }
            registry.add("spring.datasource.password") { "test" }
            registry.add("spring.datasource.driver-class-name") { "com.mysql.cj.jdbc.Driver" }

            registry.add("settlement.timezone") { "Asia/Seoul" }
            registry.add("settlement.chunk-size") { 2 }
            registry.add("settlement.stale-execution-timeout") { "PT5M" }
            registry.add("settlement.payment-db.host") { mysql.host }
            registry.add("settlement.payment-db.port") { mysql.getMappedPort(3306) }
            registry.add("settlement.payment-db.params") { "allowPublicKeyRetrieval=true&useSSL=false" }
            registry.add("payment.db.username") { "root" }
            registry.add("payment.db.password") { "test" }
            registry.add("payment.db.name") { PAYMENT_DB }
        }
    }
}
