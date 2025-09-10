package com.codejsha.bookstore.payment.infrastructure.support.utils

import org.jooq.JSON
import tools.jackson.core.type.TypeReference
import tools.jackson.databind.ObjectMapper

private val objectMapper = ObjectMapper()

private val mapTypeRef = object : TypeReference<Map<String, Any>>() {}

fun parseJson(json: JSON): Map<String, Any>? {
    return try {
        objectMapper.readValue(json.data(), mapTypeRef)
    } catch (_: Exception) {
        null
    }
}

fun toJson(value: Map<String, Any>?): JSON? {
    return value?.let { JSON.json(objectMapper.writeValueAsString(it)) }
}
