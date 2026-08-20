package com.codejsha.bookstore.admin.infrastructure.adapter.restclient

import org.junit.jupiter.api.Test
import org.springframework.data.domain.PageRequest
import org.springframework.data.domain.Pageable
import org.springframework.data.domain.Sort
import kotlin.test.assertEquals
import kotlin.test.assertNull

class RestClientPagingTest {

    @Test
    fun `an unpaged request forwards no paging at all`() {
        val pageable = Pageable.unpaged()

        assertNull(pageable.sizeParam())
        assertNull(pageable.pageParam())
        assertNull(pageable.sortParam())
    }

    @Test
    fun `a sorted page forwards size, page and a single sort order`() {
        val pageable = PageRequest.of(2, 25, Sort.by(Sort.Direction.DESC, "title"))

        assertEquals(25, pageable.sizeParam())
        assertEquals(2, pageable.pageParam())
        assertEquals("title,desc", pageable.sortParam())
    }

    @Test
    fun `only the first sort order is forwarded`() {
        val pageable = PageRequest.of(0, 10, Sort.by(Sort.Order.asc("title"), Sort.Order.desc("olKey")))

        assertEquals("title,asc", pageable.sortParam())
    }

    @Test
    fun `a page without a sort forwards no sort parameter`() {
        val pageable = PageRequest.of(0, 10)

        assertEquals(10, pageable.sizeParam())
        assertEquals(0, pageable.pageParam())
        assertNull(pageable.sortParam())
    }

    @Test
    fun `the relayed total survives the page it is wrapped in`() {
        val pageable = PageRequest.of(2, 20)

        val page = pageOf(listOf("a", "b"), pageable, 42)

        assertEquals(42, page.totalElements)
        assertEquals(listOf("a", "b"), page.content)
    }
}
