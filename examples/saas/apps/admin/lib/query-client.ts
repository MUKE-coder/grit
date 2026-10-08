import { QueryClient, type DefaultOptions } from "@tanstack/react-query";

export const adminQueryDefaults: DefaultOptions = {
  queries: {
    staleTime: 5 * 60 * 1000,
    retry: 1,
    refetchOnWindowFocus: false,
  },
  mutations: {
    retry: 0,
  },
};

// A new client, for <Providers> to create once per mount.
//
// Never a client at module scope. A module is loaded once per server process, so
// a client created there is one cache for every request the server renders: the
// moment a page prefetches or suspends on a query, one user's data is in the
// next user's HTML.
export function makeQueryClient(): QueryClient {
  return new QueryClient({ defaultOptions: adminQueryDefaults });
}
