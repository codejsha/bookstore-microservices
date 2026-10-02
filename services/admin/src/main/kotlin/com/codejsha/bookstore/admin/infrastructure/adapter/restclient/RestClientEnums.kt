package com.codejsha.bookstore.admin.infrastructure.adapter.restclient

import com.codejsha.bookstore.admin.domain.model.UnsupportedValueException

internal fun <E : Enum<E>> parseEnumValue(label: String, value: String, fromValue: (String) -> E): E =
    try {
        fromValue(value)
    } catch (_: NoSuchElementException) {
        throw UnsupportedValueException("unknown $label: $value")
    } catch (_: IllegalArgumentException) {
        throw UnsupportedValueException("unknown $label: $value")
    }
