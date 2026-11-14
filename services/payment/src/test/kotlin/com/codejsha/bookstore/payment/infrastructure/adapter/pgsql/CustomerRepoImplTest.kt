package com.codejsha.bookstore.payment.infrastructure.adapter.pgsql

import com.codejsha.bookstore.generated.infrastructure.adapter.jooq.tables.references.CUSTOMERS
import com.codejsha.bookstore.payment.domain.model.command.CustomerCreateCommand
import com.codejsha.bookstore.payment.domain.model.command.CustomerUpdateCommand
import com.codejsha.bookstore.payment.domain.model.option.CustomerQueryOption
import com.codejsha.bookstore.payment.infrastructure.support.AbstractPostgresTestSupport
import com.codejsha.bookstore.payment.infrastructure.support.PersistenceManager
import com.codejsha.bookstore.payment.support.PaymentTestFixtures
import org.jooq.DSLContext
import org.junit.jupiter.api.BeforeEach
import org.junit.jupiter.api.Test
import org.springframework.beans.factory.annotation.Autowired
import org.springframework.boot.jooq.test.autoconfigure.JooqTest
import org.springframework.context.annotation.Import
import org.springframework.dao.DataIntegrityViolationException
import org.springframework.dao.DuplicateKeyException
import org.springframework.data.domain.PageRequest
import org.springframework.data.domain.Sort
import org.springframework.transaction.annotation.Propagation
import org.springframework.transaction.annotation.Transactional
import java.util.UUID
import kotlin.test.assertEquals
import kotlin.test.assertFailsWith
import kotlin.test.assertIs
import kotlin.test.assertNotNull
import kotlin.test.assertNull
import kotlin.test.assertTrue

@JooqTest(
    properties = [
        "spring.cloud.config.enabled=false",
        "spring.cloud.config.import-check.enabled=false",
    ],
)
@Import(PersistenceManager::class, CustomerRepoImpl::class)
@Transactional(propagation = Propagation.NOT_SUPPORTED)
class CustomerRepoImplTest : AbstractPostgresTestSupport() {

    @Autowired
    private lateinit var dsl: DSLContext

    @Autowired
    private lateinit var customerRepo: CustomerRepoImpl

    private val ctx = PaymentTestFixtures.DEFAULT_CONTEXT

    @BeforeEach
    fun cleanTables() {
        dsl.deleteFrom(CUSTOMERS).execute()
    }

    @Test
    fun `create_fullCommand_persistsScalarsAndJsonbColumns`() {
        val created = customerRepo.create(
            createCommand(
                customerId = "cus_full",
                metadata = mapOf("tier" to "gold"),
                defaultBillingAddress = mapOf("city" to "Seoul"),
                defaultShippingAddress = mapOf("city" to "Busan"),
            ),
            ctx,
        )

        val found = customerRepo.findOne(created.uid, ctx)
        assertTrue(found.id > 0L)
        assertEquals(created.uid, found.uid)
        assertEquals("cus_full", found.customerId)
        assertEquals("Alice", found.name)
        assertEquals("alice@example.com", found.email)
        assertEquals("01012345678", found.phone)
        assertEquals("+82", found.phoneCountryCode)
        assertEquals("vip", found.description)
        assertEquals("gold", found.metadata?.get("tier"))
        assertEquals("Seoul", found.defaultBillingAddress?.get("city"))
        assertEquals("Busan", found.defaultShippingAddress?.get("city"))
        assertNotNull(found.createdAt)
        assertNull(found.updatedAt)
    }

    @Test
    fun `create_nullOptionalFields_storesNulls`() {
        val created = customerRepo.create(
            CustomerCreateCommand(
                customerId = "cus_min",
                name = null,
                email = null,
                phone = null,
                phoneCountryCode = null,
                description = null,
                metadata = null,
                defaultBillingAddress = null,
                defaultShippingAddress = null,
            ),
            ctx,
        )

        assertEquals("cus_min", created.customerId)
        assertNull(created.name)
        assertNull(created.email)
        assertNull(created.metadata)
        assertNull(created.defaultBillingAddress)
        assertNull(created.defaultShippingAddress)
    }

    @Test
    fun `create_duplicateCustomerId_throwsDuplicateKeyException`() {
        customerRepo.create(createCommand(customerId = "cus_dup"), ctx)

        val ex = assertFailsWith<DataIntegrityViolationException> {
            customerRepo.create(createCommand(customerId = "cus_dup"), ctx)
        }

        assertIs<DuplicateKeyException>(ex)
        assertEquals(1, dsl.fetchCount(CUSTOMERS))
    }

    @Test
    fun `findOne_unknownUid_throwsNoSuchElementException`() {
        assertFailsWith<NoSuchElementException> { customerRepo.findOne(UUID.randomUUID(), ctx) }
    }

    @Test
    fun `findByCustomerId_existingAndMissingIds_findsOnlyStoredCustomer`() {
        val created = customerRepo.create(createCommand(customerId = "cus_lookup"), ctx)

        assertEquals(created.uid, assertNotNull(customerRepo.findByCustomerId("cus_lookup", ctx)).uid)
        assertNull(customerRepo.findByCustomerId("cus_missing", ctx))
    }

    @Test
    fun `findAll_emailFilter_matchesOnlyThatEmail`() {
        customerRepo.create(createCommand(customerId = "cus_a", email = "a@example.com"), ctx)
        customerRepo.create(createCommand(customerId = "cus_b", email = "b@example.com"), ctx)

        val page = customerRepo.findAll(CustomerQueryOption(email = "b@example.com"), PageRequest.of(0, 10), ctx)

        assertEquals(1L, page.totalElements)
        assertEquals(listOf("cus_b"), page.content.map { it.customerId })
    }

