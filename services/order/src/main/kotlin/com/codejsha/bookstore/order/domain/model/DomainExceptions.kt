package com.codejsha.bookstore.order.domain.model

class OrderStateConflictException(message: String) : RuntimeException(message)

class CartStateConflictException(message: String) : RuntimeException(message)

class OrderOwnershipException(message: String) : RuntimeException(message)
