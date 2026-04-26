import { useQuery, type UseQueryOptions } from '@tanstack/vue-query'

export function useApiQuery<TQueryFnData = unknown, TError = Error>(
  queryKey: unknown[],
  queryFn: () => Promise<TQueryFnData>,
  options?: Omit<UseQueryOptions<TQueryFnData, TError>, 'queryKey' | 'queryFn'>,
) {
  return useQuery({
    queryKey,
    queryFn,
    ...options,
  })
}
