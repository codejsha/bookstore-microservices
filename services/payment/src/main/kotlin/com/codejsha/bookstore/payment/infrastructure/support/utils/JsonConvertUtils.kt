package com.codejsha.bookstore.payment.infrastructure.support.utils

import org.jooq.JSONB
import tools.jackson.core.type.TypeReference
import tools.jackson.databind.ObjectMapper

private val objectMapper = ObjectMapper()

private val mapTypeRef = object : TypeReference<Map<String, Any>>() {}

fun parseJson(json: JSONB): Map<String, Any>? {
    return try {
        objectMapper.readValue(json.data(), mapTypeRef)
    } catch (_: Exception) {
        null
    }
}

fun toJson(value: Map<String, Any>?): JSONB? {
    return value?.let { JSONB.jsonb(objectMapper.writeValueAsString(it)) }
}
