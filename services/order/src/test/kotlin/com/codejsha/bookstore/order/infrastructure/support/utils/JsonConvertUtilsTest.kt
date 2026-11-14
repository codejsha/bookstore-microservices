package com.codejsha.bookstore.order.infrastructure.support.utils

import org.jooq.JSONB
import org.junit.jupiter.api.Test
import kotlin.test.assertEquals
import kotlin.test.assertNotNull
import kotlin.test.assertNull

class JsonConvertUtilsTest {

    @Test
    fun `toJson_nullMap_yieldsNull`() {
        assertNull(toJson(null))
    }

    @Test
    fun `toJson_emptyMap_yieldsEmptyJsonObject`() {
        assertEquals("{}", toJson(emptyMap())?.data())
    }

    @Test
    fun `toJson_nestedMap_roundTripsThroughParseJson`() {
        val source = mapOf(
            "code" to "WELCOME",
            "percent" to 10,
            "stackable" to false,
            "scope" to mapOf("channel" to "web"),
            "tags" to listOf("new", "vip"),
        )

        val parsed = parseJson(assertNotNull(toJson(source)))

        assertEquals(source, parsed)
    }

    @Test
    fun `parseJson_jsonObject_yieldsMap`() {
        val parsed = parseJson(JSONB.jsonb("""{"carrier": "CJ", "weight": 1.5}"""))

        assertEquals(mapOf("carrier" to "CJ", "weight" to 1.5), parsed)
    }

    @Test
    fun `parseJson_malformedJson_yieldsNull`() {
        assertNull(parseJson(JSONB.jsonb("{not json")))
    }

    @Test
    fun `parseJson_jsonArray_yieldsNull`() {
        assertNull(parseJson(JSONB.jsonb("[1, 2, 3]")))
    }
}
