package com.codejsha.bookstore.settlement

import org.springframework.boot.Banner
import org.springframework.boot.autoconfigure.SpringBootApplication
import org.springframework.boot.runApplication

@SpringBootApplication
class SettlementApplication

fun main(args: Array<String>) {
    runApplication<SettlementApplication>(*args) {
        setBannerMode(Banner.Mode.OFF)
    }
}
