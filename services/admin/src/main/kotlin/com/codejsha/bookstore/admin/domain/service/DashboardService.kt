package com.codejsha.bookstore.admin.domain.service

import com.codejsha.bookstore.admin.application.port.restclient.CatalogClient
import com.codejsha.bookstore.admin.application.port.restclient.InventoryClient
import com.codejsha.bookstore.admin.application.port.restclient.OrderClient
import com.codejsha.bookstore.admin.application.usecase.DashboardUseCase
import com.codejsha.bookstore.admin.domain.model.Dashboard
import com.codejsha.bookstore.admin.domain.model.Metric
import com.codejsha.platform.shared.data.ActorContext
import org.slf4j.LoggerFactory
import org.springframework.stereotype.Service

@Service
class DashboardService(
    private val catalogClient: CatalogClient,
    private val orderClient: OrderClient,
    private val inventoryClient: InventoryClient,
) : DashboardUseCase {

    private val log = LoggerFactory.getLogger(javaClass)

    override suspend fun loadDashboard(context: ActorContext) = Dashboard(
        works = metric("catalog") { catalogClient.countWorks() },
        orders = metric("order") { orderClient.countOrders() },
        warehouses = metric("inventory") { inventoryClient.countWarehouses() },
    )

    private fun metric(downstream: String, count: () -> Long): Metric =
        try {
            Metric.of(count())
        } catch (e: Exception) {
            log.warn("admin: downstream {} unavailable for dashboard: {}", downstream, e.toString())
            Metric.unavailable()
        }
}
