import { Button, Container, Form } from 'react-bootstrap'
import { FormEvent, useState } from 'react'
import { Device, DeviceTypeResult, DeviceTypeSearchResults } from '../api'
import { AsyncTypeahead } from 'react-bootstrap-typeahead'
import { Option } from 'react-bootstrap-typeahead/types/types'
import { useApi } from '../ApiContext.tsx'
import { useNavigate } from 'react-router'

function DeviceAdd() {
  const api = useApi()
  const [device, setDevice] = useState<Device>({
    version: 0,
  })
  const [deviceTypeOptions, setDeviceTypeOptions] = useState<DeviceTypeSearchResults>()
  const [deviceTypeLoading, setDeviceTypeLoading] = useState(false)
  const navigate = useNavigate()
  const basePath = import.meta.env.BASE_URL

  async function lookupDeviceTypes(query: string) {
    console.log('loading', query)
    setDeviceTypeLoading(true)
    try {
      const response = await api.deviceApi.listDeviceTypes(/*{ q: query }*/)
      setDeviceTypeOptions(response)
    } finally {
      setDeviceTypeLoading(false)
    }
  }

  async function addDevice(e: FormEvent<HTMLElement>) {
    e.preventDefault()
    const rt = await api.deviceApi.createDevice({ device })
    console.log('Device added', rt)
    navigate(`${basePath}/device/${rt.id}`)
  }

  function renderDeviceTypeMenuItem(option: Option) {
    return (
      <div>
        <div>{(option as DeviceTypeResult).model}</div>
      </div>
    )
  }

  function setDeviceProperty(key: keyof Device, value: any) {
    if (value === undefined) {
      const d = { ...device }
      delete d[key]
      setDevice(d)
    } else {
      setDevice({ ...device, [key]: value })
    }
  }

  return (
    <Container>
      <h1>Add Device</h1>
      <Form onSubmit={addDevice}>
        <Form.Group className="mb-3" controlId="deviceName">
          <Form.Label>Name</Form.Label>
          <Form.Control
            type="text"
            placeholder="Enter name"
            value={device.name ?? ''}
            onChange={(e) => setDeviceProperty('name', e.target.value || undefined)}
          />
        </Form.Group>
        <Form.Group className="mb-3" controlId="deviceDescription">
          <Form.Label>Description</Form.Label>
          <Form.Control
            type="text"
            placeholder="Enter description"
            value={device.description ?? ''}
            onChange={(e) => setDeviceProperty('description', e.target.value || undefined)}
          />
        </Form.Group>
        <Form.Group className="mb-3">
          <Form.Label htmlFor="deviceType">Type</Form.Label>
          <AsyncTypeahead
            filterBy={() => true}
            isLoading={deviceTypeLoading}
            onSearch={lookupDeviceTypes}
            options={deviceTypeOptions?.items || []}
            labelKey="name"
            id="deviceType"
            onChange={(selected) =>
              setDeviceProperty(
                'deviceType',
                selected.length > 0
                  ? {
                      id: (selected[0] as DeviceTypeResult).id,
                      displayName: (selected[0] as DeviceTypeResult).model,
                    }
                  : undefined,
              )
            }
            renderMenuItemChildren={renderDeviceTypeMenuItem}
          />
        </Form.Group>
        <Button variant="primary" type="submit">
          Add
        </Button>
      </Form>
      <pre>{JSON.stringify(device, null, 2)}</pre>
    </Container>
  )
}

export default DeviceAdd
