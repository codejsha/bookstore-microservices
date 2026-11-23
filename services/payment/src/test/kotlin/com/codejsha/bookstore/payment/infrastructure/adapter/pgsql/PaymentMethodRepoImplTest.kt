package com.codejsha.bookstore.payment.infrastructure.adapter.pgsql

import com.codejsha.bookstore.generated.infrastructure.adapter.jooq.tables.references.CUSTOMERS
import com.codejsha.bookstore.generated.infrastructure.adapter.jooq.tables.references.PAYMENT_METHODS
import com.codejsha.bookstore.payment.domain.model.command.CustomerCreateCommand
import com.codejsha.bookstore.payment.domain.model.command.PaymentMethodCreateCommand
import com.codejsha.bookstore.payment.domain.model.command.PaymentMethodUpdateCommand
import com.codejsha.bookstore.payment.domain.model.option.PaymentMethodQueryOption
import com.codejsha.bookstore.payment.infrastructure.support.AbstractPostgresTestSupport
import com.codejsha.bookstore.payment.infrastructure.support.PersistenceManager
import com.codejsha.bookstore.payment.support.PaymentTestFixtures
import org.jooq.DSLContext
import org.junit.jupiter.api.BeforeEach
import org.junit.jupiter.api.Test
import org.springframework.beans.factory.annotation.Autowired
import org.springframework.boot.jooq.test.autoconfigure.JooqTest
import org.springframework.context.annotation.Import
import org.springframework.data.domain.PageRequest
import org.springframework.data.domain.Sort
import org.springframework.transaction.annotation.Propagation
import org.springframework.transaction.annotation.Transactional
import java.util.UUID
import kotlin.test.assertEquals
import kotlin.test.assertFailsWith
import kotlin.test.assertNotEquals
import kotlin.test.assertNotNull
import kotlin.test.assertNull
import kotlin.test.assertTrue

@JooqTest
@Import(PersistenceManager::class, CustomerRepoImpl::class, PaymentMethodRepoImpl::class)
@Transactional(propagation = Propagation.NOT_SUPPORTED)
class PaymentMethodRepoImplTest : AbstractPostgresTestSupport() {

    @Autowired
    private lateinit var dsl: DSLContext

    @Autowired
    private lateinit var customerRepo: CustomerRepoImpl

    @Autowired
    private lateinit var paymentMethodRepo: PaymentMethodRepoImpl

    private val ctx = PaymentTestFixtures.DEFAULT_CONTEXT

    private lateinit var customerUid: UUID
    private lateinit var otherCustomerUid: UUID

    @BeforeEach
    fun seedCustomers() {
        dsl.deleteFrom(PAYMENT_METHODS).execute()
        dsl.deleteFrom(CUSTOMERS).execute()
        customerUid = customerRepo.create(customerCommand("cus_001"), ctx).uid
        otherCustomerUid = customerRepo.create(customerCommand("cus_002"), ctx).uid
    }

    @Test
    fun `create_fullCommand_persistsCardFieldsUnderResolvedCustomerId`() {
        val created = paymentMethodRepo.create(
            customerUid,
            createCommand(isDefault = true, metadata = mapOf("fingerprint" to "fp_1")),
            ctx,
        )

        val found = paymentMethodRepo.findOne(customerUid, created.uid, ctx)
        assertTrue(found.id > 0L)
        assertEquals(created.uid, found.uid)
        assertTrue(found.paymentMethodId.matches(Regex("pm_[0-9a-f]{24}")), found.paymentMethodId)
        assertEquals("cus_001", found.customerId)
        assertEquals("card", found.paymentMethod)
        assertEquals("credit", found.paymentMethodType)
        assertEquals("Shinhan", found.paymentMethodIssuer)
        assertEquals("visa", found.cardNetwork)
        assertEquals("4242", found.cardLast4)
        assertEquals(12, found.cardExpMonth)
        assertEquals(2030, found.cardExpYear)
        assertEquals("Alice", found.cardHolderName)
        assertEquals(true, found.isDefault)
        assertEquals("fp_1", found.metadata?.get("fingerprint"))
        assertNotNull(found.createdAt)
        assertNull(found.updatedAt)
    }

    @Test
    fun `create_nullOptionalFields_storesNullsAndNonDefault`() {
        val created = paymentMethodRepo.create(
            customerUid,
            PaymentMethodCreateCommand(
                paymentMethod = "wallet",
                paymentMethodType = null,
                paymentMethodIssuer = null,
                cardNetwork = null,
                cardLast4 = null,
                cardExpMonth = null,
                cardExpYear = null,
                cardHolderName = null,
                isDefault = null,
                metadata = null,
            ),
            ctx,
        )

        assertEquals("wallet", created.paymentMethod)
        assertEquals(false, created.isDefault)
        assertNull(created.paymentMethodType)
        assertNull(created.cardLast4)
        assertNull(created.cardExpMonth)
        assertNull(created.cardExpYear)
        assertNull(created.metadata)
    }

