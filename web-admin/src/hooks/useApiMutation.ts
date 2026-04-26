import { useMutation, type UseMutationOptions } from '@tanstack/vue-query'

export function useApiMutation<TData = unknown, TVariables = void, TError = Error>(
  mutationFn: (variables: TVariables) => Promise<TData>,
  options?: UseMutationOptions<TData, TError, TVariables>,
) {
  return useMutation({
    mutationFn,
    ...options,
  })
}
