package com.codejsha.bookstore.order.application.port.support

interface TransactionRunner {

    fun <T> tx(block: () -> T): T
}
