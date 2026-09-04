package com.codejsha.bookstore.order.domain.model.command

import java.math.BigDecimal
import java.util.UUID
import kotlin.test.Test
import kotlin.test.assertEquals
import kotlin.test.assertFailsWith

class CommandValidationTest {

    private fun orderCreate(
        currency: String = "KRW",
        itemsAmount: BigDecimal = BigDecimal("10000.00"),
        discountAmount: BigDecimal = BigDecimal.ZERO,
        totalAmount: BigDecimal = BigDecimal("10000.00"),
        idempotencyKey: String = "idem-1",
    ) = OrderCreateCommand(
        userUid = UUID.randomUUID(),
        currency = currency,
        itemsAmount = itemsAmount,
        discountAmount = discountAmount,
        shippingAmount = BigDecimal.ZERO,
        taxAmount = BigDecimal.ZERO,
        totalAmount = totalAmount,
        idempotencyKey = idempotencyKey,
    )

    @Test
    fun `orderCreateCommand_whenEveryFieldValid_passes`() {
        orderCreate()
    }

    @Test
    fun `orderCreateCommand_whenCurrencyMalformed_throwsInvalidCommandException`() {
        assertFailsWith<InvalidCommandException> { orderCreate(currency = "krw") }
        assertFailsWith<InvalidCommandException> { orderCreate(currency = "KRWX") }
    }

    @Test
    fun `orderCreateCommand_whenAmountNegative_throwsInvalidCommandException`() {
        assertFailsWith<InvalidCommandException> { orderCreate(itemsAmount = BigDecimal("-1")) }
        assertFailsWith<InvalidCommandException> { orderCreate(discountAmount = BigDecimal("-1")) }
    }

    @Test
    fun `orderCreateCommand_whenIdempotencyKeyBlank_throwsInvalidCommandException`() {
        assertFailsWith<InvalidCommandException> { orderCreate(idempotencyKey = " ") }
    }

    @Test
    fun `orderUpdateCommand_whenStatusUnknown_throwsInvalidCommandException`() {
        assertFailsWith<InvalidCommandException> {
            OrderUpdateCommand(
                userUid = null, status = "NOPE", currency = null, itemsAmount = null,
                discountAmount = null, shippingAmount = null, taxAmount = null, totalAmount = null,
            )
        }
    }

    @Test
    fun `orderUpdateCommand_whenEveryFieldNull_passes`() {
        OrderUpdateCommand(
            userUid = null, status = null, currency = null, itemsAmount = null,
            discountAmount = null, shippingAmount = null, taxAmount = null, totalAmount = null,
        )
    }

    private fun orderItemCreate(
        productId: Long = 200L,
        quantity: Int = 1,
        price: BigDecimal = BigDecimal("5000.00"),
        taxRate: BigDecimal = BigDecimal.ZERO,
    ) = OrderItemCreateCommand(
        productId = productId, sku = null, productName = "Book", options = null,
        quantity = quantity, currency = "KRW", price = price, taxRate = taxRate,
    )

    @Test
    fun `orderItemCreateCommand_whenTextOverLimitOrTaxRateOutOfRange_throwsInvalidCommandException`() {
        val valid = OrderItemCreateCommand(
            productId = 200L, sku = "SKU-1", productName = "Book", options = "gift-wrap",
            quantity = 1, currency = "KRW", price = BigDecimal("5000.00"), taxRate = BigDecimal("0.1000"),
        )
        assertFailsWith<InvalidCommandException> { valid.copy(sku = "s".repeat(65)) }
        assertFailsWith<InvalidCommandException> { valid.copy(productName = "p".repeat(256)) }
        assertFailsWith<InvalidCommandException> { valid.copy(options = "o".repeat(256)) }
        assertFailsWith<InvalidCommandException> { valid.copy(options = " ") }
        assertFailsWith<InvalidCommandException> { valid.copy(taxRate = BigDecimal("1.0001")) }
        assertFailsWith<InvalidCommandException> { valid.copy(taxRate = BigDecimal("0.00001")) }
    }

