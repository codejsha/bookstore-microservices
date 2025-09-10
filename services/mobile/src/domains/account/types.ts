export interface Customer {
  uid: string;
  email?: string;
  first_name?: string;
  last_name?: string;
  phone?: string;
  roles?: string[];
}

export interface CustomerUpdateRequest {
  email?: string;
  first_name?: string;
  last_name?: string;
  phone?: string;
}

export interface PointBalance {
  user_uid: string;
  balance: number;
}

export type PointChangeType = "EARN" | "SPEND" | "ADJUST" | "EXPIRE";

export interface PointHistoryItem {
  uid: string;
  amount: number;
  change_type: PointChangeType;
  reason?: string;
  created_at: string;
}

export interface PointHistoryResp {
  total: number;
  items: PointHistoryItem[];
}
