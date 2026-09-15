package com.codejsha.bookstore.payment.support

import com.codejsha.bookstore.payment.application.port.support.TransactionRunner

class FakeTransactionRunner : TransactionRunner {
    override fun <T> tx(block: () -> T): T = block()
}
