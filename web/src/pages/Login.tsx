import { Button, Container, Spinner } from 'react-bootstrap'
import { useApi } from '@src/ApiContext.tsx'
import { useEffect, useState } from 'react'
import { UserProvider } from '@src/api'

export function Login() {
  const api = useApi()
  const [providers, setProviders] = useState<UserProvider[] | null>(null)

  useEffect(() => {
    const fetchProviders = async () => {
      const providers = await api.userApi.getUserProviders()
      setProviders(providers)
      if (providers.length === 1) {
        window.location.href = providers[0].loginURL
      }
    }
    fetchProviders().then()
  }, [api.userApi])

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
