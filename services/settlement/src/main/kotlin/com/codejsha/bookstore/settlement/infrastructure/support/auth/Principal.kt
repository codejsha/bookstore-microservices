package com.codejsha.bookstore.settlement.infrastructure.support.auth

import tools.jackson.databind.ObjectMapper
import java.util.Base64

data class Principal(
    val sub: String,
    val email: String?,
    val name: String?,
    val roles: List<String>,
    val scopes: List<String>,
    val raw: Map<String, Any?>,
) {
    fun hasRole(role: String): Boolean = role in roles

    fun hasScope(scope: String): Boolean = scope in scopes
}

object PrincipalHeaders {
    const val USER_ID = "X-User-Id"
    const val USER_EMAIL = "X-User-Email"
    const val USER_NAME = "X-User-Name"
    const val USER_ROLES = "X-User-Roles"
    const val USER_SCOPES = "X-User-Scopes"
    const val JWT_PAYLOAD = "x-jwt-payload"
}

object PrincipalParser {
    fun parse(header: (String) -> String?, objectMapper: ObjectMapper): Principal? {
        val sub = header(PrincipalHeaders.USER_ID) ?: return null
        val roles = parseClaimList(header(PrincipalHeaders.USER_ROLES), ',', objectMapper)
        val scopes = parseClaimList(header(PrincipalHeaders.USER_SCOPES), ' ', objectMapper)
        val raw: Map<String, Any?> =
            header(PrincipalHeaders.JWT_PAYLOAD)?.let {
                val padded = it.padEnd((it.length + 3) / 4 * 4, '=')
                val decoded = Base64.getUrlDecoder().decode(padded)
                @Suppress("UNCHECKED_CAST")
                objectMapper.readValue(decoded, Map::class.java) as Map<String, Any?>
            } ?: emptyMap()
        return Principal(
            sub = sub,
            email = header(PrincipalHeaders.USER_EMAIL),
            name = header(PrincipalHeaders.USER_NAME),
            roles = roles,
            scopes = scopes,
            raw = raw,
        )
    }

    private fun parseClaimList(
        value: String?,
        delimiter: Char,
        objectMapper: ObjectMapper,
    ): List<String> {
        if (value.isNullOrBlank()) return emptyList()
        decodeJsonArray(value, objectMapper)?.let { return it }
        return value.split(delimiter).map { it.trim() }.filter { it.isNotEmpty() }
    }

    private fun decodeJsonArray(value: String, objectMapper: ObjectMapper): List<String>? =
        runCatching {
            val padded = value.padEnd((value.length + 3) / 4 * 4, '=')
            val decoded = Base64.getUrlDecoder().decode(padded)
            @Suppress("UNCHECKED_CAST")
            val list = objectMapper.readValue(decoded, List::class.java) as List<Any?>
            list.mapNotNull { it?.toString()?.trim()?.takeIf(String::isNotEmpty) }
        }.getOrNull()
}
