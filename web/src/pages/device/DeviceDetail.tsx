import { useParams } from 'react-router'
import { useEffect, useRef, useState } from 'react'
import { Device } from '@src/api'
import { Breadcrumb } from 'react-bootstrap'
import ErrorPage from '../../ErrorPage.tsx'
import { useApi } from '@src/ApiContext.ts'

function DeviceDetail() {
  const { uuid } = useParams<{ uuid: string }>()
  const [device, setDevice] = useState<Device | null>(null)
  const loading = useRef(false)
  const [error, setError] = useState<unknown>(null)
  const api = useApi()
  const basePath = import.meta.env.BASE_URL

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

    fetchDevice().catch((err: unknown) => {
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
    <div>
      <Breadcrumb>
        <Breadcrumb.Item href={basePath}>Home</Breadcrumb.Item>
        <Breadcrumb.Item href={`${basePath}/device/search`}>Device</Breadcrumb.Item>
        <Breadcrumb.Item active>Detail</Breadcrumb.Item>
      </Breadcrumb>
      <h1>Device Detail</h1>
      {device && (
        <div>
          <p>Name: {device.name}</p>
          <p>Description: {device.description}</p>
          <p>Type: {device.deviceType?.displayName}</p>
        </div>
      )}
    </div>
  )
}

export default DeviceDetail
