import { createContext, FC, lazy, ReactNode, useContext, useEffect, useState } from 'react'
import { CurrentUser } from '@src/api'
import { useApi } from '@src/ApiContext.tsx'
import { Spinner } from 'react-bootstrap'

const Login = lazy(() => import('./pages/Login.tsx'))

const UserContext = createContext<CurrentUser | null>(null)

export const UserProvider: FC<{ children: ReactNode; forceUser?: CurrentUser }> = ({ children, forceUser }) => {
  const api = useApi()
  const [user, setUser] = useState<CurrentUser | null>(null)
  const [loading, setLoading] = useState(true)
  useEffect(() => {
    if (forceUser) {
      setUser(forceUser)
      setLoading(false)
      return
    }
    const fetchUser = async () => {
      try {
        const user = await api.userApi.getCurrentUser()
        setUser(user)
      } finally {
        setLoading(false)
      }
    }
    fetchUser().then()
    const id = setInterval(fetchUser, 1000 * 60)
    return () => clearInterval(id)
  }, [api.userApi, forceUser])

  if (loading) {
    return (
      <div>
        <Spinner />
        Loading...
      </div>
    )
  } else if (!user?.loggedIn) {
    return <Login />
  }
  return <UserContext.Provider value={user}>{children}</UserContext.Provider>
}

// eslint-disable-next-line react-refresh/only-export-components
export const useUser = () => {
  const context = useContext(UserContext)
  if (!context) {
    throw new Error('useApi must be used within a UserContext')
  }
  return context
}
