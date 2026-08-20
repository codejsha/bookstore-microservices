package com.codejsha.bookstore.payment.application.port.support

interface TransactionRunner {

    suspend fun <T> tx(block: () -> T): T
}
