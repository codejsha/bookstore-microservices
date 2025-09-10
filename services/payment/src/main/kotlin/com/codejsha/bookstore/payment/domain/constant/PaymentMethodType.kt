package com.codejsha.bookstore.payment.domain.constant

enum class PaymentMethodType(
    val value: String
) {
    CARD("card"),
    WALLET("wallet"),
    BANK_TRANSFER("bank_transfer"),
    PAY_LATER("pay_later");

    companion object {
        fun fromValue(value: String): PaymentMethodType =
            entries.first { it.value == value }
    }
}
