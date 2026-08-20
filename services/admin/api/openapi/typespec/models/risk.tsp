import "@typespec/http";

using Http;

namespace AdminService;

@doc("Enforcement level for a flagged principal")
enum AdminRiskLevel {
  restrict,
  block,
}

model AdminRiskFlagRequest {
  level: AdminRiskLevel;
  reason: string;

  @doc("Time-to-live in seconds; defaults to 86400 (24h), capped at 30 days")
  ttl_seconds?: int64;
}

model AdminRiskEntryResponse {
  user_uid: uuid;
  level: AdminRiskLevel;
  reason: string;
  flagged_by?: string;
  flagged_at: utcDateTime;
  expires_at: utcDateTime;
}

model AdminRiskEntryListResponse {
  entries: AdminRiskEntryResponse[];
}
