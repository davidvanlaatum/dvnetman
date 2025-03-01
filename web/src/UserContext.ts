import { createContext, useContext } from 'react'
import { CurrentUser } from '@src/api'

export const UserContext = createContext<CurrentUser | null>(null)

export const useUser = () => {
  const context = useContext(UserContext)
  if (!context) {
    throw new Error('useApi must be used within a UserProvider')
  }
  return context
}
