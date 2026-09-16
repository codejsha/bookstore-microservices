package com.codejsha.bookstore.order.infrastructure.support

import org.springframework.http.MediaType
import org.springframework.util.LinkedMultiValueMap
import org.springframework.web.client.RestClient
import tools.jackson.core.type.TypeReference
import tools.jackson.databind.ObjectMapper
import java.time.Duration

class KeycloakTemporalTokenFetcher(
    private val restClient: RestClient,
    private val tokenUrl: String,
    private val clientId: String,
    private val clientSecret: String,
) : TemporalTokenFetcher {

    override fun fetch(): TemporalAccessToken {
        val form = LinkedMultiValueMap<String, String>().apply {
            add("grant_type", "client_credentials")
            add("client_id", clientId)
            add("client_secret", clientSecret)
        }
        val body = restClient.post()
            .uri(tokenUrl)
            .contentType(MediaType.APPLICATION_FORM_URLENCODED)
            .accept(MediaType.APPLICATION_JSON)
            .body(form)
            .retrieve()
            .body(String::class.java)
        val payload = objectMapper.readValue(checkNotNull(body) { "token endpoint returned an empty body" }, mapTypeRef)
        val accessToken = payload["access_token"] as? String
        val expiresIn = (payload["expires_in"] as? Number)?.toLong()
        check(!accessToken.isNullOrBlank()) { "token endpoint response has no access_token" }
        check(expiresIn != null && expiresIn > 0) { "token endpoint response has no positive expires_in" }
        return TemporalAccessToken(accessToken, Duration.ofSeconds(expiresIn))
    }

    private companion object {
        val objectMapper = ObjectMapper()
        val mapTypeRef = object : TypeReference<Map<String, Any>>() {}
    }
}
