package com.codejsha.bookstore.settlement.infrastructure.support

import com.codejsha.bookstore.settlement.application.port.support.TransactionRunner
import org.jooq.DSLContext
import org.springframework.stereotype.Component

@Component
class JooqTransactionRunner(private val dsl: DSLContext) : TransactionRunner {

    override fun <T> tx(block: () -> T): T = dsl.transactionResult { _ -> block() }
}
