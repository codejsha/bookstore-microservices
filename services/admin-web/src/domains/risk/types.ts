export type RiskLevel = "restrict" | "block";

export interface RiskEntry {
  user_uid: string;
  level: RiskLevel;
  reason: string;
  flagged_by?: string;
  flagged_at: string;
  expires_at: string;
}

export interface RiskFlagRequest {
  level: RiskLevel;
  reason: string;
  ttl_seconds?: number;
}

export interface RiskEntryList {
  entries: RiskEntry[];
}
