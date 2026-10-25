from uuid import UUID


def uuid_to_bytes(uid: UUID) -> bytes:
    return uid.bytes


def bytes_to_uuid(data: bytes) -> UUID:
    return UUID(bytes=data)
