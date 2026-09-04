package com.codejsha.bookstore.payment.domain.model.command

import kotlin.test.Test
import kotlin.test.assertFailsWith

class CommandValidationTest {

    private fun paymentCreate(
        amount: Long = 10000L,
        currency: String = "KRW",
        idempotencyKey: String? = "idem-1",
    ) = PaymentCreateCommand(
        amount = amount, currency = currency, customerId = null, paymentMethod = null,
        paymentMethodType = null, authenticationType = null, setupFutureUsage = null,
        description = null, returnUrl = null, billingAddress = null, shippingAddress = null,
        metadata = null, idempotencyKey = idempotencyKey,
    )

    @Test
    fun `paymentCreateCommand_whenAmountCurrencyOrKeyInvalid_throwsInvalidCommandException`() {
        paymentCreate()
        paymentCreate(amount = 0L)
        assertFailsWith<InvalidCommandException> { paymentCreate(amount = -1L) }
        assertFailsWith<InvalidCommandException> { paymentCreate(currency = "krw") }
        assertFailsWith<InvalidCommandException> { paymentCreate(idempotencyKey = " ") }
    }

    @Test
    fun `paymentCreateCommand_whenFieldOverColumnWidth_throwsInvalidCommandException`() {
        assertFailsWith<InvalidCommandException> { paymentCreate(idempotencyKey = "k".repeat(101)) }
        paymentCreate(idempotencyKey = "k".repeat(100))
        assertFailsWith<InvalidCommandException> {
            PaymentCreateCommand(
                amount = 100L, currency = "KRW", customerId = "c".repeat(65), paymentMethod = null,
                paymentMethodType = null, authenticationType = null, setupFutureUsage = null,
                description = null, returnUrl = null, billingAddress = null, shippingAddress = null,
                metadata = null, idempotencyKey = "idem-1",
            )
        }
        assertFailsWith<InvalidCommandException> {
            PaymentCreateCommand(
                amount = 100L, currency = "KRW", customerId = null, paymentMethod = null,
                paymentMethodType = null, authenticationType = null, setupFutureUsage = null,
                description = "d".repeat(501), returnUrl = null, billingAddress = null, shippingAddress = null,
                metadata = null, idempotencyKey = "idem-1",
            )
        }
        assertFailsWith<InvalidCommandException> {
            PaymentCreateCommand(
                amount = 100L, currency = "KRW", customerId = null, paymentMethod = null,
                paymentMethodType = null, authenticationType = null, setupFutureUsage = null,
                description = null, returnUrl = "u".repeat(1025), billingAddress = null, shippingAddress = null,
                metadata = null, idempotencyKey = "idem-1",
            )
        }
        assertFailsWith<InvalidCommandException> {
            PaymentCreateCommand(
                amount = 100L, currency = "KRW", customerId = null, paymentMethod = null,
                paymentMethodType = null, authenticationType = null, setupFutureUsage = null,
                description = null, returnUrl = null, billingAddress = null, shippingAddress = null,
                metadata = null, idempotencyKey = "idem-1", paymentId = "p".repeat(65),
            )
        }
    }

    @Test
    fun `refundCreateCommand_whenKeyOrGatewayIdOverColumnWidth_throwsInvalidCommandException`() {
        assertFailsWith<InvalidCommandException> {
            RefundCreateCommand(
                paymentId = "pay_1", amount = 500L, currency = "KRW",
                reason = null, refundType = null, metadata = null,
                idempotencyKey = "k".repeat(101),
            )
        }
        assertFailsWith<InvalidCommandException> {
            RefundCreateCommand(
                paymentId = "pay_1", amount = 500L, currency = "KRW",
                reason = "r".repeat(501), refundType = null, metadata = null,
            )
        }
        assertFailsWith<InvalidCommandException> {
            RefundCreateCommand(
                paymentId = "pay_1", amount = 500L, currency = "KRW",
                reason = null, refundType = null, metadata = null, refundId = "r".repeat(65),
            )
        }
    }

    @Test
    fun `customerCreateCommand_whenFieldOverColumnWidth_throwsInvalidCommandException`() {
        assertFailsWith<InvalidCommandException> {
            CustomerCreateCommand(
                customerId = "c".repeat(65), name = null, email = null, phone = null,
                phoneCountryCode = null, description = null, metadata = null,
                defaultBillingAddress = null, defaultShippingAddress = null,
            )
        }
        assertFailsWith<InvalidCommandException> {
            CustomerCreateCommand(
                customerId = "cus_1", name = "n".repeat(256), email = null, phone = null,
                phoneCountryCode = null, description = null, metadata = null,
                defaultBillingAddress = null, defaultShippingAddress = null,
            )
        }
        assertFailsWith<InvalidCommandException> {
            CustomerUpdateCommand(
                name = null, email = null, phone = "p".repeat(33), phoneCountryCode = null,
                description = null, metadata = null,
                defaultBillingAddress = null, defaultShippingAddress = null,
            )
        }
        assertFailsWith<InvalidCommandException> {
            CustomerUpdateCommand(
                name = null, email = null, phone = null, phoneCountryCode = "+8200000",
                description = "d".repeat(501), metadata = null,
                defaultBillingAddress = null, defaultShippingAddress = null,
            )
        }
    }

