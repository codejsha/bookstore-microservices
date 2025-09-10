/* eslint-disable */

// @ts-nocheck

import { Route as rootRouteImport } from "./routes/__root"
import { Route as UnauthorizedRouteImport } from "./routes/unauthorized"
import { Route as AuthedRouteImport } from "./routes/_authed"
import { Route as AuthedIndexRouteImport } from "./routes/_authed/index"
import { Route as AuthedRiskIndexRouteImport } from "./routes/_authed/risk/index"
import { Route as AuthedCatalogIndexRouteImport } from "./routes/_authed/catalog/index"
import { Route as AuthedCatalogNewRouteImport } from "./routes/_authed/catalog/new"
import { Route as AuthedCatalogUidRouteImport } from "./routes/_authed/catalog/$uid"

const UnauthorizedRoute = UnauthorizedRouteImport.update({
  id: "/unauthorized",
  path: "/unauthorized",
  getParentRoute: () => rootRouteImport,
} as any)
const AuthedRoute = AuthedRouteImport.update({
  id: "/_authed",
  getParentRoute: () => rootRouteImport,
} as any)
const AuthedIndexRoute = AuthedIndexRouteImport.update({
  id: "/",
  path: "/",
  getParentRoute: () => AuthedRoute,
} as any)
const AuthedRiskIndexRoute = AuthedRiskIndexRouteImport.update({
  id: "/risk/",
  path: "/risk/",
  getParentRoute: () => AuthedRoute,
} as any)
const AuthedCatalogIndexRoute = AuthedCatalogIndexRouteImport.update({
  id: "/catalog/",
  path: "/catalog/",
  getParentRoute: () => AuthedRoute,
} as any)
const AuthedCatalogNewRoute = AuthedCatalogNewRouteImport.update({
  id: "/catalog/new",
  path: "/catalog/new",
  getParentRoute: () => AuthedRoute,
} as any)
const AuthedCatalogUidRoute = AuthedCatalogUidRouteImport.update({
  id: "/catalog/$uid",
  path: "/catalog/$uid",
  getParentRoute: () => AuthedRoute,
} as any)

export interface FileRoutesByFullPath {
  "/": typeof AuthedIndexRoute
  "/unauthorized": typeof UnauthorizedRoute
  "/catalog/$uid": typeof AuthedCatalogUidRoute
  "/catalog/new": typeof AuthedCatalogNewRoute
  "/risk/": typeof AuthedRiskIndexRoute
  "/catalog/": typeof AuthedCatalogIndexRoute
}

export interface FileRoutesByTo {
  "/unauthorized": typeof UnauthorizedRoute
  "/": typeof AuthedIndexRoute
  "/catalog/$uid": typeof AuthedCatalogUidRoute
  "/catalog/new": typeof AuthedCatalogNewRoute
  "/risk": typeof AuthedRiskIndexRoute
  "/catalog": typeof AuthedCatalogIndexRoute
}

export interface FileRoutesById {
  __root__: typeof rootRouteImport
  "/_authed": typeof AuthedRouteWithChildren
  "/unauthorized": typeof UnauthorizedRoute
  "/_authed/": typeof AuthedIndexRoute
  "/_authed/catalog/$uid": typeof AuthedCatalogUidRoute
  "/_authed/catalog/new": typeof AuthedCatalogNewRoute
  "/_authed/risk/": typeof AuthedRiskIndexRoute
  "/_authed/catalog/": typeof AuthedCatalogIndexRoute
}

export interface FileRouteTypes {
  fileRoutesByFullPath: FileRoutesByFullPath
  fullPaths:
      | "/"
      | "/unauthorized"
      | "/catalog/$uid"
      | "/catalog/new"
      | "/risk/"
      | "/catalog/"
  fileRoutesByTo: FileRoutesByTo
  to:
      | "/unauthorized"
      | "/"
      | "/catalog/$uid"
      | "/catalog/new"
      | "/risk"
      | "/catalog"
  id:
      | "__root__"
      | "/_authed"
      | "/unauthorized"
      | "/_authed/"
      | "/_authed/catalog/$uid"
      | "/_authed/catalog/new"
      | "/_authed/risk/"
      | "/_authed/catalog/"
  fileRoutesById: FileRoutesById
}

export interface RootRouteChildren {
  AuthedRoute: typeof AuthedRouteWithChildren
  UnauthorizedRoute: typeof UnauthorizedRoute
}

declare module "@tanstack/react-router" {
  interface FileRoutesByPath {
    "/unauthorized": {
      id: "/unauthorized"
      path: "/unauthorized"
      fullPath: "/unauthorized"
      preLoaderRoute: typeof UnauthorizedRouteImport
      parentRoute: typeof rootRouteImport
    }
    "/_authed": {
      id: "/_authed"
      path: ""
      fullPath: "/"
      preLoaderRoute: typeof AuthedRouteImport
      parentRoute: typeof rootRouteImport
    }
    "/_authed/": {
      id: "/_authed/"
      path: "/"
      fullPath: "/"
      preLoaderRoute: typeof AuthedIndexRouteImport
      parentRoute: typeof AuthedRoute
    }
    "/_authed/risk/": {
      id: "/_authed/risk/"
      path: "/risk"
      fullPath: "/risk/"
      preLoaderRoute: typeof AuthedRiskIndexRouteImport
      parentRoute: typeof AuthedRoute
    }
    "/_authed/catalog/": {
      id: "/_authed/catalog/"
      path: "/catalog"
      fullPath: "/catalog/"
      preLoaderRoute: typeof AuthedCatalogIndexRouteImport
      parentRoute: typeof AuthedRoute
    }
    "/_authed/catalog/new": {
      id: "/_authed/catalog/new"
      path: "/catalog/new"
      fullPath: "/catalog/new"
      preLoaderRoute: typeof AuthedCatalogNewRouteImport
      parentRoute: typeof AuthedRoute
    }
    "/_authed/catalog/$uid": {
      id: "/_authed/catalog/$uid"
      path: "/catalog/$uid"
      fullPath: "/catalog/$uid"
      preLoaderRoute: typeof AuthedCatalogUidRouteImport
      parentRoute: typeof AuthedRoute
    }
  }
}

interface AuthedRouteChildren {
  AuthedIndexRoute: typeof AuthedIndexRoute
  AuthedRiskIndexRoute: typeof AuthedRiskIndexRoute
  AuthedCatalogUidRoute: typeof AuthedCatalogUidRoute
  AuthedCatalogNewRoute: typeof AuthedCatalogNewRoute
  AuthedCatalogIndexRoute: typeof AuthedCatalogIndexRoute
}

const AuthedRouteChildren: AuthedRouteChildren = {
  AuthedIndexRoute: AuthedIndexRoute,
  AuthedRiskIndexRoute: AuthedRiskIndexRoute,
  AuthedCatalogUidRoute: AuthedCatalogUidRoute,
  AuthedCatalogNewRoute: AuthedCatalogNewRoute,
  AuthedCatalogIndexRoute: AuthedCatalogIndexRoute,
}

const AuthedRouteWithChildren =
    AuthedRoute._addFileChildren(AuthedRouteChildren)

const rootRouteChildren: RootRouteChildren = {
  AuthedRoute: AuthedRouteWithChildren,
  UnauthorizedRoute: UnauthorizedRoute,
}
export const routeTree = rootRouteImport
    ._addFileChildren(rootRouteChildren)
    ._addFileTypes<FileRouteTypes>()
