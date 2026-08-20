package com.codejsha.bookstore.settlement.application.port.support

interface TransactionRunner {

    suspend fun <T> tx(block: () -> T): T
}
