package com.codejsha.bookstore.admin.config

import org.springframework.boot.context.properties.ConfigurationPropertiesScan
import org.springframework.context.annotation.Configuration

@Configuration
@ConfigurationPropertiesScan(basePackages = ["com.codejsha.bookstore.admin.config.properties"])
class PropertyScan
