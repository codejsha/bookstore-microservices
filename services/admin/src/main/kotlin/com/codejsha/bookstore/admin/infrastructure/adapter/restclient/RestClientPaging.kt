package com.codejsha.bookstore.admin.infrastructure.adapter.restclient

import org.springframework.data.domain.Page
import org.springframework.data.domain.PageImpl
import org.springframework.data.domain.Pageable
import org.springframework.data.domain.Sort

internal fun Pageable.sizeParam(): Int? = if (isPaged) pageSize else null

internal fun Pageable.pageParam(): Int? = if (isPaged) pageNumber else null

internal fun Pageable.sortParam(): String? = if (isPaged) sort.toSortParam() else null

internal fun Sort.toSortParam(): String? =
    firstOrNull()?.let { "${it.property},${it.direction.name.lowercase()}" }

internal fun <T : Any> pageOf(items: List<T>, pageable: Pageable, total: Long): Page<T> =
    PageImpl(items, pageable, total)
