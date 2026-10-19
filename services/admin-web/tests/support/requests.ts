export const ME_URL = /\/api\/v1\/admin\/me$/;

export function problem(status: number, title: string, detail?: string) {
  return {
    status,
    contentType: "application/json",
    body: JSON.stringify({ type: "about:blank", title, status, detail }),
  };
}

export function serverError() {
  return problem(500, "Internal Server Error");
}
