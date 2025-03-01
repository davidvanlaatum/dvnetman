import { Button, Container, Spinner } from 'react-bootstrap'
import { useEffect, useState } from 'react'
import { UserProvider } from '@src/api'
import { useApi } from '@src/ApiContext.ts'
import ErrorPage from '@src/ErrorPage.tsx'

export function Login() {
  const api = useApi()
  const [providers, setProviders] = useState<UserProvider[] | null>(null)
  const [error, setError] = useState<unknown>(null)

  useEffect(() => {
    const fetchProviders = async () => {
      const providers = await api.userApi.getUserProviders()
      setProviders(providers)
      if (providers.length === 1) {
        window.location.href = providers[0].loginURL
      }
    }
    fetchProviders().catch((err: unknown) => {
      setError(err)
    })
  }, [api.userApi])

  if (error) {
    return <ErrorPage error={error} />
  }

  return (
    <Container>
      <h1>Login</h1>
      <div className="d-grid gap-2">
        {(!providers || providers.length == 1) && <Spinner />}
        {(providers?.length ?? 0) > 1 &&
          providers?.map((provider) => (
            <Button variant="primary" size="lg" key={provider.provider} as={'a'} href={provider.loginURL}>
              {provider.loginButtonImageURL ? (
                <img src={provider.loginButtonImageURL} alt={provider.provider} />
              ) : (
                'Login with ' + provider.displayName
              )}
            </Button>
          ))}
      </div>
    </Container>
  )
}

export default Login
