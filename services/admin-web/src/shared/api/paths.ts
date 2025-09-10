const V1 = "/api/v1";

export const paths = {
  admin: {
    me: `${V1}/admin/me`,
    dashboard: `${V1}/admin/dashboard`,
  },
  risk: {
    list: `${V1}/admin/risk`,
    entry: (uid: string) => `${V1}/admin/risk/${uid}`,
  },
  catalog: {
    works: `${V1}/admin/catalog/works`,
    work: (uid: string) => `${V1}/admin/catalog/works/${uid}`,
    authors: `${V1}/admin/catalog/authors`,
    subjects: `${V1}/admin/catalog/subjects`,
  },
} as const;
