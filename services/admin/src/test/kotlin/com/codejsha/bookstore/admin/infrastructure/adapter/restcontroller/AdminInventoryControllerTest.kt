package com.codejsha.bookstore.admin.infrastructure.adapter.restcontroller

import com.codejsha.bookstore.admin.application.usecase.InventoryUseCase
import com.codejsha.bookstore.admin.domain.model.external.Stock
import com.codejsha.bookstore.admin.domain.model.external.StockAtWarehouse
import com.codejsha.bookstore.admin.domain.model.external.Warehouse
import com.codejsha.bookstore.admin.domain.model.option.StockQueryOption
import com.codejsha.bookstore.admin.infrastructure.support.auth.ForbiddenException
import com.codejsha.bookstore.admin.infrastructure.support.auth.HttpPrincipalResolver
import com.codejsha.platform.shared.data.ActorContext
import com.codejsha.platform.shared.data.ActorType
import kotlinx.coroutines.runBlocking
import org.junit.jupiter.api.AfterEach
import org.junit.jupiter.api.Test
import org.mockito.BDDMockito.given
import org.mockito.Mockito.mock
import org.mockito.Mockito.verifyNoInteractions
import org.springframework.data.domain.PageImpl
import org.springframework.data.domain.Pageable
import org.springframework.mock.web.MockHttpServletRequest
import org.springframework.web.context.request.RequestContextHolder
import org.springframework.web.context.request.ServletRequestAttributes
import tools.jackson.databind.ObjectMapper
import java.time.OffsetDateTime
import java.time.ZoneOffset
import kotlin.test.assertEquals
import kotlin.test.assertFailsWith
import kotlin.test.assertTrue

class AdminInventoryControllerTest {

    private val resolver = HttpPrincipalResolver(ObjectMapper())
    private val controllerContext = ActorContext(actorId = 0L, actorType = ActorType.USER)

    private fun bindPrincipal(roles: String?) {
        val request = MockHttpServletRequest()
        request.addHeader("X-User-Id", "11111111-1111-1111-1111-111111111111")
        roles?.let { request.addHeader("X-User-Roles", it) }
        RequestContextHolder.setRequestAttributes(ServletRequestAttributes(request))
    }

    @AfterEach
    fun clearPrincipal() {
        RequestContextHolder.resetRequestAttributes()
    }

    @Test
    fun `every inventory endpoint rejects a caller without the ADMIN role`() {
        bindPrincipal(roles = "VIEW")
        val useCase = mock(InventoryUseCase::class.java)
        val controller = AdminInventoryController(useCase, resolver)

        assertFailsWith<ForbiddenException> { controller.adminInventoryListWarehouses(null, null) }
        assertFailsWith<ForbiddenException> { controller.adminInventoryReadWarehouse(WAREHOUSE_UID) }
        assertFailsWith<ForbiddenException> { controller.adminInventoryListStocks(null, null, null) }
        verifyNoInteractions(useCase)
    }

    @Test
    fun `listWarehouses maps the name filter and the warehouse fields`(): Unit = runBlocking {
        bindPrincipal(roles = "ADMIN")
        val useCase = mock(InventoryUseCase::class.java)
        val controller = AdminInventoryController(useCase, resolver)
        val unpaged = Pageable.unpaged()
        given(useCase.findAllWarehouses("Seoul", unpaged, controllerContext))
            .willReturn(PageImpl(listOf(warehouse()), unpaged, 1))

        val body = controller.adminInventoryListWarehouses("Seoul", null).body!!

        assertEquals(1L, body.total)
        val item = body.items.first()
        assertEquals(WAREHOUSE_UID, item.uid)
        assertEquals("Seoul Central", item.name)
        assertEquals(5_000, item.capacity)
    }

    @Test
    fun `listStocks maps the per-warehouse breakdown`(): Unit = runBlocking {
        bindPrincipal(roles = "ADMIN")
        val useCase = mock(InventoryUseCase::class.java)
        val controller = AdminInventoryController(useCase, resolver)
        val unpaged = Pageable.unpaged()
        given(useCase.findAllStocks(StockQueryOption(editionUid = EDITION_UID), unpaged, controllerContext))
            .willReturn(PageImpl(listOf(stock()), unpaged, 1))

        val body = controller.adminInventoryListStocks(EDITION_UID, null, null).body!!

        val item = body.items.first()
        assertEquals(EDITION_UID, item.editionUid)
        assertEquals(12, item.totalQuantity)
        assertEquals(2, item.warehouses.size)
        assertEquals("Seoul Central", item.warehouses.first().warehouseName)
        assertEquals(7, item.warehouses.first().quantity)
    }

    @Test
    fun `an edition held in no warehouse keeps an empty breakdown`(): Unit = runBlocking {
        bindPrincipal(roles = "ADMIN")
        val useCase = mock(InventoryUseCase::class.java)
        val controller = AdminInventoryController(useCase, resolver)
        val unpaged = Pageable.unpaged()
        given(useCase.findAllStocks(StockQueryOption(), unpaged, controllerContext))
            .willReturn(PageImpl(listOf(stock(warehouses = emptyList(), total = 0)), unpaged, 1))

        val body = controller.adminInventoryListStocks(null, null, null).body!!

        assertEquals(0, body.items.first().totalQuantity)
        assertTrue(body.items.first().warehouses.isEmpty())
    }

    private companion object {
        private const val WAREHOUSE_UID = "77777777-7777-7777-7777-777777777777"
        private const val OTHER_WAREHOUSE_UID = "88888888-8888-8888-8888-888888888888"
        private const val EDITION_UID = "99999999-9999-9999-9999-999999999999"
        private val CREATED_AT: OffsetDateTime = OffsetDateTime.of(2026, 2, 1, 0, 0, 0, 0, ZoneOffset.UTC)

        private fun warehouse() = Warehouse(
            uid = WAREHOUSE_UID,
            name = "Seoul Central",
            address = "1 Main Street",
            capacity = 5_000,
            createdAt = CREATED_AT,
            updatedAt = null,
        )

        private fun stock(
            warehouses: List<StockAtWarehouse> = listOf(
                StockAtWarehouse(WAREHOUSE_UID, "Seoul Central", 7),
                StockAtWarehouse(OTHER_WAREHOUSE_UID, "Busan Depot", 5),
            ),
            total: Int = 12,
        ) = Stock(
            uid = "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
            editionUid = EDITION_UID,
            totalQuantity = total,
            warehouses = warehouses,
        )
    }
}
