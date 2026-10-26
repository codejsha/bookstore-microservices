from enum import StrEnum


class StockChangeType(StrEnum):
    INBOUND = "INBOUND"
    OUTBOUND = "OUTBOUND"
    ADJUSTMENT = "ADJUSTMENT"
    RESERVATION = "RESERVATION"
    RELEASE = "RELEASE"