    @Test
    fun `refundCreateCommand_whenAmountNotPositive_throwsInvalidCommandException`() {
        RefundCreateCommand(
            paymentId = "pay_1", amount = 500L, currency = "KRW",
            reason = null, refundType = null, metadata = null,
        )
        assertFailsWith<InvalidCommandException> {
            RefundCreateCommand(
                paymentId = "pay_1", amount = 0L, currency = "KRW",
                reason = null, refundType = null, metadata = null,
            )
        }
        assertFailsWith<InvalidCommandException> {
            RefundCreateCommand(
                paymentId = " ", amount = 500L, currency = "KRW",
                reason = null, refundType = null, metadata = null,
            )
        }
    }

    @Test
    fun `paymentAttemptCommand_whenIdOrStatusMissing_throwsInvalidCommandException`() {
        PaymentAttemptCreateCommand(paymentId = "pay_1", amount = 100L, currency = "USD", status = "succeeded")
        assertFailsWith<InvalidCommandException> {
            PaymentAttemptCreateCommand(paymentId = "pay_1", amount = 100L, currency = "USD", status = " ")
        }
    }

    // Setup and revoke share the rule set, so both commands run here.
    @Test
    fun `mandateCommands_whenRequiredFieldMissing_throwInvalidCommandException`() {
        MandateSetupCommand(customerId = "cus_1", paymentMethodToken = "tok_1", currency = "KRW")
        assertFailsWith<InvalidCommandException> {
            MandateSetupCommand(customerId = "cus_1", paymentMethodToken = " ", currency = "KRW")
        }
        assertFailsWith<InvalidCommandException> {
            MandateSetupCommand(customerId = "cus_1", paymentMethodToken = "tok_1", currency = "KRW", mandateAmountMinor = -1L)
        }
        assertFailsWith<InvalidCommandException> {
            MandateCreateCommand(
                mandateId = "", customerId = "cus_1", paymentMethodId = null,
                mandateType = "multi_use", mandateStatus = "active", mandateAmount = null,
                mandateCurrency = null, setupFutureUsage = null, customerAcceptanceType = null,
            )
        }
    }

    // Create and update share the email rule, so both commands run here.
    @Test
    fun `customerCommands_whenEmailMalformed_throwInvalidCommandException`() {
        CustomerCreateCommand(
            customerId = "cus_1", name = "Alice", email = "a@b.co", phone = null,
            phoneCountryCode = null, description = null, metadata = null,
            defaultBillingAddress = null, defaultShippingAddress = null,
        )
        assertFailsWith<InvalidCommandException> {
            CustomerCreateCommand(
                customerId = "cus_1", name = null, email = "nope", phone = null,
                phoneCountryCode = null, description = null, metadata = null,
                defaultBillingAddress = null, defaultShippingAddress = null,
            )
        }
        assertFailsWith<InvalidCommandException> {
            CustomerUpdateCommand(
                name = " ", email = null, phone = null, phoneCountryCode = null,
                description = null, metadata = null,
                defaultBillingAddress = null, defaultShippingAddress = null,
            )
        }
    }

    @Test
    fun `webhookEventCommand_whenIdentifierMissing_throwsInvalidCommandException`() {
        WebhookEventCreateCommand(
            eventId = "evt_1", eventType = "payment_succeeded", objectType = "payment",
            objectId = "pay_1", payload = "{}", signature = null,
        )
        assertFailsWith<InvalidCommandException> {
            WebhookEventCreateCommand(
                eventId = " ", eventType = "payment_succeeded", objectType = "payment",
                objectId = "pay_1", payload = "{}", signature = null,
            )
        }
    }

    // Create and update share the card rules, so both commands run here.
    @Test
    fun `paymentMethodCommands_whenCardFieldInvalid_throwInvalidCommandException`() {
        PaymentMethodCreateCommand(
            paymentMethod = "card", paymentMethodType = "credit", paymentMethodIssuer = null,
            cardNetwork = null, cardLast4 = "4242", cardExpMonth = 12, cardExpYear = 2030,
            cardHolderName = null, isDefault = null, metadata = null,
        )
        assertFailsWith<InvalidCommandException> {
            PaymentMethodCreateCommand(
                paymentMethod = "card", paymentMethodType = null, paymentMethodIssuer = null,
                cardNetwork = null, cardLast4 = "42", cardExpMonth = null, cardExpYear = null,
                cardHolderName = null, isDefault = null, metadata = null,
            )
        }
        assertFailsWith<InvalidCommandException> {
            PaymentMethodUpdateCommand(
                paymentMethod = null, paymentMethodType = null, paymentMethodIssuer = null,
                cardNetwork = null, cardLast4 = null, cardExpMonth = 13, cardExpYear = null,
                cardHolderName = null, isDefault = null, metadata = null,
            )
        }
    }
}
