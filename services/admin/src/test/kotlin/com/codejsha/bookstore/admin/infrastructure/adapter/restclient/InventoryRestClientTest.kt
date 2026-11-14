package com.codejsha.bookstore.admin.infrastructure.adapter.restclient

import com.codejsha.bookstore.admin.domain.model.external.Stock
import com.codejsha.bookstore.admin.domain.model.external.StockAtWarehouse
import com.codejsha.bookstore.admin.domain.model.external.Warehouse
import com.codejsha.bookstore.admin.domain.model.option.StockQueryOption
import org.junit.jupiter.api.AfterAll
import org.junit.jupiter.api.AfterEach
import org.junit.jupiter.api.Test
import org.junit.jupiter.api.TestInstance
import org.springframework.data.domain.Pageable
import org.springframework.http.HttpStatus
import org.springframework.web.client.HttpClientErrorException
import org.springframework.web.client.HttpServerErrorException
import org.springframework.web.context.request.RequestContextHolder
import kotlin.test.assertEquals
import kotlin.test.assertFailsWith

@TestInstance(TestInstance.Lifecycle.PER_CLASS)
class InventoryRestClientTest {

    private val downstream = DownstreamStub()

    private val client = InventoryRestClient(
        downstream.config.inventoryWarehouseApi(),
        downstream.config.inventoryStockApi(),
    )

    @AfterEach
    fun clearRequest() {
        RequestContextHolder.resetRequestAttributes()
    }

    @AfterAll
    fun stopDownstream() {
        downstream.stop()
    }

    @Test
    fun `countWarehouses_existingWarehouses_requestsOneItemAndRelaysTotal`() {
        downstream.respond("""{"total":12,"items":[$FULL_WAREHOUSE_JSON]}""")

        val total = client.countWarehouses()

        assertEquals(12, total)
        val request = downstream.onlyRequest()
        assertEquals("GET", request.method)
        assertEquals("/api/v1/warehouses", request.path)
        assertEquals(mapOf("size" to "1"), request.query)
    }

    @Test
    fun `findAllWarehouses_nameFilterWithSortedPage_forwardsParametersAndMapsWarehouses`() {
        downstream.respond("""{"total":112,"items":[$FULL_WAREHOUSE_JSON,$BARE_WAREHOUSE_JSON]}""")

        val page = client.findAllWarehouses("Seoul", SORTED_PAGE)

        val request = downstream.onlyRequest()
        assertEquals("GET", request.method)
        assertEquals("/api/v1/warehouses", request.path)
        assertEquals(mapOf("name" to "Seoul") + SORTED_PAGE_QUERY, request.query)
        assertEquals(112, page.totalElements)
        assertEquals(SORTED_PAGE, page.pageable)
        assertEquals(listOf(FULL_WAREHOUSE, BARE_WAREHOUSE), page.content)
    }

    @Test
    fun `findAllWarehouses_noNameUnpaged_sendsNoQueryParameters`() {
        downstream.respond("""{"total":0,"items":[]}""")

        val page = client.findAllWarehouses(null, Pageable.unpaged())

        assertEquals(emptyMap(), downstream.onlyRequest().query)
        assertEquals(emptyList(), page.content)
    }

    @Test
    fun `findWarehouse_existingUid_readsWarehouseByUid`() {
        downstream.respond(FULL_WAREHOUSE_JSON)

        val warehouse = client.findWarehouse("warehouse-1")

        val request = downstream.onlyRequest()
        assertEquals("GET", request.method)
        assertEquals("/api/v1/warehouses/warehouse-1", request.path)
        assertEquals(FULL_WAREHOUSE, warehouse)
    }

    @Test
    fun `findWarehouse_optionalFieldsAbsent_mapsThemToNull`() {
        downstream.respond(BARE_WAREHOUSE_JSON)

        assertEquals(BARE_WAREHOUSE, client.findWarehouse("warehouse-2"))
    }

    @Test
    fun `findWarehouse_unknownUid_propagatesNotFound`() {
        downstream.respond(PROBLEM_NOT_FOUND, status = 404)

        val failure = assertFailsWith<HttpClientErrorException.NotFound> { client.findWarehouse("missing") }

        assertEquals(PROBLEM_NOT_FOUND, failure.responseBodyAsString)
    }

