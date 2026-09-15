package com.codejsha.bookstore.payment.infrastructure.support

import com.codejsha.bookstore.payment.application.port.support.TransactionRunner
import org.jooq.DSLContext
import org.springframework.stereotype.Component

@Component
class JooqTransactionRunner(private val dsl: DSLContext) : TransactionRunner {

    override fun <T> tx(block: () -> T): T = dsl.transactionResult { _ -> block() }
}
