package com.codejsha.bookstore.admin.domain.service

import com.codejsha.bookstore.admin.application.port.restclient.CatalogClient
import com.codejsha.bookstore.admin.application.port.restclient.InventoryClient
import com.codejsha.bookstore.admin.application.port.restclient.OrderClient
import com.codejsha.bookstore.admin.domain.model.command.WorkCreateCommand
import com.codejsha.bookstore.admin.domain.model.command.WorkUpdateCommand
import com.codejsha.bookstore.admin.domain.model.option.OrderQueryOption
import com.codejsha.bookstore.admin.domain.model.option.StockQueryOption
import com.codejsha.bookstore.admin.domain.model.option.WorkQueryOption
import com.codejsha.platform.shared.data.ActorContext
import com.codejsha.platform.shared.data.ActorType
import kotlinx.coroutines.runBlocking
import org.junit.jupiter.api.Test
import org.springframework.data.domain.Pageable
import org.springframework.web.client.RestClientException
import kotlin.test.assertEquals
import kotlin.test.assertFalse
import kotlin.test.assertTrue

class DashboardServiceTest {

    private val context = ActorContext(actorId = 0L, actorType = ActorType.USER)

    private class CountOnlyCatalogClient(private val works: () -> Long) : CatalogClient {
        override fun countWorks() = works()

        override fun findAllWorks(option: WorkQueryOption, pageable: Pageable) = unsupported()

        override fun findWork(uid: String) = unsupported()

        override fun createWork(command: WorkCreateCommand) = unsupported()

        override fun updateWork(uid: String, command: WorkUpdateCommand) = unsupported()

        override fun findAllAuthors(name: String?, pageable: Pageable) = unsupported()

        override fun findAllSubjects(name: String?, pageable: Pageable) = unsupported()

        private fun unsupported(): Nothing = error("the dashboard only reads counts")
    }

    private class CountOnlyOrderClient(private val orders: () -> Long) : OrderClient {
        override fun countOrders() = orders()

        override fun findAllOrders(option: OrderQueryOption, pageable: Pageable) = unsupported()

        override fun findOrder(uid: String) = unsupported()

        override fun cancelOrder(uid: String) = unsupported()

        private fun unsupported(): Nothing = error("the dashboard only reads counts")
    }

    private class CountOnlyInventoryClient(private val warehouses: () -> Long) : InventoryClient {
        override fun countWarehouses() = warehouses()

        override fun findAllWarehouses(name: String?, pageable: Pageable) = unsupported()

        override fun findWarehouse(uid: String) = unsupported()

        override fun findAllStocks(option: StockQueryOption, pageable: Pageable) = unsupported()

        private fun unsupported(): Nothing = error("the dashboard only reads counts")
    }

    private fun service(
        works: () -> Long = { 1 },
        orders: () -> Long = { 2 },
        warehouses: () -> Long = { 3 },
    ) = DashboardService(
        catalogClient = CountOnlyCatalogClient(works),
        orderClient = CountOnlyOrderClient(orders),
        inventoryClient = CountOnlyInventoryClient(warehouses),
    )

    @Test
    fun `collects a count from each downstream`(): Unit = runBlocking {
        val dashboard = service().loadDashboard(context)

        assertEquals(1, dashboard.works.count)
        assertEquals(2, dashboard.orders.count)
        assertEquals(3, dashboard.warehouses.count)
        assertTrue(dashboard.works.available)
    }

    @Test
    fun `one unreachable downstream degrades only its own metric`(): Unit = runBlocking {
        val dashboard =
            service(orders = { throw RestClientException("connection refused") }).loadDashboard(context)

        assertFalse(dashboard.orders.available)
        assertEquals(0, dashboard.orders.count)
        assertTrue(dashboard.works.available)
        assertEquals(1, dashboard.works.count)
        assertTrue(dashboard.warehouses.available)
    }

    @Test
    fun `an Error is not swallowed as an unavailable metric`(): Unit = runBlocking {
        val svc = service(works = { throw OutOfMemoryError("heap") })

        assertFailsWithOutOfMemory { svc.loadDashboard(context) }
    }

    private suspend fun assertFailsWithOutOfMemory(block: suspend () -> Unit) {
        try {
            block()
        } catch (_: OutOfMemoryError) {
            return
        }
        throw AssertionError("expected OutOfMemoryError to propagate, not degrade to an unavailable metric")
    }
}
