import { isEqual } from 'lodash'

export function URLSearchParamsEqual(a?: URLSearchParams, b?: URLSearchParams): boolean {
  if (a === b) return true
  if (a == null || b == null) {
    return false
  }
  const aa = Array.from(a.entries())
  const ba = Array.from(b.entries())
  return isEqual(aa, ba)
}
