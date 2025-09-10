package com.codejsha.bookstore.payment.infrastructure.support

import com.zaxxer.hikari.HikariDataSource
import jakarta.annotation.PostConstruct
import jakarta.annotation.PreDestroy
import org.slf4j.LoggerFactory
import org.springframework.stereotype.Component
import java.nio.file.FileSystems
import java.nio.file.Files
import java.nio.file.Path
import java.nio.file.Paths
import java.nio.file.StandardWatchEventKinds.ENTRY_CREATE
import java.nio.file.StandardWatchEventKinds.ENTRY_MODIFY
import java.util.Properties
import java.util.concurrent.Executors
import java.util.concurrent.ScheduledExecutorService
import java.util.concurrent.TimeUnit
import javax.sql.DataSource

@Component
class VaultCredentialWatcher(private val dataSource: DataSource) {

    private val log = LoggerFactory.getLogger(VaultCredentialWatcher::class.java)

    private val secretsFile: Path = Paths.get(SECRETS_FILE_PATH)
    private val watchExecutor: ScheduledExecutorService =
        Executors.newSingleThreadScheduledExecutor { r ->
            Thread(r, "vault-cred-watcher").apply { isDaemon = true }
        }

    @Volatile
    private var lastModifiedTime: Long = -1L

    @PostConstruct
    fun start() {
        if (!Files.exists(secretsFile)) {
            log.info(
                "vault-cred-watcher: {} not found — credential rotation watcher is inactive (local dev mode)",
                SECRETS_FILE_PATH,
            )
            return
        }

        lastModifiedTime = Files.getLastModifiedTime(secretsFile).toMillis()

        startWatchService()
        startPollFallback()
        log.info("vault-cred-watcher: watching {} for Vault credential rotation", SECRETS_FILE_PATH)
    }

    @PreDestroy
    fun stop() {
        watchExecutor.shutdownNow()
    }

    private fun startWatchService() {
        val watchService = FileSystems.getDefault().newWatchService()
        val parentDir = secretsFile.parent ?: Paths.get("/vault/secrets")
        parentDir.register(watchService, ENTRY_MODIFY, ENTRY_CREATE)

        watchExecutor.submit {
            try {
                while (!Thread.currentThread().isInterrupted) {
                    val key = watchService.take()
                    val relevant = key.pollEvents().any { event ->
                        val changed = parentDir.resolve(event.context() as Path)
                        Files.isSameFile(changed, secretsFile)
                    }
                    if (relevant) {
                        reloadCredentials("watch-event")
                    }
                    if (!key.reset()) break
                }
            } catch (_: InterruptedException) {
                Thread.currentThread().interrupt()
            } finally {
                watchService.close()
            }
        }
    }

    private fun startPollFallback() {
        watchExecutor.scheduleWithFixedDelay(
            {
                if (!Files.exists(secretsFile)) return@scheduleWithFixedDelay
                val currentMtime = Files.getLastModifiedTime(secretsFile).toMillis()
                if (currentMtime != lastModifiedTime) {
                    reloadCredentials("mtime-poll")
                }
            },
            POLL_INTERVAL_SECONDS,
            POLL_INTERVAL_SECONDS,
            TimeUnit.SECONDS,
        )
    }

    private fun reloadCredentials(trigger: String) {
        try {
            val props = Properties()
            Files.newBufferedReader(secretsFile).use { props.load(it) }

            val newUsername = props.getProperty("db.username")
            val newPassword = props.getProperty("db.password")

            if (newUsername.isNullOrBlank() || newPassword.isNullOrBlank()) {
                log.warn(
                    "vault-cred-watcher [{}]: db.username or db.password missing in {} — skipping update",
                    trigger,
                    SECRETS_FILE_PATH,
                )
                return
            }

            val hikari = unwrapHikari()
            if (hikari == null) {
                log.warn("vault-cred-watcher [{}]: DataSource is not a HikariDataSource — cannot rotate credentials", trigger)
                return
            }

            hikari.hikariConfigMXBean.setUsername(newUsername)
            hikari.hikariConfigMXBean.setPassword(newPassword)
            hikari.hikariPoolMXBean?.softEvictConnections()

            lastModifiedTime = Files.getLastModifiedTime(secretsFile).toMillis()

            log.info(
                "vault-cred-watcher [{}]: HikariCP credentials updated for user '{}' and pool soft-eviction triggered",
                trigger,
                newUsername,
            )
        } catch (e: Exception) {
            log.error("vault-cred-watcher [{}]: failed to reload credentials — {}", trigger, e.message, e)
        }
    }

    private fun unwrapHikari(): HikariDataSource? =
        when (dataSource) {
            is HikariDataSource -> dataSource
            else ->
                try {
                    dataSource.unwrap(HikariDataSource::class.java)
                } catch (_: Exception) {
                    null
                }
        }

    companion object {
        private const val SECRETS_FILE_PATH = "/vault/secrets/db.properties"
        private const val POLL_INTERVAL_SECONDS = 30L
    }
}
