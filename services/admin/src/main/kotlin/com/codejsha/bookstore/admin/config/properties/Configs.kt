package com.codejsha.bookstore.admin.config.properties

import org.springframework.boot.context.properties.ConfigurationProperties

@ConfigurationProperties(prefix = "downstream")
data class RestClientProperties(
    val catalogBaseUrl: String,
    val customerBaseUrl: String,
    val identityBaseUrl: String,
    val fulfillmentBaseUrl: String,
    val orderBaseUrl: String,
    val paymentBaseUrl: String,
    val settlementBaseUrl: String,
)
