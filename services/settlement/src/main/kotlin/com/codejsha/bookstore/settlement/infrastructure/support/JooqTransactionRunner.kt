package com.codejsha.bookstore.settlement.infrastructure.support

import com.codejsha.bookstore.settlement.application.port.support.TransactionRunner
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.withContext
import org.jooq.DSLContext
import org.springframework.stereotype.Component

@Component
class JooqTransactionRunner(private val dsl: DSLContext) : TransactionRunner {

    override suspend fun <T> tx(block: () -> T): T = withContext(Dispatchers.IO) {
        dsl.transactionResult { _ -> block() }
    }
}
