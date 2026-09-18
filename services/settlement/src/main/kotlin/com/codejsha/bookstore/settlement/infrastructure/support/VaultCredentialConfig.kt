package com.codejsha.bookstore.settlement.infrastructure.support

import com.codejsha.platform.shared.vault.VaultCredentialBinding
import com.codejsha.platform.shared.vault.VaultCredentialWatcher
import com.zaxxer.hikari.HikariDataSource
import org.springframework.beans.factory.annotation.Qualifier
import org.springframework.context.annotation.Bean
import org.springframework.context.annotation.Configuration
import java.nio.file.Path
import java.nio.file.Paths
import javax.sql.DataSource

@Configuration
class VaultCredentialConfig {

    @Bean(initMethod = "start", destroyMethod = "stop")
    fun vaultCredentialWatcher(
        dataSource: DataSource,
        @Qualifier("paymentDataSource") paymentDataSource: HikariDataSource,
    ): VaultCredentialWatcher =
        VaultCredentialWatcher(
            listOf(
                VaultCredentialBinding(DB_SECRETS_FILE, "db.username", "db.password") {
                    dataSource.unwrapHikari()
                },
                VaultCredentialBinding(PAYMENT_DB_SECRETS_FILE, PAYMENT_DB_USERNAME_KEY, PAYMENT_DB_PASSWORD_KEY) {
                    paymentDataSource
                },
            ),
        )

    @Bean(initMethod = "start", destroyMethod = "stop")
    fun readinessVaultCredentialWatcher(readinessDataSource: ReadinessDataSource): VaultCredentialWatcher =
        VaultCredentialWatcher(
            listOf(
                VaultCredentialBinding(DB_SECRETS_FILE, "db.username", "db.password") {
                    readinessDataSource.unwrapHikari()
                },
            ),
        )

    private fun DataSource.unwrapHikari(): HikariDataSource? =
        when (this) {
            is HikariDataSource -> this
            else -> runCatching { unwrap(HikariDataSource::class.java) }.getOrNull()
        }

    companion object {
        val DB_SECRETS_FILE: Path = Paths.get("/vault/secrets/db.properties")
        val PAYMENT_DB_SECRETS_FILE: Path = Paths.get("/vault/secrets/payment-db.properties")
        const val PAYMENT_DB_USERNAME_KEY = "payment.db.username"
        const val PAYMENT_DB_PASSWORD_KEY = "payment.db.password"
    }
}
