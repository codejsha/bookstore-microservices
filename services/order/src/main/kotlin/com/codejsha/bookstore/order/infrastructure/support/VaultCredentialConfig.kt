package com.codejsha.bookstore.order.infrastructure.support

import com.codejsha.platform.shared.vault.VaultCredentialBinding
import com.codejsha.platform.shared.vault.VaultCredentialWatcher
import com.zaxxer.hikari.HikariDataSource
import org.springframework.context.annotation.Bean
import org.springframework.context.annotation.Configuration
import java.nio.file.Paths
import javax.sql.DataSource

@Configuration
class VaultCredentialConfig {

    @Bean(initMethod = "start", destroyMethod = "stop")
    fun vaultCredentialWatcher(dataSource: DataSource): VaultCredentialWatcher =
        VaultCredentialWatcher(
            listOf(
                VaultCredentialBinding(Paths.get(DB_SECRETS_FILE), "db.username", "db.password") {
                    dataSource.unwrapHikari()
                },
            ),
        )

    private fun DataSource.unwrapHikari(): HikariDataSource? =
        when (this) {
            is HikariDataSource -> this
            else -> runCatching { unwrap(HikariDataSource::class.java) }.getOrNull()
        }

    companion object {
        private const val DB_SECRETS_FILE = "/vault/secrets/db.properties"
    }
}
