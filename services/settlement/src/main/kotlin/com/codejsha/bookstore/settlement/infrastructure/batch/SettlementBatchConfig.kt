package com.codejsha.bookstore.settlement.infrastructure.batch

import com.codejsha.bookstore.settlement.application.port.PgReconciliationPort
import com.codejsha.bookstore.settlement.config.properties.SettlementBatchProperties
import com.codejsha.bookstore.settlement.domain.constant.SettlementSourceType
import com.codejsha.bookstore.settlement.domain.constant.SettlementStatus
import com.codejsha.bookstore.settlement.infrastructure.support.utils.uuidToBytes
import org.springframework.batch.core.configuration.annotation.StepScope
import org.springframework.batch.core.job.Job
import org.springframework.batch.core.job.builder.JobBuilder
import org.springframework.batch.core.repository.JobRepository
import org.springframework.batch.core.step.Step
import org.springframework.batch.core.step.builder.StepBuilder
import org.springframework.batch.core.step.tasklet.Tasklet
import org.springframework.batch.infrastructure.item.ItemProcessor
import org.springframework.batch.infrastructure.item.database.ItemPreparedStatementSetter
import org.springframework.batch.infrastructure.item.database.JdbcBatchItemWriter
import org.springframework.batch.infrastructure.item.database.JdbcPagingItemReader
import org.springframework.batch.infrastructure.item.database.Order
import org.springframework.batch.infrastructure.item.database.builder.JdbcBatchItemWriterBuilder
import org.springframework.batch.infrastructure.item.database.builder.JdbcPagingItemReaderBuilder
import org.springframework.batch.infrastructure.item.database.support.MySqlPagingQueryProvider
import org.springframework.batch.infrastructure.repeat.RepeatStatus
import org.springframework.beans.factory.annotation.Qualifier
import org.springframework.beans.factory.annotation.Value
import org.springframework.context.annotation.Bean
import org.springframework.context.annotation.Configuration
import org.springframework.dao.TransientDataAccessException
import org.springframework.jdbc.core.JdbcTemplate
import org.springframework.transaction.PlatformTransactionManager
import java.time.LocalDate
import java.time.LocalDateTime
import java.time.ZoneId
import java.time.ZoneOffset
import java.util.UUID
import kotlin.uuid.Uuid
import kotlin.uuid.toJavaUuid
import javax.sql.DataSource

