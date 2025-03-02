import { FC, lazy, ReactNode, useEffect, useState } from 'react'
import { CurrentUser } from '@src/api'
import { Spinner } from 'react-bootstrap'
import { UserContext } from '@src/UserContext.ts'
import { useApi } from '@src/ApiContext.ts'

const Login = lazy(() => import('./pages/Login.tsx'))

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
    void fetchUser().then()
    const id = setInterval(() => {
      void fetchUser()
    }, 1000 * 60)
    const cleanup = api.onUnauthorized(() => {
      void fetchUser()
    })
    return () => {
      clearInterval(id)
      cleanup()
    }
  }, [api, forceUser])

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
