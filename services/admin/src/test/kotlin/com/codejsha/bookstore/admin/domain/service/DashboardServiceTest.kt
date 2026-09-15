package com.codejsha.bookstore.admin.domain.service

import com.codejsha.bookstore.admin.application.port.restclient.CatalogClient
import com.codejsha.bookstore.admin.application.port.restclient.InventoryClient
import com.codejsha.bookstore.admin.application.port.restclient.OrderClient
import com.codejsha.bookstore.admin.application.port.support.ConcurrentCallContext
import com.codejsha.bookstore.admin.domain.model.command.WorkCreateCommand
import com.codejsha.bookstore.admin.domain.model.command.WorkUpdateCommand
import com.codejsha.bookstore.admin.domain.model.option.OrderQueryOption
import com.codejsha.bookstore.admin.domain.model.option.StockQueryOption
import com.codejsha.bookstore.admin.domain.model.option.WorkQueryOption
import com.codejsha.platform.shared.data.ActorContext
import com.codejsha.platform.shared.data.ActorType
import kotlinx.coroutines.Dispatchers
import org.junit.jupiter.api.Test
import org.springframework.context.annotation.AnnotationConfigApplicationContext
import org.springframework.data.domain.Pageable
import org.springframework.web.client.RestClientException
import java.time.Duration
import java.util.concurrent.CountDownLatch
import java.util.concurrent.TimeUnit
import java.util.function.Supplier
import kotlin.coroutines.CoroutineContext
import kotlin.test.assertEquals
import kotlin.test.assertFalse
import kotlin.test.assertTrue

class DashboardServiceTest {

    private val context = ActorContext(actorId = 0L, actorType = ActorType.USER)

    private val callContext = object : ConcurrentCallContext {
        override fun capture(): CoroutineContext = Dispatchers.IO
    }

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
        metricTimeout: Duration = Duration.ofSeconds(5),
    ) = DashboardService(
        catalogClient = CountOnlyCatalogClient(works),
        orderClient = CountOnlyOrderClient(orders),
        inventoryClient = CountOnlyInventoryClient(warehouses),
        callContext = callContext,
        metricTimeout = metricTimeout,
    )

    private fun rendezvous(inFlight: CountDownLatch, value: Long): () -> Long = {
        inFlight.countDown()
        check(inFlight.await(2, TimeUnit.SECONDS)) { "downstream calls did not overlap" }
        value
    }

    @Test
    fun `collects a count from each downstream`() {
        val dashboard = service().loadDashboard(context)

        assertEquals(1, dashboard.works.count)
        assertEquals(2, dashboard.orders.count)
        assertEquals(3, dashboard.warehouses.count)
        assertTrue(dashboard.works.available)
    }

    @Test
    fun `one unreachable downstream degrades only its own metric`() {
        val dashboard =
            service(orders = { throw RestClientException("connection refused") }).loadDashboard(context)

        assertFalse(dashboard.orders.available)
        assertEquals(0, dashboard.orders.count)
        assertTrue(dashboard.works.available)
        assertEquals(1, dashboard.works.count)
        assertTrue(dashboard.warehouses.available)
    }

    @Test
    fun `an Error is not swallowed as an unavailable metric`() {
        val svc = service(works = { throw OutOfMemoryError("heap") })

        assertFailsWithOutOfMemory { svc.loadDashboard(context) }
    }

    @Test
    fun `loadDashboard_threeDownstreams_callsThemConcurrently`() {
        val inFlight = CountDownLatch(3)
        val dashboard = service(
            works = rendezvous(inFlight, 1),
            orders = rendezvous(inFlight, 2),
            warehouses = rendezvous(inFlight, 3),
        ).loadDashboard(context)

        assertTrue(dashboard.works.available)
        assertTrue(dashboard.orders.available)
        assertTrue(dashboard.warehouses.available)
        assertEquals(listOf(1L, 2L, 3L), listOf(dashboard.works.count, dashboard.orders.count, dashboard.warehouses.count))
    }

    @Test
    fun `loadDashboard_downstreamHangs_degradesOnlyThatMetricAtTheDeadline`() {
        val interrupted = CountDownLatch(1)
        val svc = service(
            orders = {
                try {
                    Thread.sleep(60_000)
                    2L
                } catch (e: InterruptedException) {
                    interrupted.countDown()
                    throw e
                }
            },
            metricTimeout = Duration.ofMillis(200),
        )

        val startedAt = System.nanoTime()
        val dashboard = svc.loadDashboard(context)
        val elapsed = Duration.ofNanos(System.nanoTime() - startedAt)

        assertFalse(dashboard.orders.available)
        assertTrue(dashboard.works.available)
        assertTrue(dashboard.warehouses.available)
        assertTrue(elapsed < Duration.ofSeconds(5), "expected the metric deadline to cut the hung call, took $elapsed")
        assertTrue(interrupted.await(2, TimeUnit.SECONDS), "the hung downstream call was not interrupted")
    }

    @Test
    fun `springContext_noMetricTimeoutBean_createsServiceWithDefaultDeadline`() {
        AnnotationConfigApplicationContext().use { spring ->
            spring.registerBean(CatalogClient::class.java, Supplier<CatalogClient> { CountOnlyCatalogClient { 1 } })
            spring.registerBean(OrderClient::class.java, Supplier<OrderClient> { CountOnlyOrderClient { 2 } })
            spring.registerBean(InventoryClient::class.java, Supplier<InventoryClient> { CountOnlyInventoryClient { 3 } })
            spring.registerBean(ConcurrentCallContext::class.java, Supplier { callContext })
            spring.register(DashboardService::class.java)
            spring.refresh()

            val dashboard = spring.getBean(DashboardService::class.java).loadDashboard(context)

            assertEquals(1, dashboard.works.count)
        }
    }

    private fun assertFailsWithOutOfMemory(block: () -> Unit) {
        try {
            block()
        } catch (_: OutOfMemoryError) {
            return
        }
        throw AssertionError("expected OutOfMemoryError to propagate, not degrade to an unavailable metric")
    }
}
