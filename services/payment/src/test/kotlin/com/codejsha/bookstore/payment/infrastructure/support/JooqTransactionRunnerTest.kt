package com.codejsha.bookstore.payment.infrastructure.support

import org.junit.jupiter.api.Test
import org.mockito.ArgumentMatchers.any
import org.mockito.BDDMockito.given
import org.mockito.Mockito.mock
import org.mockito.Mockito.verify
import org.mockito.Mockito.verifyNoMoreInteractions
import org.springframework.transaction.PlatformTransactionManager
import org.springframework.transaction.TransactionDefinition
import org.springframework.transaction.support.SimpleTransactionStatus
import kotlin.test.assertEquals
import kotlin.test.assertFailsWith
import kotlin.test.assertNull
import kotlin.test.assertSame

class JooqTransactionRunnerTest {

    private val transactionManager: PlatformTransactionManager = mock(PlatformTransactionManager::class.java)
    private val status = SimpleTransactionStatus()
    private val runner = JooqTransactionRunner(transactionManager)

    @Test
    fun `tx_blockCompletes_commitsAndYieldsBlockValue`() {
        given(transactionManager.getTransaction(any(TransactionDefinition::class.java))).willReturn(status)

        val result = runner.tx { "committed" }

        assertEquals("committed", result)
        verify(transactionManager).getTransaction(any(TransactionDefinition::class.java))
        verify(transactionManager).commit(status)
        verifyNoMoreInteractions(transactionManager)
    }

    @Test
    fun `tx_blockYieldsNull_commitsAndYieldsNull`() {
        given(transactionManager.getTransaction(any(TransactionDefinition::class.java))).willReturn(status)

        val result = runner.tx<String?> { null }

        assertNull(result)
        verify(transactionManager).commit(status)
    }

    @Test
    fun `tx_blockThrows_rollsBackAndPropagatesSameException`() {
        given(transactionManager.getTransaction(any(TransactionDefinition::class.java))).willReturn(status)
        val failure = IllegalStateException("duplicate payment")

        val thrown = assertFailsWith<IllegalStateException> { runner.tx<Unit> { throw failure } }

        assertSame(failure, thrown)
        verify(transactionManager).getTransaction(any(TransactionDefinition::class.java))
        verify(transactionManager).rollback(status)
        verifyNoMoreInteractions(transactionManager)
    }
}
