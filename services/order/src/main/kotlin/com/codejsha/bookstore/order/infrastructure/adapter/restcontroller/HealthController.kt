package com.codejsha.bookstore.order.infrastructure.adapter.restcontroller

import org.springframework.http.ResponseEntity
import org.springframework.web.bind.annotation.GetMapping
import org.springframework.web.bind.annotation.RestController

@RestController
class HealthController {
    @GetMapping("/health")
    fun health(): ResponseEntity<Void> = ResponseEntity.ok().build()
}
