from enum import StrEnum


class FreightStatus(StrEnum):
    ESTIMATED = "ESTIMATED"
    CONFIRMED = "CONFIRMED"
    INVOICED = "INVOICED"
    PAID = "PAID"
