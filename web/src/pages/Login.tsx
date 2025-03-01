import { Button, Container, Spinner } from 'react-bootstrap'
import { useEffect, useState } from 'react'
import { UserProvider } from '@src/api'
import { useApi } from '@src/ApiContext.ts'
import ErrorPage from '@src/ErrorPage.tsx'

export function Login() {
  const api = useApi()
  const [providers, setProviders] = useState<UserProvider[] | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<unknown>(null)

  useEffect(() => {
    api.userApi
      .getUserProviders()
      .then((providers) => {
        setProviders(providers)
        if (providers.length === 1) {
          window.location.href = providers[0].loginURL
        } else if (providers.length === 0) {
          setError(new Error('No login providers available'))
        }
        setLoading(false)
      })
      .catch((err: unknown) => {
        setError(err)
      })
  }, [api.userApi])

  if (error) {
    return <ErrorPage error={error} />
  }

  const providerList = () => {
    if (loading) {
      return <Spinner />
    }
    return providers?.map((provider) => (
      <Button variant="primary" size="lg" key={provider.provider} as={'a'} href={provider.loginURL}>
        {provider.loginButtonImageURL ? (
          <img src={provider.loginButtonImageURL} alt={provider.provider} />
        ) : (
          'Login with ' + provider.displayName
        )}
      </Button>
    ))
  }

  return (
    <Container>
      <h1>Login</h1>
      <div className="d-grid gap-2">{providerList()}</div>
    </Container>
  )
}

export default Login
