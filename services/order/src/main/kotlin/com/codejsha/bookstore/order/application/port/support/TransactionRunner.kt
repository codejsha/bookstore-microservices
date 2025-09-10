package com.codejsha.bookstore.order.application.port.support

interface TransactionRunner {

    suspend fun <T> tx(block: () -> T): T
}
