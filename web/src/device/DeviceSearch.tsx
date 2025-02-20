import { useEffect, useRef, useState } from 'react'
import { Spinner } from 'react-bootstrap'
import { useApi } from '../ApiContext.tsx'
import { Device } from '../api'

function DeviceSearch() {
  const loading = useRef(false)
  const api = useApi()
  const [devices, setDevices] = useState<Device[] | null>(null)

  useEffect(() => {
    if (!loading.current) {
      loading.current = true
      api.deviceApi.listDevices({}).then((response) => {
        setDevices(response.items)
        loading.current = false
      })
    }
  }, [api.deviceApi])

  if (loading.current) {
    return (
      <div>
        <Spinner />
      </div>
    )
  }

  return (
    <div>
      <h2>Search</h2>
      {devices?.map((device) => (
        <div key={device.id}>
          <p>Name: {device.name}</p>
          <p>{device.description}</p>
        </div>
      ))}

      <pre>{JSON.stringify(devices, null, 2)}</pre>
    </div>
  )
}

export default DeviceSearch
