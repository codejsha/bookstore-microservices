package com.codejsha.bookstore.payment.domain.model.command

class InvalidCommandException(message: String) : RuntimeException(message)

internal const val MAX_IDEMPOTENCY_KEY = 100
internal const val MAX_ERROR_CODE = 64
internal const val MAX_ERROR_MESSAGE = 1024
internal const val MAX_CONNECTOR = 64

private val CURRENCY_PATTERN = Regex("^[A-Z]{3}$")
private val EMAIL_PATTERN = Regex("^[^@\\s]+@[^@\\s]+\\.[^@\\s]+$")
private val CARD_LAST4_PATTERN = Regex("^\\d{4}$")

internal fun requireCommand(condition: Boolean, lazyMessage: () -> String) {
    if (!condition) throw InvalidCommandException(lazyMessage())
}

internal fun requireNonBlank(field: String, value: String) =
    requireCommand(value.isNotBlank()) { "$field must not be blank" }

internal fun requireNonBlankIfPresent(field: String, value: String?) {
    if (value != null) requireNonBlank(field, value)
}

internal fun requireMaxLength(field: String, value: String, max: Int) =
    requireCommand(value.length <= max) { "$field must be at most $max characters, was ${value.length}" }

internal fun requireMaxLengthIfPresent(field: String, value: String?, max: Int) {
    if (value != null) requireMaxLength(field, value, max)
}

internal fun truncate(value: String?, max: Int): String? = value?.take(max)

internal fun requireCurrency(field: String, value: String) =
    requireCommand(CURRENCY_PATTERN.matches(value)) { "$field must be a 3-letter ISO 4217 code, was $value" }

internal fun requireCurrencyIfPresent(field: String, value: String?) {
    if (value != null) requireCurrency(field, value)
}

internal fun requireNonNegative(field: String, value: Long) =
    requireCommand(value >= 0) { "$field must not be negative, was $value" }

internal fun requireNonNegativeIfPresent(field: String, value: Long?) {
    if (value != null) requireNonNegative(field, value)
}

internal fun requirePositive(field: String, value: Long) =
    requireCommand(value > 0) { "$field must be positive, was $value" }

internal fun requireEmailIfPresent(field: String, value: String?) {
    if (value != null) {
        requireCommand(EMAIL_PATTERN.matches(value)) { "$field must be a valid email address" }
    }
}

internal fun requireCardLast4IfPresent(field: String, value: String?) {
    if (value != null) {
        requireCommand(CARD_LAST4_PATTERN.matches(value)) { "$field must be 4 digits" }
    }
}

internal fun requireRangeIfPresent(field: String, value: Int?, min: Int, max: Int) {
    if (value != null) {
        requireCommand(value in min..max) { "$field must be between $min and $max, was $value" }
    }
}
