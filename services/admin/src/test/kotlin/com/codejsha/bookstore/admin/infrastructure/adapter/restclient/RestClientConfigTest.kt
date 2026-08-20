package com.codejsha.bookstore.admin.infrastructure.adapter.restclient

import com.codejsha.bookstore.admin.config.properties.RestClientProperties
import com.codejsha.bookstore.admin.infrastructure.support.auth.RequestTokenResolver
import org.junit.jupiter.api.Test
import kotlin.test.assertNotNull

class RestClientConfigTest {

    private val config = RestClientConfig(
        config = RestClientProperties(
            catalogBaseUrl = "http://catalog",
            customerBaseUrl = "http://customer",
            identityBaseUrl = "http://identity",
            inventoryBaseUrl = "http://inventory",
            orderBaseUrl = "http://order",
            paymentBaseUrl = "http://payment",
            settlementBaseUrl = "http://settlement",
        ),
        tokenPropagation = BearerTokenPropagationInterceptor(RequestTokenResolver()),
    )

    @Test
    fun `builds a proxy for every generated downstream interface`() {
        assertNotNull(config.catalogWorkApi())
        assertNotNull(config.catalogAuthorApi())
        assertNotNull(config.catalogSubjectApi())
        assertNotNull(config.identityUserApi())
        assertNotNull(config.orderOrderApi())
        assertNotNull(config.inventoryWarehouseApi())
        assertNotNull(config.inventoryStockApi())
        assertNotNull(config.paymentPaymentApi())
        assertNotNull(config.paymentRefundApi())
        assertNotNull(config.settlementSettlementApi())
    }
}
