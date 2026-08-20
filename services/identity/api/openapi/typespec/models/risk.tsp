import "@typespec/http";

using Http;

namespace IdentityService;

@doc("Enforcement level for a flagged principal")
enum RiskLevel {
  restrict,
  block,
}

model RiskFlagRequest {
  level: RiskLevel;
  reason: string;

  @doc("Time-to-live in seconds; defaults to 86400 (24h), capped at 30 days")
  ttl_seconds?: int64;
}

model RiskEntryResponse {
  user_uid: uuid;
  level: RiskLevel;
  reason: string;
  flagged_by?: string;
  flagged_at: utcDateTime;
  expires_at: utcDateTime;
}

model RiskEntryListResponse {
  entries: RiskEntryResponse[];
}
