package com.codejsha.bookstore.payment.config

import org.springframework.boot.context.properties.ConfigurationPropertiesScan
import org.springframework.context.annotation.Configuration

@Configuration
@ConfigurationPropertiesScan(basePackages = ["com.codejsha.bookstore.payment.config.properties"])
class PropertyScan