    @Test
    fun `create_twoMethodsForOneCustomer_assignsDistinctIdentifiers`() {
        val first = paymentMethodRepo.create(customerUid, createCommand(), ctx)
        val second = paymentMethodRepo.create(customerUid, createCommand(), ctx)

        assertNotEquals(first.uid, second.uid)
        assertNotEquals(first.paymentMethodId, second.paymentMethodId)
        assertEquals(2, dsl.fetchCount(PAYMENT_METHODS))
    }

    @Test
    fun `create_unknownCustomer_throwsNoSuchElementExceptionAndInsertsNothing`() {
        assertFailsWith<NoSuchElementException> {
            paymentMethodRepo.create(UUID.randomUUID(), createCommand(), ctx)
        }

        assertEquals(0, dsl.fetchCount(PAYMENT_METHODS))
    }

    @Test
    fun `create_softDeletedCustomer_throwsNoSuchElementException`() {
        customerRepo.delete(customerUid, ctx)

        assertFailsWith<NoSuchElementException> { paymentMethodRepo.create(customerUid, createCommand(), ctx) }
    }

    @Test
    fun `findOne_methodOwnedByAnotherCustomer_throwsNoSuchElementException`() {
        val created = paymentMethodRepo.create(customerUid, createCommand(), ctx)

        assertFailsWith<NoSuchElementException> { paymentMethodRepo.findOne(otherCustomerUid, created.uid, ctx) }
    }

    @Test
    fun `findOne_unknownUid_throwsNoSuchElementException`() {
        assertFailsWith<NoSuchElementException> { paymentMethodRepo.findOne(customerUid, UUID.randomUUID(), ctx) }
    }

    @Test
    fun `findAllByCustomer_methodsOfTwoCustomers_listsOnlyOwnedOnes`() {
        val own = paymentMethodRepo.create(customerUid, createCommand(), ctx)
        paymentMethodRepo.create(otherCustomerUid, createCommand(), ctx)

        val page = paymentMethodRepo.findAllByCustomer(customerUid, PaymentMethodQueryOption(), PageRequest.of(0, 10), ctx)

        assertEquals(1L, page.totalElements)
        assertEquals(listOf(own.uid), page.content.map { it.uid })
    }

    @Test
    fun `findAllByCustomer_paymentMethodFilter_matchesOnlyThatKind`() {
        paymentMethodRepo.create(customerUid, createCommand(paymentMethod = "card"), ctx)
        val wallet = paymentMethodRepo.create(customerUid, createCommand(paymentMethod = "wallet"), ctx)

        val page = paymentMethodRepo.findAllByCustomer(
            customerUid,
            PaymentMethodQueryOption(paymentMethod = "wallet"),
            PageRequest.of(0, 10),
            ctx,
        )

        assertEquals(1L, page.totalElements)
        assertEquals(listOf(wallet.uid), page.content.map { it.uid })
    }

    @Test
    fun `findAllByCustomer_pagedSortedByCardLast4Desc_slicesInOrderAndCountsAll`() {
        listOf("1111", "3333", "2222").forEach {
            paymentMethodRepo.create(customerUid, createCommand(cardLast4 = it), ctx)
        }
        val sort = Sort.by(Sort.Direction.DESC, "card_last4")

        val first = paymentMethodRepo.findAllByCustomer(customerUid, PaymentMethodQueryOption(), PageRequest.of(0, 2, sort), ctx)
        val second = paymentMethodRepo.findAllByCustomer(customerUid, PaymentMethodQueryOption(), PageRequest.of(1, 2, sort), ctx)

        assertEquals(3L, first.totalElements)
        assertEquals(listOf("3333", "2222"), first.content.map { it.cardLast4 })
        assertEquals(3L, second.totalElements)
        assertEquals(listOf("1111"), second.content.map { it.cardLast4 })
    }

    @Test
    fun `findAllByCustomer_unknownCustomer_throwsNoSuchElementException`() {
        assertFailsWith<NoSuchElementException> {
            paymentMethodRepo.findAllByCustomer(UUID.randomUUID(), PaymentMethodQueryOption(), PageRequest.of(0, 10), ctx)
        }
    }

    @Test
    fun `update_partialCommand_changesOnlyGivenFieldsAndStampsUpdatedAt`() {
        val created = paymentMethodRepo.create(
            customerUid,
            createCommand(isDefault = false, metadata = mapOf("fingerprint" to "fp_1")),
            ctx,
        )

        val updated = paymentMethodRepo.update(customerUid, created.uid, updateCommand(isDefault = true), ctx)

        assertEquals(true, updated.isDefault)
        assertEquals(created.paymentMethodId, updated.paymentMethodId)
        assertEquals("card", updated.paymentMethod)
        assertEquals("4242", updated.cardLast4)
        assertEquals(12, updated.cardExpMonth)
        assertEquals("fp_1", updated.metadata?.get("fingerprint"))
        assertNotNull(updated.updatedAt)
    }

