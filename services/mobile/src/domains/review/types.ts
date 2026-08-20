export interface Review {
  uid: string;
  user_uid: string;
  book_uid: string;
  rating: number;
  title?: string;
  content?: string;
  created_at: string;
  updated_at?: string;
}

export interface ReviewFindAllResp {
  total: number;
  items: Review[];
}

export interface ReviewCreateReq {
  book_uid: string;
  rating: number;
  title?: string;
  content?: string;
}
