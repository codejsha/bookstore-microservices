package com.codejsha.bookstore.commonlib.domain.exception

class ResourceNotFoundException : RuntimeException {
    constructor() : super("Resource not found")
    constructor(message: String) : super(message)
    constructor(message: String, cause: Throwable) : super(message, cause)
}