    @Test
    fun `findAll_customerIdFilter_matchesOnlyThatCustomer`() {
        customerRepo.create(createCommand(customerId = "cus_a"), ctx)
        customerRepo.create(createCommand(customerId = "cus_b"), ctx)

        val page = customerRepo.findAll(CustomerQueryOption(customerId = "cus_a"), PageRequest.of(0, 10), ctx)

        assertEquals(1L, page.totalElements)
        assertEquals(listOf("cus_a"), page.content.map { it.customerId })
    }

    @Test
    fun `findAll_secondPageSortedByCustomerIdDesc_slicesInOrderAndCountsAll`() {
        listOf("cus_1", "cus_2", "cus_3").forEach { customerRepo.create(createCommand(customerId = it), ctx) }
        val sort = Sort.by(Sort.Direction.DESC, "customer_id")

        val first = customerRepo.findAll(CustomerQueryOption(), PageRequest.of(0, 2, sort), ctx)
        val second = customerRepo.findAll(CustomerQueryOption(), PageRequest.of(1, 2, sort), ctx)

        assertEquals(3L, first.totalElements)
        assertEquals(listOf("cus_3", "cus_2"), first.content.map { it.customerId })
        assertEquals(3L, second.totalElements)
        assertEquals(listOf("cus_1"), second.content.map { it.customerId })
    }

    @Test
    fun `update_partialCommand_changesOnlyGivenFieldsAndStampsUpdatedAt`() {
        val created = customerRepo.create(
            createCommand(customerId = "cus_upd", metadata = mapOf("tier" to "gold")),
            ctx,
        )

        val updated = customerRepo.update(created.uid, updateCommand(name = "Bob"), ctx)

        assertEquals("Bob", updated.name)
        assertEquals(created.email, updated.email)
        assertEquals(created.phone, updated.phone)
        assertEquals("gold", updated.metadata?.get("tier"))
        assertEquals(created.createdAt, updated.createdAt)
        assertNotNull(updated.updatedAt)
    }

    @Test
    fun `update_everyField_overwritesScalarsAndJsonbColumns`() {
        val created = customerRepo.create(createCommand(customerId = "cus_all"), ctx)

        val updated = customerRepo.update(
            created.uid,
            CustomerUpdateCommand(
                name = "Carol",
                email = "carol@example.com",
                phone = "0299998888",
                phoneCountryCode = "+1",
                description = "changed",
                metadata = mapOf("tier" to "silver"),
                defaultBillingAddress = mapOf("city" to "Incheon"),
                defaultShippingAddress = mapOf("city" to "Daegu"),
            ),
            ctx,
        )

        assertEquals("cus_all", updated.customerId)
        assertEquals("Carol", updated.name)
        assertEquals("carol@example.com", updated.email)
        assertEquals("0299998888", updated.phone)
        assertEquals("+1", updated.phoneCountryCode)
        assertEquals("changed", updated.description)
        assertEquals("silver", updated.metadata?.get("tier"))
        assertEquals("Incheon", updated.defaultBillingAddress?.get("city"))
        assertEquals("Daegu", updated.defaultShippingAddress?.get("city"))
    }

    @Test
    fun `update_unknownUid_throwsIllegalStateException`() {
        assertFailsWith<IllegalStateException> {
            customerRepo.update(UUID.randomUUID(), updateCommand(name = "Bob"), ctx)
        }
    }

    @Test
    fun `delete_storedCustomer_keepsRowAndHidesItFromEveryRead`() {
        val created = customerRepo.create(createCommand(customerId = "cus_del"), ctx)

        customerRepo.delete(created.uid, ctx)

        assertEquals(1, dsl.fetchCount(CUSTOMERS, CUSTOMERS.UID.eq(created.uid).and(CUSTOMERS.DELETED_AT.isNotNull)))
        assertFailsWith<NoSuchElementException> { customerRepo.findOne(created.uid, ctx) }
        assertNull(customerRepo.findByCustomerId("cus_del", ctx))
        assertEquals(0L, customerRepo.findAll(CustomerQueryOption(), PageRequest.of(0, 10), ctx).totalElements)
        assertFailsWith<IllegalStateException> { customerRepo.update(created.uid, updateCommand(name = "Bob"), ctx) }
    }

    @Test
    fun `delete_alreadyDeletedCustomer_keepsOriginalDeletedAt`() {
        val created = customerRepo.create(createCommand(customerId = "cus_del_twice"), ctx)
        customerRepo.delete(created.uid, ctx)
        val firstDeletedAt = deletedAtOf(created.uid)

        Thread.sleep(5)
        customerRepo.delete(created.uid, ctx)

        assertEquals(firstDeletedAt, deletedAtOf(created.uid))
    }

    private fun deletedAtOf(uid: UUID) =
        dsl.select(CUSTOMERS.DELETED_AT).from(CUSTOMERS).where(CUSTOMERS.UID.eq(uid)).fetchOne(CUSTOMERS.DELETED_AT)

    private fun createCommand(
        customerId: String,
        email: String? = "alice@example.com",
        metadata: Map<String, Any>? = null,
        defaultBillingAddress: Map<String, Any>? = null,
        defaultShippingAddress: Map<String, Any>? = null,
    ) = CustomerCreateCommand(
        customerId = customerId,
        name = "Alice",
        email = email,
        phone = "01012345678",
        phoneCountryCode = "+82",
        description = "vip",
        metadata = metadata,
        defaultBillingAddress = defaultBillingAddress,
        defaultShippingAddress = defaultShippingAddress,
    )

    private fun updateCommand(name: String?) = CustomerUpdateCommand(
        name = name,
        email = null,
        phone = null,
        phoneCountryCode = null,
        description = null,
        metadata = null,
        defaultBillingAddress = null,
        defaultShippingAddress = null,
    )
}
