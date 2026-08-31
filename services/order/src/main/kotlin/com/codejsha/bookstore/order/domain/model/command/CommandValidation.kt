package com.codejsha.bookstore.order.domain.model.command

import java.math.BigDecimal

class InvalidCommandException(message: String) : RuntimeException(message)

private val CURRENCY_PATTERN = Regex("^[A-Z]{3}$")
private val COUNTRY_PATTERN = Regex("^[A-Z]{2}$")

internal val MAX_MONEY_AMOUNT: BigDecimal = BigDecimal("1E+15")

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

internal fun requireCountry(field: String, value: String) =
    requireCommand(COUNTRY_PATTERN.matches(value)) { "$field must be a 2-letter ISO 3166-1 alpha-2 code, was $value" }

internal fun requireCountryIfPresent(field: String, value: String?) {
    if (value != null) requireCountry(field, value)
}

internal fun requireCurrency(field: String, value: String) =
    requireCommand(CURRENCY_PATTERN.matches(value)) { "$field must be a 3-letter ISO 4217 code, was $value" }

internal fun requireCurrencyIfPresent(field: String, value: String?) {
    if (value != null) requireCurrency(field, value)
}

internal fun requireNonNegative(field: String, value: BigDecimal) =
    requireCommand(value.signum() >= 0) { "$field must not be negative, was $value" }

internal fun requireNonNegativeIfPresent(field: String, value: BigDecimal?) {
    if (value != null) requireNonNegative(field, value)
}

internal fun requireAmountInRange(field: String, value: BigDecimal) {
    requireNonNegative(field, value)
    requireCommand(value < MAX_MONEY_AMOUNT) { "$field must be less than $MAX_MONEY_AMOUNT, was $value" }
}

internal fun requireAmountInRangeIfPresent(field: String, value: BigDecimal?) {
    if (value != null) requireAmountInRange(field, value)
}

internal fun requireAmountMagnitudeInRange(field: String, value: BigDecimal) =
    requireCommand(value.abs() < MAX_MONEY_AMOUNT) {
        "$field magnitude must be less than $MAX_MONEY_AMOUNT, was $value"
    }

internal fun requireRate(field: String, value: BigDecimal, maxScale: Int) {
    requireNonNegative(field, value)
    requireCommand(value <= BigDecimal.ONE) { "$field must not exceed 1, was $value" }
    requireCommand(value.stripTrailingZeros().scale() <= maxScale) {
        "$field must have at most $maxScale decimal places, was $value"
    }
}

internal fun requireRateIfPresent(field: String, value: BigDecimal?, maxScale: Int) {
    if (value != null) requireRate(field, value, maxScale)
}

internal fun requirePositive(field: String, value: Int) =
    requireCommand(value > 0) { "$field must be positive, was $value" }

internal fun requirePositive(field: String, value: Long) =
    requireCommand(value > 0) { "$field must be positive, was $value" }
