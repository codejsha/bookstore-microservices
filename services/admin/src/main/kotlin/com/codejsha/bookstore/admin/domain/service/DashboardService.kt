package com.codejsha.bookstore.admin.domain.service

import com.codejsha.bookstore.admin.application.port.restclient.CatalogClient
import com.codejsha.bookstore.admin.application.port.restclient.InventoryClient
import com.codejsha.bookstore.admin.application.port.restclient.OrderClient
import com.codejsha.bookstore.admin.application.port.support.ConcurrentCallContext
import com.codejsha.bookstore.admin.application.usecase.DashboardUseCase
import com.codejsha.bookstore.admin.domain.model.Dashboard
import com.codejsha.bookstore.admin.domain.model.Metric
import com.codejsha.platform.shared.data.ActorContext
import kotlinx.coroutines.CancellationException
import kotlinx.coroutines.async
import kotlinx.coroutines.runBlocking
import kotlinx.coroutines.runInterruptible
import kotlinx.coroutines.withTimeoutOrNull
import org.slf4j.LoggerFactory
import org.springframework.stereotype.Service
import java.time.Duration

@Service
class DashboardService(
    private val catalogClient: CatalogClient,
    private val orderClient: OrderClient,
    private val inventoryClient: InventoryClient,
    private val callContext: ConcurrentCallContext,
    private val metricTimeout: Duration = DEFAULT_METRIC_TIMEOUT,
) : DashboardUseCase {

    private val log = LoggerFactory.getLogger(javaClass)

    override fun loadDashboard(context: ActorContext): Dashboard = runBlocking(callContext.capture()) {
        val works = async { metric("catalog") { catalogClient.countWorks() } }
        val orders = async { metric("order") { orderClient.countOrders() } }
        val warehouses = async { metric("inventory") { inventoryClient.countWarehouses() } }
        Dashboard(works = works.await(), orders = orders.await(), warehouses = warehouses.await())
    }

    private suspend fun metric(downstream: String, count: () -> Long): Metric {
        val value = try {
            withTimeoutOrNull(metricTimeout.toMillis()) { runInterruptible { count() } }
        } catch (e: CancellationException) {
            throw e
        } catch (e: Exception) {
            log.warn("admin: downstream {} unavailable for dashboard: {}", downstream, e.toString())
            return Metric.unavailable()
        }
        if (value == null) {
            log.warn("admin: downstream {} did not answer within {} for dashboard", downstream, metricTimeout)
            return Metric.unavailable()
        }
        return Metric.of(value)
    }

    private companion object {
        private val DEFAULT_METRIC_TIMEOUT: Duration = Duration.ofSeconds(5)
    }
}