    @Test
    fun `orderCreateCommand_whenAmountAtColumnCeiling_throwsInvalidCommandException`() {
        assertFailsWith<InvalidCommandException> { orderCreate(itemsAmount = BigDecimal("1E+15")) }
        assertFailsWith<InvalidCommandException> { orderCreate(totalAmount = BigDecimal("1000000000000000")) }
        orderCreate(itemsAmount = BigDecimal("999999999999999.9999"), totalAmount = BigDecimal("999999999999999.9999"))
    }

    @Test
    fun `orderCreateCommand_whenIdempotencyKeyOverColumnWidth_throwsInvalidCommandException`() {
        assertFailsWith<InvalidCommandException> { orderCreate(idempotencyKey = "k".repeat(101)) }
        orderCreate(idempotencyKey = "k".repeat(100))
    }

    // The composed key (user uid prefix + client key) must still fit the 100-char column.
    @Test
    fun `cartCheckoutCommand_whenKeyLeavesNoRoomForUserUidPrefix_throwsInvalidCommandException`() {
        val key = "k".repeat(CHECKOUT_IDEMPOTENCY_KEY_MAX_LENGTH)
        CartCheckoutCommand(currency = "KRW", idempotencyKey = key, shipping = null)
        assertFailsWith<InvalidCommandException> {
            CartCheckoutCommand(currency = "KRW", idempotencyKey = "k".repeat(64), shipping = null)
        }
        val composed = "${UUID.randomUUID()}:$key"
        assertEquals(100, composed.length)
    }

    @Test
    fun `orderAdjustmentCreateCommand_whenLabelOverLimitOrAmountOutOfRange_throwsInvalidCommandException`() {
        val valid = OrderAdjustmentCreateCommand(
            type = "COUPON", label = "WELCOME", amount = BigDecimal("-1000.00"), meta = null,
        )
        assertFailsWith<InvalidCommandException> { valid.copy(label = "l".repeat(101)) }
        assertFailsWith<InvalidCommandException> { valid.copy(amount = BigDecimal("-1E+15")) }
        assertFailsWith<InvalidCommandException> { valid.copy(amount = BigDecimal("1E+15")) }
    }

    @Test
    fun `cartAddItemCommand_whenProductNameOverLimit_throwsInvalidCommandException`() {
        val valid = CartAddItemCommand(
            productId = 200L, productName = "Book", quantity = 1,
            currency = "KRW", price = BigDecimal("5000.00"),
        )
        assertFailsWith<InvalidCommandException> { valid.copy(productName = "p".repeat(256)) }
        assertFailsWith<InvalidCommandException> { valid.copy(price = BigDecimal("1E+15")) }
    }

    @Test
    fun `orderItemCreateCommand_whenIdQuantityOrPriceInvalid_throwsInvalidCommandException`() {
        orderItemCreate()
        assertFailsWith<InvalidCommandException> { orderItemCreate(productId = 0L) }
        assertFailsWith<InvalidCommandException> { orderItemCreate(quantity = 0) }
        assertFailsWith<InvalidCommandException> { orderItemCreate(price = BigDecimal("-1")) }
        assertFailsWith<InvalidCommandException> { orderItemCreate(taxRate = BigDecimal("-0.1")) }
    }

    private fun shippingCreate() = OrderShippingCreateCommand(
        recipientName = "Alice", recipientPhone = "010-1", addressLine1 = "1 Foo St",
        addressLine2 = null, city = "Seoul", state = "Seoul", postalCode = "12345",
        country = "KR", shippingMethod = "STANDARD",
    )

