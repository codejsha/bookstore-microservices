package com.codejsha.bookstore.settlement.application.port.support

interface TransactionRunner {

    fun <T> tx(block: () -> T): T
}
