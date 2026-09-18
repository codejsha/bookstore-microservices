package com.codejsha.bookstore.settlement.infrastructure.adapter.mysql

import com.codejsha.bookstore.settlement.application.port.SettlementRunLock
import com.codejsha.bookstore.settlement.config.properties.SettlementBatchProperties
import com.codejsha.bookstore.settlement.infrastructure.support.utils.uuidToBytes
import jakarta.annotation.PreDestroy
import org.slf4j.LoggerFactory
import org.springframework.dao.DuplicateKeyException
import org.springframework.jdbc.core.JdbcTemplate
import org.springframework.stereotype.Repository
import java.time.LocalDate
import java.time.LocalDateTime
import java.time.ZoneOffset
import java.util.concurrent.ConcurrentHashMap
import java.util.concurrent.Executors
import java.util.concurrent.ScheduledFuture
import java.util.concurrent.TimeUnit
import kotlin.uuid.Uuid
import kotlin.uuid.toJavaUuid
import javax.sql.DataSource

@Repository
class MysqlSettlementRunLock(
    dataSource: DataSource,
    private val properties: SettlementBatchProperties,
) : SettlementRunLock {

    private val log = LoggerFactory.getLogger(javaClass)
    private val jdbc = JdbcTemplate(dataSource)
    private val ownerHost = System.getenv("HOSTNAME") ?: "unknown"
    private val heartbeats = ConcurrentHashMap<String, ScheduledFuture<*>>()
    private val scheduler = Executors.newSingleThreadScheduledExecutor { runnable ->
        Thread(runnable, "settlement-lock-heartbeat").apply { isDaemon = true }
    }

    override fun tryAcquire(targetDate: LocalDate, ownerToken: String): Boolean {
        val now = LocalDateTime.now(ZoneOffset.UTC)
        val acquired = insert(targetDate, ownerToken, now) || takeOverLapsed(targetDate, ownerToken, now)
        if (acquired) {
            startHeartbeat(targetDate, ownerToken)
            log.info("Acquired the settlement run lock for {} as {}", targetDate, ownerToken)
        }
        return acquired
    }

    override fun isHeldBy(targetDate: LocalDate, ownerToken: String): Boolean {
        val held = jdbc.queryForObject(
            """
            SELECT COUNT(*) FROM settlement_run_lock
             WHERE target_date = ? AND owner_token = ?
               AND released_at IS NULL AND deleted_at IS NULL AND expires_at > ?
            """.trimIndent(),
            Int::class.java,
            targetDate, ownerToken, LocalDateTime.now(ZoneOffset.UTC),
        ) ?: 0
        return held > 0
    }

    override fun release(targetDate: LocalDate, ownerToken: String): Boolean {
        stopHeartbeat(targetDate, ownerToken)
        val now = LocalDateTime.now(ZoneOffset.UTC)
        return jdbc.update(
            """
            UPDATE settlement_run_lock
               SET released_at = ?, expires_at = ?, updated_at = ?
             WHERE target_date = ? AND owner_token = ? AND released_at IS NULL
            """.trimIndent(),
            now, now, now, targetDate, ownerToken,
        ) > 0
    }

    private fun insert(targetDate: LocalDate, ownerToken: String, now: LocalDateTime): Boolean =
        try {
            jdbc.update(
                """
                INSERT INTO settlement_run_lock
                    (uid, target_date, owner_token, owner_host, acquired_at, heartbeat_at, expires_at, created_at)
                VALUES (?, ?, ?, ?, ?, ?, ?, ?)
                """.trimIndent(),
                uuidToBytes(Uuid.generateV7().toJavaUuid()), targetDate, ownerToken, ownerHost,
                now, now, now.plus(properties.lock.ttl), now,
            ) > 0
        } catch (e: DuplicateKeyException) {
            false
        }

    private fun takeOverLapsed(targetDate: LocalDate, ownerToken: String, now: LocalDateTime): Boolean =
        jdbc.update(
            """
            UPDATE settlement_run_lock
               SET uid = ?, owner_token = ?, owner_host = ?, acquired_at = ?, heartbeat_at = ?,
                   expires_at = ?, released_at = NULL, deleted_at = NULL, updated_at = ?
             WHERE target_date = ?
               AND (released_at IS NOT NULL OR deleted_at IS NOT NULL OR expires_at <= ?)
            """.trimIndent(),
            uuidToBytes(Uuid.generateV7().toJavaUuid()), ownerToken, ownerHost, now, now,
            now.plus(properties.lock.ttl), now, targetDate, now,
        ) > 0

    private fun renew(targetDate: LocalDate, ownerToken: String): Boolean {
        val now = LocalDateTime.now(ZoneOffset.UTC)
        return jdbc.update(
            """
            UPDATE settlement_run_lock
               SET heartbeat_at = ?, expires_at = ?, updated_at = ?
             WHERE target_date = ? AND owner_token = ? AND released_at IS NULL AND deleted_at IS NULL
            """.trimIndent(),
            now, now.plus(properties.lock.ttl), now, targetDate, ownerToken,
        ) > 0
    }

    private fun startHeartbeat(targetDate: LocalDate, ownerToken: String) {
        val interval = properties.lock.heartbeatInterval.toMillis().coerceAtLeast(1_000L)
        val future = scheduler.scheduleWithFixedDelay(
            { heartbeat(targetDate, ownerToken) },
            interval, interval, TimeUnit.MILLISECONDS,
        )
        heartbeats.put(heartbeatKey(targetDate, ownerToken), future)?.cancel(false)
    }

    private fun heartbeat(targetDate: LocalDate, ownerToken: String) {
        try {
            if (!renew(targetDate, ownerToken)) {
                log.warn("The settlement run lock for {} is no longer owned by {}", targetDate, ownerToken)
                stopHeartbeat(targetDate, ownerToken)
            }
        } catch (e: Exception) {
            log.warn("Failed to renew the settlement run lock for {} as {}", targetDate, ownerToken, e)
        }
    }

    private fun stopHeartbeat(targetDate: LocalDate, ownerToken: String) {
        heartbeats.remove(heartbeatKey(targetDate, ownerToken))?.cancel(false)
    }

    private fun heartbeatKey(targetDate: LocalDate, ownerToken: String): String = "$targetDate/$ownerToken"

    @PreDestroy
    fun shutdown() {
        heartbeats.values.forEach { it.cancel(false) }
        heartbeats.clear()
        scheduler.shutdownNow()
    }
}
