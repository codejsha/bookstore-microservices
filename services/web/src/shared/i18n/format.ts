import type { Locale } from "./locales";

type DateInput = string | number | Date;

const DATE_TIME_OPTIONS: Intl.DateTimeFormatOptions = {
  year: "numeric",
  month: "numeric",
  day: "numeric",
  hour: "numeric",
  minute: "numeric",
  second: "numeric",
};

const NON_BREAKING_SPACE = / /g;

export interface Formatters {
  date: (value: DateInput) => string;
  dateTime: (value: DateInput) => string;
  number: (value: number, options?: Intl.NumberFormatOptions) => string;
  price: (amount: number, currency: string) => string;
  money: (amount: number, currency: string) => string;
}

export function createFormatters(locale: Locale): Formatters {
  const dateFormat = new Intl.DateTimeFormat(locale);
  const dateTimeFormat = new Intl.DateTimeFormat(locale, DATE_TIME_OPTIONS);
  const plainNumberFormat = new Intl.NumberFormat(locale);
  const currencyFormats = new Map<string, Intl.NumberFormat>();

  const currencyFormat = (
    currency: string,
    currencyDisplay: "symbol" | "code",
  ) => {
    const key = `${currency}:${currencyDisplay}`;
    let format = currencyFormats.get(key);
    if (!format) {
      format = new Intl.NumberFormat(locale, {
        style: "currency",
        currency,
        currencyDisplay,
      });
      currencyFormats.set(key, format);
    }
    return format;
  };

  return {
    date: (value) => dateFormat.format(new Date(value)),
    dateTime: (value) => dateTimeFormat.format(new Date(value)),
    number: (value, options) =>
      options
        ? new Intl.NumberFormat(locale, options).format(value)
        : plainNumberFormat.format(value),
    price: (amount, currency) =>
      currencyFormat(currency, "symbol").format(amount),
    money: (amount, currency) =>
      currencyFormat(currency, "code")
        .format(amount)
        .replace(NON_BREAKING_SPACE, " "),
  };
}
