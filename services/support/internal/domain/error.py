class DomainError(Exception):
    pass


class UnknownReferenceError(DomainError):
    pass


class ConflictError(DomainError):
    pass