    @Test
    fun `orderShippingCreateCommand_whenRequiredFieldBlank_throwsInvalidCommandException`() {
        val valid = shippingCreate()
        assertFailsWith<InvalidCommandException> { valid.copy(recipientName = " ") }
        assertFailsWith<InvalidCommandException> { valid.copy(postalCode = "") }
        assertFailsWith<InvalidCommandException> { valid.copy(state = " ") }
        assertFailsWith<InvalidCommandException> { valid.copy(country = " ") }
        assertFailsWith<InvalidCommandException> { valid.copy(addressLine2 = " ") }
    }

    @Test
    fun `orderShippingCreateCommand_whenCountryNotAlpha2_throwsInvalidCommandException`() {
        val valid = shippingCreate()
        assertFailsWith<InvalidCommandException> { valid.copy(country = "kr") }
        assertFailsWith<InvalidCommandException> { valid.copy(country = "KOR") }
        assertFailsWith<InvalidCommandException> { valid.copy(country = "K") }
    }

    @Test
    fun `orderShippingCreateCommand_whenFieldOverLimit_throwsInvalidCommandException`() {
        val valid = shippingCreate()
        assertFailsWith<InvalidCommandException> { valid.copy(recipientName = "a".repeat(256)) }
        assertFailsWith<InvalidCommandException> { valid.copy(recipientPhone = "0".repeat(21)) }
        assertFailsWith<InvalidCommandException> { valid.copy(addressLine1 = "a".repeat(256)) }
        assertFailsWith<InvalidCommandException> { valid.copy(addressLine2 = "a".repeat(256)) }
        assertFailsWith<InvalidCommandException> { valid.copy(city = "a".repeat(101)) }
        assertFailsWith<InvalidCommandException> { valid.copy(state = "a".repeat(101)) }
        assertFailsWith<InvalidCommandException> { valid.copy(postalCode = "1".repeat(21)) }
        assertFailsWith<InvalidCommandException> { valid.copy(shippingMethod = "a".repeat(101)) }
    }

    @Test
    fun `orderShippingUpdateCommand_whenFieldPresentAndOutOfBounds_throwsInvalidCommandException`() {
        val valid = OrderShippingUpdateCommand(
            recipientName = null, recipientPhone = null, addressLine1 = null, addressLine2 = null,
            city = null, state = null, postalCode = null, country = null, shippingMethod = null,
        )
        assertFailsWith<InvalidCommandException> { valid.copy(state = " ") }
        assertFailsWith<InvalidCommandException> { valid.copy(state = "a".repeat(101)) }
        assertFailsWith<InvalidCommandException> { valid.copy(country = "kr") }
        assertFailsWith<InvalidCommandException> { valid.copy(addressLine2 = " ") }
    }

    @Test
    fun `cartAddItemCommand_whenIdQuantityOrPriceInvalid_throwsInvalidCommandException`() {
        val valid = CartAddItemCommand(
            productId = 200L, productName = "Book", quantity = 1,
            currency = "KRW", price = BigDecimal("5000.00"),
        )
        assertFailsWith<InvalidCommandException> { valid.copy(quantity = 0) }
        assertFailsWith<InvalidCommandException> { valid.copy(productId = -1L) }
        assertFailsWith<InvalidCommandException> { valid.copy(currency = "won") }
    }

    @Test
    fun `cartCheckoutCommand_whenIdempotencyKeyBlank_throwsInvalidCommandException`() {
        CartCheckoutCommand(currency = "KRW", idempotencyKey = "k-1", shipping = null)
        assertFailsWith<InvalidCommandException> {
            CartCheckoutCommand(currency = "KRW", idempotencyKey = " ", shipping = null)
        }
    }

    @Test
    fun `orderAdjustmentCreateCommand_whenTypeUnknown_throwsInvalidCommandException`() {
        OrderAdjustmentCreateCommand(type = "COUPON", label = "WELCOME", amount = BigDecimal("-1000.00"), meta = null)
        assertFailsWith<InvalidCommandException> {
            OrderAdjustmentCreateCommand(type = "BOGUS", label = null, amount = BigDecimal.ONE, meta = null)
        }
    }
}
