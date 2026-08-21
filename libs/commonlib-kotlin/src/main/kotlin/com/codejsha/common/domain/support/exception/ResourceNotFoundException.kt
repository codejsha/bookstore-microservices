package com.codejsha.common.domain.support.exception

class ResourceNotFoundException : RuntimeException {
    constructor() : super("Resource not found")
    constructor(message: String) : super(message)
    constructor(message: String, cause: Throwable) : super(message, cause)
}