    @Test
    fun `findWarehouse_inboundBearerToken_replaysItDownstream`() {
        bindInboundAuthorization(CALLER_AUTHORIZATION)
        downstream.respond(FULL_WAREHOUSE_JSON)

        client.findWarehouse("warehouse-1")

        assertEquals(CALLER_AUTHORIZATION, downstream.onlyRequest().authorization)
    }

    @Test
    fun `findAllStocks_editionAndWarehouseWithSortedPage_forwardsParametersAndMapsStocks`() {
        downstream.respond(
            """{"total":90,"items":[
                {
                    "uid":"stock-1",
                    "edition_uid":"edition-1",
                    "total_quantity":42,
                    "warehouses":[
                        {"warehouse_uid":"warehouse-1","warehouse_name":"Seoul DC","quantity":30},
                        {"warehouse_uid":"warehouse-2","warehouse_name":"Busan DC","quantity":12}
                    ]
                },
                {"uid":"stock-2","edition_uid":"edition-2","total_quantity":0,"warehouses":[]}
            ]}""",
        )

        val page = client.findAllStocks(
            StockQueryOption(editionUid = "edition-1", warehouseUid = "warehouse-1"),
            SORTED_PAGE,
        )

        val request = downstream.onlyRequest()
        assertEquals("GET", request.method)
        assertEquals("/api/v1/stocks", request.path)
        assertEquals(
            mapOf("edition_uid" to "edition-1", "warehouse_uid" to "warehouse-1") + SORTED_PAGE_QUERY,
            request.query,
        )
        assertEquals(90, page.totalElements)
        assertEquals(
            listOf(
                Stock(
                    uid = "stock-1",
                    editionUid = "edition-1",
                    totalQuantity = 42,
                    warehouses = listOf(
                        StockAtWarehouse("warehouse-1", "Seoul DC", 30),
                        StockAtWarehouse("warehouse-2", "Busan DC", 12),
                    ),
                ),
                Stock(uid = "stock-2", editionUid = "edition-2", totalQuantity = 0, warehouses = emptyList()),
            ),
            page.content,
        )
    }

    @Test
    fun `findAllStocks_noFilterUnpaged_sendsNoQueryParameters`() {
        downstream.respond("""{"total":0,"items":[]}""")

        val page = client.findAllStocks(StockQueryOption(), Pageable.unpaged())

        assertEquals(emptyMap(), downstream.onlyRequest().query)
        assertEquals(emptyList(), page.content)
    }

    @Test
    fun `findAllStocks_downstreamBadGateway_propagatesServerError`() {
        downstream.respond(status = 502)

        val failure = assertFailsWith<HttpServerErrorException> {
            client.findAllStocks(StockQueryOption(), Pageable.unpaged())
        }

        assertEquals(HttpStatus.BAD_GATEWAY, failure.statusCode)
    }

    private companion object {
        private const val FULL_WAREHOUSE_JSON =
            """{
                "uid":"warehouse-1",
                "name":"Seoul DC",
                "address":"1 Hangang-daero",
                "capacity":5000,
                "created_at":"2026-01-02T03:04:00Z",
                "updated_at":"2026-02-03T04:05:00Z"
            }"""

        private const val BARE_WAREHOUSE_JSON =
            """{"uid":"warehouse-2","name":"Busan DC","capacity":0,"created_at":"2026-03-04T05:06:00Z"}"""

        private val FULL_WAREHOUSE = Warehouse(
            uid = "warehouse-1",
            name = "Seoul DC",
            address = "1 Hangang-daero",
            capacity = 5000,
            createdAt = utc(2026, 1, 2, 3, 4),
            updatedAt = utc(2026, 2, 3, 4, 5),
        )

        private val BARE_WAREHOUSE = Warehouse(
            uid = "warehouse-2",
            name = "Busan DC",
            address = null,
            capacity = 0,
            createdAt = utc(2026, 3, 4, 5, 6),
            updatedAt = null,
        )
    }
}
