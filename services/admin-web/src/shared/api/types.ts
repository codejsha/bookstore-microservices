export interface Paged<T> {
  total: number;
  items: T[];
}

export interface ListParams {
  page: number;
  size: number;
  sort?: string;
}

type Swap<V, R> = undefined extends V ? R | undefined : R;

export type Wire<
  T,
  DateKeys extends keyof T = never,
  IntKeys extends keyof T = never,
> = {
  [P in keyof T]: P extends DateKeys
    ? Swap<T[P], string>
    : P extends IntKeys
      ? Swap<T[P], number>
      : T[P];
};
