export interface Author {
  uid: string;
  name: string;
  bio?: string;
}

export interface Publisher {
  uid: string;
  name: string;
}

export interface Category {
  uid: string;
  name: string;
  sortOrder?: number;
}

export type BookStatus = "ACTIVE" | "DRAFT" | "OUT_OF_PRINT";

export interface Book {
  uid: string;
  title: string;
  sku?: string;
  isbn10?: string;
  isbn13?: string;
  price: number;
  description?: string;
  category: Category;
  publisher: Publisher;
  authors: Author[];
  status: BookStatus;
  publishedAt?: string;
  createdAt?: string;
  updatedAt?: string;
}

export interface BookFindAllResp {
  total: number;
  items: Book[];
}

export interface CategoryFindAllResp {
  total: number;
  items: Category[];
}

export interface BookFindAllParams {
  title?: string;
  isbn?: string;
  category_uid?: string;
  publisher_uid?: string;
  author_uid?: string;
  status?: string;
  size?: number;
  page?: number;
  sort?: string;
}
