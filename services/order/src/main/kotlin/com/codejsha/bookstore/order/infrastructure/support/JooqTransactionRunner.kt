package com.codejsha.bookstore.order.infrastructure.support

import com.codejsha.bookstore.order.application.port.support.TransactionRunner
import org.springframework.stereotype.Component
import org.springframework.transaction.PlatformTransactionManager
import org.springframework.transaction.support.TransactionTemplate

@Component
class JooqTransactionRunner(transactionManager: PlatformTransactionManager) : TransactionRunner {

    private val template = TransactionTemplate(transactionManager)

    @Suppress("UNCHECKED_CAST")
    override fun <T> tx(block: () -> T): T = template.execute { block() } as T
}
