import { Button, Container, Form } from 'react-bootstrap'
import { FormEvent, useState } from 'react'
import { Device } from '@src/api'
import { useNavigate } from 'react-router'
import { useApi } from '@src/ApiContext.ts'
import TypeToFindDeviceType from '@src/components/deviceType/TypeToFindDeviceType.tsx'

function DeviceAdd() {
  const api = useApi()
  const [device, setDevice] = useState<Partial<Device>>({})
  const [isAdding, setIsAdding] = useState(false)
  const navigate = useNavigate()
  const basePath = import.meta.env.BASE_URL

  function addDevice(e: FormEvent<HTMLElement>) {
    e.preventDefault()
    if (!isAdding) {
      setIsAdding(true)
      api.deviceApi
        .createDevice({ device: { ...device, version: 0 } })
        .then((rt) => navigate(`${basePath}/device/${rt.id}`))
        .catch((err: unknown) => {
          console.log(err)
        })
        .finally(() => {
          setIsAdding(false)
        })
    }
  }

  function setDeviceProperty(key: keyof Device, value: string | object | undefined) {
    if (value === undefined) {
      const { [key]: _, ...d } = device
      setDevice(d as Device)
    } else {
      setDevice({ ...device, [key]: value })
    }
  }

  return (
    <Container>
      <h1>Add Device</h1>
      <Form onSubmit={addDevice}>
        <fieldset disabled={isAdding}>
          <Form.Group className="mb-3" controlId="deviceName">
            <Form.Label>Name</Form.Label>
            <Form.Control
              type="text"
              placeholder="Enter name"
              value={device.name ?? ''}
              onChange={(e) => {
                setDeviceProperty('name', e.target.value || undefined)
              }}
            />
          </Form.Group>
          <Form.Group className="mb-3" controlId="deviceDescription">
            <Form.Label>Description</Form.Label>
            <Form.Control
              type="text"
              placeholder="Enter description"
              value={device.description ?? ''}
              onChange={(e) => {
                setDeviceProperty('description', e.target.value || undefined)
              }}
            />
          </Form.Group>
          <Form.Group className="mb-3">
            <Form.Label htmlFor="deviceType">Type</Form.Label>
            <TypeToFindDeviceType
              id={'deviceType'}
              onSelect={(selected) => {
                setDeviceProperty('deviceType', selected.length > 0 ? { id: selected[0] } : undefined)
              }}
            />
          </Form.Group>
          <Button variant="primary" type="submit">
            Add
          </Button>
        </fieldset>
      </Form>
      <details>
        <summary>Debug</summary>
        <pre>{JSON.stringify(device, null, 2)}</pre>
      </details>
    </Container>
  )
}

export default DeviceAdd