    @Test
    fun `update_everyField_overwritesAllMutableColumns`() {
        val created = paymentMethodRepo.create(customerUid, createCommand(), ctx)

        val updated = paymentMethodRepo.update(
            customerUid,
            created.uid,
            PaymentMethodUpdateCommand(
                paymentMethod = "card_redirect",
                paymentMethodType = "debit",
                paymentMethodIssuer = "Kookmin",
                cardNetwork = "mastercard",
                cardLast4 = "0005",
                cardExpMonth = 1,
                cardExpYear = 2031,
                cardHolderName = "Bob",
                isDefault = true,
                metadata = mapOf("fingerprint" to "fp_2"),
            ),
            ctx,
        )

        assertEquals("card_redirect", updated.paymentMethod)
        assertEquals("debit", updated.paymentMethodType)
        assertEquals("Kookmin", updated.paymentMethodIssuer)
        assertEquals("mastercard", updated.cardNetwork)
        assertEquals("0005", updated.cardLast4)
        assertEquals(1, updated.cardExpMonth)
        assertEquals(2031, updated.cardExpYear)
        assertEquals("Bob", updated.cardHolderName)
        assertEquals(true, updated.isDefault)
        assertEquals("fp_2", updated.metadata?.get("fingerprint"))
        assertEquals("cus_001", updated.customerId)
    }

    @Test
    fun `update_methodOwnedByAnotherCustomer_throwsIllegalStateExceptionAndLeavesRow`() {
        val created = paymentMethodRepo.create(customerUid, createCommand(isDefault = false), ctx)

        assertFailsWith<IllegalStateException> {
            paymentMethodRepo.update(otherCustomerUid, created.uid, updateCommand(isDefault = true), ctx)
        }

        val untouched = paymentMethodRepo.findOne(customerUid, created.uid, ctx)
        assertEquals(false, untouched.isDefault)
        assertNull(untouched.updatedAt)
    }

    @Test
    fun `update_unknownUid_throwsIllegalStateException`() {
        assertFailsWith<IllegalStateException> {
            paymentMethodRepo.update(customerUid, UUID.randomUUID(), updateCommand(isDefault = true), ctx)
        }
    }

    @Test
    fun `delete_storedMethod_keepsRowAndHidesItFromReads`() {
        val created = paymentMethodRepo.create(customerUid, createCommand(), ctx)

        paymentMethodRepo.delete(customerUid, created.uid, ctx)

        assertEquals(
            1,
            dsl.fetchCount(PAYMENT_METHODS, PAYMENT_METHODS.UID.eq(created.uid).and(PAYMENT_METHODS.DELETED_AT.isNotNull)),
        )
        assertFailsWith<NoSuchElementException> { paymentMethodRepo.findOne(customerUid, created.uid, ctx) }
        assertEquals(
            0L,
            paymentMethodRepo.findAllByCustomer(customerUid, PaymentMethodQueryOption(), PageRequest.of(0, 10), ctx).totalElements,
        )
        assertFailsWith<IllegalStateException> {
            paymentMethodRepo.update(customerUid, created.uid, updateCommand(isDefault = true), ctx)
        }
    }

    @Test
    fun `delete_methodOwnedByAnotherCustomer_leavesRowVisible`() {
        val created = paymentMethodRepo.create(customerUid, createCommand(), ctx)

        paymentMethodRepo.delete(otherCustomerUid, created.uid, ctx)

        assertEquals(created.uid, paymentMethodRepo.findOne(customerUid, created.uid, ctx).uid)
    }

    @Test
    fun `delete_unknownCustomer_throwsNoSuchElementException`() {
        val created = paymentMethodRepo.create(customerUid, createCommand(), ctx)

        assertFailsWith<NoSuchElementException> { paymentMethodRepo.delete(UUID.randomUUID(), created.uid, ctx) }
    }

    private fun customerCommand(customerId: String) = CustomerCreateCommand(
        customerId = customerId,
        name = "Alice",
        email = "alice@example.com",
        phone = null,
        phoneCountryCode = null,
        description = null,
        metadata = null,
        defaultBillingAddress = null,
        defaultShippingAddress = null,
    )

    private fun createCommand(
        paymentMethod: String = "card",
        cardLast4: String? = "4242",
        isDefault: Boolean? = null,
        metadata: Map<String, Any>? = null,
    ) = PaymentMethodCreateCommand(
        paymentMethod = paymentMethod,
        paymentMethodType = "credit",
        paymentMethodIssuer = "Shinhan",
        cardNetwork = "visa",
        cardLast4 = cardLast4,
        cardExpMonth = 12,
        cardExpYear = 2030,
        cardHolderName = "Alice",
        isDefault = isDefault,
        metadata = metadata,
    )

    private fun updateCommand(isDefault: Boolean?) = PaymentMethodUpdateCommand(
        paymentMethod = null,
        paymentMethodType = null,
        paymentMethodIssuer = null,
        cardNetwork = null,
        cardLast4 = null,
        cardExpMonth = null,
        cardExpYear = null,
        cardHolderName = null,
        isDefault = isDefault,
        metadata = null,
    )
}
