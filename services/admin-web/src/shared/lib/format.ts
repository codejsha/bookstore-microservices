export const EMPTY = "—";

const dateTimeFormat = new Intl.DateTimeFormat(undefined, {
  dateStyle: "medium",
  timeStyle: "short",
});

const dateFormat = new Intl.DateTimeFormat(undefined, { dateStyle: "medium" });

export function formatDateTime(value: string | null | undefined): string {
  if (!value) return EMPTY;
  const date = new Date(value);
  return Number.isNaN(date.getTime()) ? value : dateTimeFormat.format(date);
}

export function formatDate(value: string | null | undefined): string {
  if (!value) return EMPTY;
  const [year, month, day] = value.split("-").map(Number);
  if (!year || !month || !day) return value;
  return dateFormat.format(new Date(year, month - 1, day));
}

function currencyFormat(currency: string): Intl.NumberFormat | null {
  try {
    return new Intl.NumberFormat(undefined, {
      style: "currency",
      currency,
      currencyDisplay: "code",
    });
  } catch {
    return null;
  }
}

export function formatMoney(
  amount: number | null | undefined,
  currency: string,
): string {
  if (amount === null || amount === undefined) return EMPTY;
  const format = currencyFormat(currency);
  return format ? format.format(amount) : `${amount} ${currency}`;
}

export function formatMinorUnits(
  amount: number | null | undefined,
  currency: string,
): string {
  if (amount === null || amount === undefined) return EMPTY;
  const format = currencyFormat(currency);
  if (!format) return `${amount} ${currency} (minor units)`;
  const digits = format.resolvedOptions().maximumFractionDigits ?? 0;
  return format.format(amount / 10 ** digits);
}

const percentFormat = new Intl.NumberFormat(undefined, {
  style: "percent",
  maximumFractionDigits: 2,
});

export function formatRate(rate: number | null | undefined): string {
  return rate === null || rate === undefined
    ? EMPTY
    : percentFormat.format(rate);
}

export function orEmpty(value: string | number | null | undefined): string {
  return value === null || value === undefined || value === ""
    ? EMPTY
    : String(value);
}
