package com.codejsha.bookstore.settlement.config.properties

import org.springframework.boot.context.properties.ConfigurationProperties
import java.time.Duration

@ConfigurationProperties(prefix = "settlement")
data class SettlementBatchProperties(
    val timezone: String = "Asia/Seoul",
    val chunkSize: Int = 500,
    val staleExecutionTimeout: Duration = Duration.ofMinutes(30),
    val lock: SettlementRunLockProperties = SettlementRunLockProperties(),
    val paymentDb: PaymentDbConnection = PaymentDbConnection(),
    val reconcile: ReconcileProperties = ReconcileProperties(),
)

data class SettlementRunLockProperties(
    val ttl: Duration = Duration.ofMinutes(5),
    val heartbeatInterval: Duration = Duration.ofMinutes(1),
)

data class PaymentDbConnection(
    val host: String = "localhost",
    val port: Int = 3306,
    val params: String = "allowPublicKeyRetrieval=true&characterEncoding=utf8&useSSL=false",
)

data class ReconcileProperties(
    val hyperswitch: HyperswitchReconcileProperties = HyperswitchReconcileProperties(),
)

data class HyperswitchReconcileProperties(
    val enabled: Boolean = false,
)

@ConfigurationProperties(prefix = "payment.db")
data class PaymentDbProperties(
    val username: String = "",
    val password: String = "",
    val name: String = "payment_db",
)