@Configuration
class SettlementBatchConfig(
    private val properties: SettlementBatchProperties,
) {

    private val zone: ZoneId get() = ZoneId.of(properties.timezone)

    // ─── Job ──────────────────────────────────────────────────────────────────

    @Bean
    fun dailySettlementJob(
        jobRepository: JobRepository,
        listener: SettlementJobListener,
        cleanupStep: Step,
        extractPaymentsStep: Step,
        extractRefundsStep: Step,
        aggregateStep: Step,
        reconcileStep: Step,
    ): Job =
        JobBuilder("dailySettlementJob", jobRepository)
            .listener(listener)
            .start(cleanupStep)
            .next(extractPaymentsStep)
            .next(extractRefundsStep)
            .next(aggregateStep)
            .next(reconcileStep)
            .build()

    // ─── Step 0: cleanup / idempotency (delete-insert) ─────────────────────────

    @Bean
    fun cleanupStep(
        jobRepository: JobRepository,
        transactionManager: PlatformTransactionManager,
        cleanupTasklet: Tasklet,
    ): Step =
        StepBuilder("cleanupStep", jobRepository)
            .tasklet(cleanupTasklet, transactionManager)
            .build()

    @Bean
    @StepScope
    fun cleanupTasklet(
        dataSource: DataSource,
        @Value("#{jobParameters['targetDate']}") targetDate: LocalDate,
        @Value("#{jobParameters['runUid']}") runUid: String?,
    ): Tasklet =
        Tasklet { contribution, chunkContext ->
            val jdbc = JdbcTemplate(dataSource)
            val now = LocalDateTime.now(ZoneOffset.UTC)
            jdbc.update("DELETE FROM settlement_detail WHERE settlement_date = ?", targetDate)
            jdbc.update("DELETE FROM daily_settlement WHERE settlement_date = ?", targetDate)

            val jobExecutionId = chunkContext.stepContext.stepExecution.jobExecutionId
            val uid = runUid?.let { UUID.fromString(it) } ?: Uuid.generateV7().toJavaUuid()
            jdbc.update("DELETE FROM settlement_job_run WHERE job_execution_id = ?", jobExecutionId)
            jdbc.update(
                """
                INSERT INTO settlement_job_run
                    (uid, target_date, job_execution_id, status, started_at, created_at)
                VALUES (?, ?, ?, 'STARTED', ?, ?)
                """.trimIndent(),
                uuidToBytes(uid), targetDate, jobExecutionId, now, now,
            )
            contribution.incrementWriteCount(1L)
            RepeatStatus.FINISHED
        }

    // ─── Step 1a: extract captured payments ────────────────────────────────────

    @Bean
    fun extractPaymentsStep(
        jobRepository: JobRepository,
        transactionManager: PlatformTransactionManager,
        paymentItemReader: JdbcPagingItemReader<PaymentRow>,
        paymentItemProcessor: ItemProcessor<PaymentRow, SettlementDetailInsert>,
        settlementDetailWriter: JdbcBatchItemWriter<SettlementDetailInsert>,
    ): Step =
        StepBuilder("extractPaymentsStep", jobRepository)
            .chunk<PaymentRow, SettlementDetailInsert>(properties.chunkSize, transactionManager)
            .reader(paymentItemReader)
            .processor(paymentItemProcessor)
            .writer(settlementDetailWriter)
            .faultTolerant()
            .retry(TransientDataAccessException::class.java)
            .retryLimit(3)
            .build()

    @Bean
    @StepScope
    fun paymentItemReader(
        @Qualifier("paymentDataSource") paymentDataSource: DataSource,
        @Value("#{jobParameters['targetDate']}") targetDate: LocalDate,
    ): JdbcPagingItemReader<PaymentRow> {
        val window = SettlementWindow.of(targetDate, zone)
        val provider = MySqlPagingQueryProvider().apply {
            setSelectClause("payment_id, status, amount, amount_captured, currency, payment_method, captured_at")
            setFromClause("payments")
            setWhereClause(
                "status IN ('succeeded', 'partially_captured') " +
                    "AND captured_at >= :startAt AND captured_at < :endAt AND deleted_at IS NULL",
            )
            setSortKeys(mapOf("payment_id" to Order.ASCENDING))
        }
        return JdbcPagingItemReaderBuilder<PaymentRow>()
            .name("paymentItemReader")
            .dataSource(paymentDataSource)
            .queryProvider(provider)
            .parameterValues(mapOf("startAt" to window.startUtc, "endAt" to window.endUtc))
            .pageSize(properties.chunkSize)
            .rowMapper { rs, _ ->
                PaymentRow(
                    paymentId = rs.getString("payment_id"),
                    status = rs.getString("status"),
                    amount = rs.getLong("amount"),
                    amountCaptured = rs.getLong("amount_captured").let { if (rs.wasNull()) null else it },
                    currency = rs.getString("currency"),
                    paymentMethod = rs.getString("payment_method"),
                    capturedAt = rs.getObject("captured_at", LocalDateTime::class.java),
                )
            }
            .build()
    }

    @Bean
    @StepScope
    fun paymentItemProcessor(
        @Value("#{jobParameters['targetDate']}") targetDate: LocalDate,
    ): ItemProcessor<PaymentRow, SettlementDetailInsert> =
        ItemProcessor { row ->
            SettlementDetailInsert(
                uid = Uuid.generateV7().toJavaUuid(),
                settlementDate = targetDate,
                sourceType = SettlementSourceType.PAYMENT.value,
                sourceId = row.paymentId,
                paymentId = row.paymentId,
                amount = row.amountCaptured
                    ?: if (row.status == "partially_captured") 0L else row.amount,
                currency = row.currency,
                paymentMethod = row.paymentMethod,
                occurredAt = row.capturedAt,
                createdAt = LocalDateTime.now(ZoneOffset.UTC),
            )
        }

    // ─── Step 1b: extract succeeded refunds ────────────────────────────────────

    @Bean
    fun extractRefundsStep(
        jobRepository: JobRepository,
        transactionManager: PlatformTransactionManager,
        refundItemReader: JdbcPagingItemReader<RefundRow>,
        refundItemProcessor: ItemProcessor<RefundRow, SettlementDetailInsert>,
        settlementDetailWriter: JdbcBatchItemWriter<SettlementDetailInsert>,
    ): Step =
        StepBuilder("extractRefundsStep", jobRepository)
            .chunk<RefundRow, SettlementDetailInsert>(properties.chunkSize, transactionManager)
            .reader(refundItemReader)
            .processor(refundItemProcessor)
            .writer(settlementDetailWriter)
            .faultTolerant()
            .retry(TransientDataAccessException::class.java)
            .retryLimit(3)
            .build()

    @Bean
    @StepScope
    fun refundItemReader(
        @Qualifier("paymentDataSource") paymentDataSource: DataSource,
        @Value("#{jobParameters['targetDate']}") targetDate: LocalDate,
    ): JdbcPagingItemReader<RefundRow> {
        val window = SettlementWindow.of(targetDate, zone)
        val provider = MySqlPagingQueryProvider().apply {
            setSelectClause(
                "r.refund_id AS refund_id, r.payment_id AS payment_id, r.amount AS amount, " +
                    "r.currency AS currency, p.payment_method AS payment_method, r.created_at AS created_at",
            )
            setFromClause("refunds r JOIN payments p ON p.payment_id = r.payment_id")
            setWhereClause(
                "r.status = 'succeeded' " +
                    "AND r.created_at >= :startAt AND r.created_at < :endAt AND r.deleted_at IS NULL",
            )
            setSortKeys(mapOf("r.refund_id" to Order.ASCENDING))
        }
        return JdbcPagingItemReaderBuilder<RefundRow>()
            .name("refundItemReader")
            .dataSource(paymentDataSource)
            .queryProvider(provider)
            .parameterValues(mapOf("startAt" to window.startUtc, "endAt" to window.endUtc))
            .pageSize(properties.chunkSize)
            .rowMapper { rs, _ ->
                RefundRow(
                    refundId = rs.getString("refund_id"),
                    paymentId = rs.getString("payment_id"),
                    amount = rs.getLong("amount"),
                    currency = rs.getString("currency"),
                    paymentMethod = rs.getString("payment_method"),
                    createdAt = rs.getObject("created_at", LocalDateTime::class.java),
                )
            }
            .build()
    }

    @Bean
    @StepScope
    fun refundItemProcessor(
        @Value("#{jobParameters['targetDate']}") targetDate: LocalDate,
    ): ItemProcessor<RefundRow, SettlementDetailInsert> =
        ItemProcessor { row ->
            SettlementDetailInsert(
                uid = Uuid.generateV7().toJavaUuid(),
                settlementDate = targetDate,
                sourceType = SettlementSourceType.REFUND.value,
                sourceId = row.refundId,
                paymentId = row.paymentId,
                amount = -row.amount,
                currency = row.currency,
                paymentMethod = row.paymentMethod,
                occurredAt = row.createdAt,
                createdAt = LocalDateTime.now(ZoneOffset.UTC),
            )
        }

    // ─── Shared writer → settlement_detail (primary DataSource) ────────────────

    @Bean
    fun settlementDetailWriter(dataSource: DataSource): JdbcBatchItemWriter<SettlementDetailInsert> {
        val setter = ItemPreparedStatementSetter<SettlementDetailInsert> { item, ps ->
            ps.setBytes(1, uuidToBytes(item.uid))
            ps.setObject(2, item.settlementDate)
            ps.setString(3, item.sourceType)
            ps.setString(4, item.sourceId)
            ps.setString(5, item.paymentId)
            ps.setLong(6, item.amount)
            ps.setString(7, item.currency)
            ps.setString(8, item.paymentMethod)
            ps.setObject(9, item.occurredAt)
            ps.setObject(10, item.createdAt)
        }
        return JdbcBatchItemWriterBuilder<SettlementDetailInsert>()
            .dataSource(dataSource)
            .sql(
                """
                INSERT INTO settlement_detail
                    (uid, settlement_date, source_type, source_id, payment_id,
                     amount, currency, payment_method, occurred_at, created_at)
                VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
                """.trimIndent(),
            )
            .itemPreparedStatementSetter(setter)
            .assertUpdates(true)
            .build()
    }

    // ─── Step 2: aggregate settlement_detail → daily_settlement ────────────────

    @Bean
    fun aggregateStep(
        jobRepository: JobRepository,
        transactionManager: PlatformTransactionManager,
        aggregateTasklet: Tasklet,
    ): Step =
        StepBuilder("aggregateStep", jobRepository)
            .tasklet(aggregateTasklet, transactionManager)
            .build()

    @Bean
    @StepScope
    fun aggregateTasklet(
        dataSource: DataSource,
        @Value("#{jobParameters['targetDate']}") targetDate: LocalDate,
    ): Tasklet =
        Tasklet { contribution, _ ->
            val jdbc = JdbcTemplate(dataSource)
            val now = LocalDateTime.now(ZoneOffset.UTC)
            val buckets = jdbc.query(
                """
                SELECT currency,
                       payment_method,
                       SUM(CASE WHEN source_type = 'PAYMENT' THEN amount ELSE 0 END) AS gross,
                       SUM(CASE WHEN source_type = 'REFUND' THEN -amount ELSE 0 END) AS refund,
                       SUM(CASE WHEN source_type = 'PAYMENT' THEN 1 ELSE 0 END) AS pay_count,
                       SUM(CASE WHEN source_type = 'REFUND' THEN 1 ELSE 0 END) AS refund_count
                  FROM settlement_detail
                 WHERE settlement_date = ? AND deleted_at IS NULL
                 GROUP BY currency, payment_method
                """.trimIndent(),
                { rs, _ ->
                    AggregateBucket(
                        currency = rs.getString("currency"),
                        paymentMethod = rs.getString("payment_method"),
                        gross = rs.getLong("gross"),
                        refund = rs.getLong("refund"),
                        paymentCount = rs.getInt("pay_count"),
                        refundCount = rs.getInt("refund_count"),
                    )
                },
                targetDate,
            )

            buckets.forEach { b ->
                val fee = 0L
                val net = b.gross - b.refund - fee
                jdbc.update(
                    """
                    INSERT INTO daily_settlement
                        (uid, settlement_date, currency, payment_method,
                         gross_amount, refund_amount, fee_amount, net_amount,
                         payment_count, refund_count, status, created_at)
                    VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
                    """.trimIndent(),
                    uuidToBytes(Uuid.generateV7().toJavaUuid()), targetDate, b.currency, b.paymentMethod,
                    b.gross, b.refund, fee, net,
                    b.paymentCount, b.refundCount, SettlementStatus.OPEN.value, now,
                )
            }
            contribution.incrementWriteCount(buckets.size.toLong())
            RepeatStatus.FINISHED
        }

    // ─── Step 3: reconcile ─────────────────────────────────────────────────────

    @Bean
    fun reconcileStep(
        jobRepository: JobRepository,
        transactionManager: PlatformTransactionManager,
        reconcileTasklet: Tasklet,
    ): Step =
        StepBuilder("reconcileStep", jobRepository)
            .tasklet(reconcileTasklet, transactionManager)
            .build()

    @Bean
    @StepScope
    fun reconcileTasklet(
        dataSource: DataSource,
        pgReconciliationPort: PgReconciliationPort,
        @Value("#{jobParameters['targetDate']}") targetDate: LocalDate,
    ): Tasklet =
        Tasklet { contribution, chunkContext ->
            val jdbc = JdbcTemplate(dataSource)
            val window = SettlementWindow.of(targetDate, zone)
            val jobExecutionId = chunkContext.stepContext.stepExecution.jobExecutionId

            resetDiscrepancies(dataSource, targetDate)

            val detailByCurrency = jdbc.query(
                """
                SELECT currency,
                       SUM(CASE WHEN source_type = 'PAYMENT' THEN amount ELSE 0 END) AS gross,
                       SUM(CASE WHEN source_type = 'REFUND' THEN -amount ELSE 0 END) AS refund,
                       SUM(CASE WHEN source_type = 'PAYMENT' THEN 1 ELSE 0 END) AS pay_count,
                       SUM(CASE WHEN source_type = 'REFUND' THEN 1 ELSE 0 END) AS refund_count
                  FROM settlement_detail
                 WHERE settlement_date = ? AND deleted_at IS NULL
                 GROUP BY currency
                """.trimIndent(),
                { rs, _ ->
                    rs.getString("currency") to CurrencyTotals(
                        gross = rs.getLong("gross"),
                        refund = rs.getLong("refund"),
                        paymentCount = rs.getLong("pay_count"),
                        refundCount = rs.getLong("refund_count"),
                    )
                },
                targetDate,
            ).toMap()
            val settledByCurrency = jdbc.query(
                """
                SELECT currency,
                       COALESCE(SUM(gross_amount), 0)  AS gross,
                       COALESCE(SUM(refund_amount), 0) AS refund,
                       COALESCE(SUM(payment_count), 0) AS pay_count,
                       COALESCE(SUM(refund_count), 0)  AS refund_count
                  FROM daily_settlement
                 WHERE settlement_date = ? AND deleted_at IS NULL
                 GROUP BY currency
                """.trimIndent(),
                { rs, _ ->
                    rs.getString("currency") to CurrencyTotals(
                        gross = rs.getLong("gross"),
                        refund = rs.getLong("refund"),
                        paymentCount = rs.getLong("pay_count"),
                        refundCount = rs.getLong("refund_count"),
                    )
                },
                targetDate,
            ).toMap()

            val internalMismatches = (detailByCurrency.keys + settledByCurrency.keys)
                .filter { detailByCurrency[it] != settledByCurrency[it] }
                .sorted()
            if (internalMismatches.isNotEmpty()) {
                val message = "Internal reconcile mismatch for $targetDate: " +
                    internalMismatches.joinToString("; ") { currency ->
                        "$currency detail=${detailByCurrency[currency]} vs settled=${settledByCurrency[currency]}"
                    }
                markDiscrepancy(dataSource, targetDate, jobExecutionId, message, internalMismatches)
                error(message)
            }

            if (properties.reconcile.hyperswitch.enabled) {
                val pgTotals = pgReconciliationPort
                    .fetchDailyTotals(targetDate, window.startUtc, window.endUtc)
                    .associateBy { it.currency }
                val pgMismatches = (settledByCurrency.keys + pgTotals.keys)
                    .filter { currency ->
                        val settled = settledByCurrency[currency]
                        val pg = pgTotals[currency]
                        settled == null || pg == null ||
                            pg.grossAmount != settled.gross || pg.count != settled.paymentCount
                    }
                    .sorted()
                if (pgMismatches.isNotEmpty()) {
                    val message = "Hyperswitch reconcile mismatch for $targetDate: " +
                        pgMismatches.joinToString("; ") { currency ->
                            "$currency settled=${settledByCurrency[currency]} vs pg=${pgTotals[currency]}"
                        }
                    markDiscrepancy(dataSource, targetDate, jobExecutionId, message, pgMismatches)
                    error(message)
                }
            }

            val now = LocalDateTime.now(ZoneOffset.UTC)
            val confirmed = jdbc.update(
                """
                UPDATE daily_settlement
                   SET status = ?, updated_at = ?
                 WHERE settlement_date = ? AND status = ? AND deleted_at IS NULL
                """.trimIndent(),
                SettlementStatus.CONFIRMED.value, now, targetDate, SettlementStatus.OPEN.value,
            )
            contribution.incrementWriteCount(confirmed.toLong())
            RepeatStatus.FINISHED
        }

    private fun resetDiscrepancies(
        dataSource: DataSource,
        targetDate: LocalDate,
    ) {
        val now = LocalDateTime.now(ZoneOffset.UTC)
        dataSource.connection.use { conn ->
            conn.autoCommit = true
            conn.prepareStatement(
                """
                UPDATE daily_settlement
                   SET status = ?, updated_at = ?
                 WHERE settlement_date = ? AND status = ? AND deleted_at IS NULL
                """.trimIndent(),
            ).use { ps ->
                ps.setString(1, SettlementStatus.OPEN.value)
                ps.setObject(2, now)
                ps.setObject(3, targetDate)
                ps.setString(4, SettlementStatus.DISCREPANCY.value)
                ps.executeUpdate()
            }
        }
    }

    private fun markDiscrepancy(
        dataSource: DataSource,
        targetDate: LocalDate,
        jobExecutionId: Long,
        message: String,
        currencies: Collection<String>,
    ) {
        if (currencies.isEmpty()) return
        val now = LocalDateTime.now(ZoneOffset.UTC)
        dataSource.connection.use { conn ->
            conn.autoCommit = false
            try {
                val placeholders = currencies.joinToString(",") { "?" }
                conn.prepareStatement(
                    """
                    UPDATE daily_settlement
                       SET status = ?, updated_at = ?
                     WHERE settlement_date = ? AND currency IN ($placeholders) AND deleted_at IS NULL
                    """.trimIndent(),
                ).use { ps ->
                    ps.setString(1, SettlementStatus.DISCREPANCY.value)
                    ps.setObject(2, now)
                    ps.setObject(3, targetDate)
                    currencies.forEachIndexed { i, currency -> ps.setString(4 + i, currency) }
                    ps.executeUpdate()
                }
                conn.prepareStatement(
                    "UPDATE settlement_job_run SET message = ?, updated_at = ? WHERE job_execution_id = ?",
                ).use { ps ->
                    ps.setString(1, message)
                    ps.setObject(2, now)
                    ps.setLong(3, jobExecutionId)
                    ps.executeUpdate()
                }
                conn.commit()
            } catch (e: Exception) {
                conn.rollback()
                throw e
            }
        }
    }

    private data class AggregateBucket(
        val currency: String,
        val paymentMethod: String?,
        val gross: Long,
        val refund: Long,
        val paymentCount: Int,
        val refundCount: Int,
    )

    private data class CurrencyTotals(
        val gross: Long,
        val refund: Long,
        val paymentCount: Long,
        val refundCount: Long,
    )
}
