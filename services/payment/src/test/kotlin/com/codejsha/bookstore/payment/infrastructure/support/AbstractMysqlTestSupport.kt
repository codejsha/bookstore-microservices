package com.codejsha.bookstore.payment.infrastructure.support

import org.springframework.context.ApplicationContextInitializer
import org.springframework.context.ConfigurableApplicationContext
import org.springframework.test.context.ContextConfiguration
import org.testcontainers.containers.MySQLContainer

@ContextConfiguration(initializers = [AbstractMysqlTestSupport.Initializer::class])
abstract class AbstractMysqlTestSupport {
    companion object {
        val container = MySQLContainer<Nothing>("mysql:9.3.0").apply {
            portBindings = listOf("3306:3306")
            withDatabaseName("testdb")
            withUsername("root")
            withPassword("test")
        }
    }

    internal class Initializer : ApplicationContextInitializer<ConfigurableApplicationContext> {
        override fun initialize(configurableApplicationContext: ConfigurableApplicationContext) {
            container.start()
        }
    }
}
