package com.codejsha.bookstore.payment.application.port.support

interface TransactionRunner {

    fun <T> tx(block: () -> T): T
}
