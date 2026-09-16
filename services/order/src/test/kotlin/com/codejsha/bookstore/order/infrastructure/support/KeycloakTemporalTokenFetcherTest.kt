package com.codejsha.bookstore.order.infrastructure.support

import org.junit.jupiter.api.Test
import org.springframework.http.HttpMethod
import org.springframework.http.MediaType
import org.springframework.test.web.client.MockRestServiceServer
import org.springframework.test.web.client.match.MockRestRequestMatchers.content
import org.springframework.test.web.client.match.MockRestRequestMatchers.method
import org.springframework.test.web.client.match.MockRestRequestMatchers.requestTo
import org.springframework.test.web.client.response.MockRestResponseCreators.withServerError
import org.springframework.test.web.client.response.MockRestResponseCreators.withSuccess
import org.springframework.util.LinkedMultiValueMap
import org.springframework.web.client.RestClient
import org.springframework.web.client.RestClientException
import java.time.Duration
import kotlin.test.assertEquals
import kotlin.test.assertFailsWith

class KeycloakTemporalTokenFetcherTest {

    private val builder = RestClient.builder()
    private val server = MockRestServiceServer.bindTo(builder).build()
    private val fetcher = KeycloakTemporalTokenFetcher(builder.build(), TOKEN_URL, "temporal-worker-order", "s3cret")

    @Test
    fun fetch_clientCredentialsGrant_parsesAccessTokenAndLifetime() {
        server.expect(requestTo(TOKEN_URL))
            .andExpect(method(HttpMethod.POST))
            .andExpect(
                content().formData(
                    LinkedMultiValueMap<String, String>().apply {
                        add("grant_type", "client_credentials")
                        add("client_id", "temporal-worker-order")
                        add("client_secret", "s3cret")
                    },
                ),
            )
            .andRespond(
                withSuccess(
                    """{"access_token":"abc.def.ghi","expires_in":300,"token_type":"Bearer"}""",
                    MediaType.APPLICATION_JSON,
                ),
            )

        val token = fetcher.fetch()

        assertEquals(TemporalAccessToken("abc.def.ghi", Duration.ofSeconds(300)), token)
        server.verify()
    }

    @Test
    fun fetch_responseWithoutAccessToken_throwsIllegalState() {
        server.expect(requestTo(TOKEN_URL))
            .andRespond(withSuccess("""{"expires_in":300}""", MediaType.APPLICATION_JSON))

        assertFailsWith<IllegalStateException> { fetcher.fetch() }
    }

    @Test
    fun fetch_serverError_throwsRestClientException() {
        server.expect(requestTo(TOKEN_URL)).andRespond(withServerError())

        assertFailsWith<RestClientException> { fetcher.fetch() }
    }

    private companion object {
        const val TOKEN_URL = "http://keycloak.test/realms/bookstore/protocol/openid-connect/token"
    }
}
