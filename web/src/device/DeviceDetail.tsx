import { useParams } from 'react-router'
import { useEffect, useRef, useState } from 'react'
import { Device } from '../api'
import { Container } from 'react-bootstrap'
import { useApi } from '../ApiContext.tsx'
import ErrorPage from '../ErrorPage.tsx'

function DeviceDetail() {
  const { uuid } = useParams<{ uuid: string }>()
  const [device, setDevice] = useState<Device | null>(null)
  const loading = useRef(false)
  const [error, setError] = useState<Error | null>(null)
  const api = useApi()

  useEffect(() => {
    async function fetchDevice() {
      if (uuid && !device && !loading.current && !error) {
        loading.current = true
        try {
          const response = await api.deviceApi.getDevice({ id: uuid })
          setDevice(response)
        } finally {
          loading.current = false
        }
      }
    }

    fetchDevice().catch((err) => {
      setError(err)
    })
  }, [api.deviceApi, device, device?.id, error, uuid])

  if (loading.current) {
    return <div>Loading...</div>
  }

  if (error) {
    return <ErrorPage error={error} />
  }

  return (
    <Container>
      <h1>Device Detail</h1>
      {device && (
        <div>
          <p>Name: {device.name}</p>
          <p>Description: {device.description}</p>
          <p>Type: {device.deviceType?.displayName}</p>
        </div>
      )}
    </Container>
  )
}

export default DeviceDetail
