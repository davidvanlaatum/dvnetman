import { Stats } from './api'
import { useEffect, useState } from 'react'
import { useApi } from '@src/ApiContext.ts'

function Index() {
  const api = useApi()
  const [stats, setStats] = useState<Stats | null>(null)
  const [initialized, setInitialized] = useState(false)

  useEffect(() => {
    const f = () => {
      api.statsApi.getStats().then((response) => {
        setStats(response)
      })
    }
    const id = setInterval(f, 30000)
    if (!initialized) {
      f()
      setInitialized(true)
    }
    return () => clearInterval(id)
  }, [api.statsApi, initialized])

  return (
    <div>
      <h1>Index</h1>
      <pre>{JSON.stringify(stats, null, 2)}</pre>
    </div>
  )
}

export default Index
