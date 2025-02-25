import { DeviceSearchBody } from '@src/api'
import { ChangeEvent, FC, FormEvent, Fragment, useEffect, useState } from 'react'
import { Button, ButtonGroup, Col, Form, Row } from 'react-bootstrap'
import { TypeToFindDeviceType } from '@src/components/deviceType/TypeToFindDeviceType.tsx'
import { isEqual } from 'lodash'

interface DeviceSearchFiltersProps {
  onSearch: (query: DeviceSearchBody) => void
  searchOpts?: DeviceSearchBody
}

export const DeviceSearchFilters: FC<DeviceSearchFiltersProps> = ({ onSearch, searchOpts }) => {
  const [searchOptsState, setSearchOptsState] = useState<DeviceSearchBody>()

  useEffect(() => {
    setSearchOptsState((current) => {
      if (isEqual(current, searchOpts)) {
        return current
      }
      return searchOpts
    })
  }, [searchOpts])

  function doSearch(e: FormEvent<HTMLFormElement>) {
    e.preventDefault()
    onSearch(searchOptsState ?? {})
  }

  function onNameChange(e: ChangeEvent<HTMLInputElement>) {
    setSearchOptsState((current) => {
      const newValue: DeviceSearchBody = { ...current, nameRegex: e.target.value }
      if (isEqual(newValue, current)) {
        return current
      } else if (newValue.nameRegex === '') {
        delete newValue.nameRegex
      }
      return newValue
    })
  }

  function onDeviceTypeChange(deviceType: string[]) {
    setSearchOptsState((current) => {
      const newValue: DeviceSearchBody = { ...current, deviceType }
      if (isEqual(newValue, current)) {
        return current
      } else if (newValue.deviceType?.length === 0) {
        delete newValue.deviceType
      }
      return newValue
    })
  }

  return (
    <div>
      <Form onSubmit={doSearch}>
        <Row>
          <Form.Group as={Fragment} controlId="deviceName">
            <Form.Label column={true} lg={1}>
              Name
            </Form.Label>
            <Col lg={1}>
              <Form.Control type="text" name="name" value={searchOptsState?.nameRegex ?? ''} onChange={onNameChange} />
            </Col>
          </Form.Group>
          <Form.Group as={Fragment} controlId="deviceType">
            <Form.Label column={true} lg={1}>
              Device Type
            </Form.Label>
            <Col lg={1}>
              <TypeToFindDeviceType
                id="deviceType"
                onSelect={onDeviceTypeChange}
                selected={searchOptsState?.deviceType ?? []}
                asyncProps={{ multiple: true }}
              />
            </Col>
          </Form.Group>
        </Row>
        <ButtonGroup className="float-end">
          <Button type="submit" size="sm" className="bi bi-search" />
          <Button type="button" size="sm" variant={'secondary'} onClick={() => setSearchOptsState({})}>
            Clear
          </Button>
        </ButtonGroup>
        {import.meta.env.DEV && (
          <details>
            <summary>Debug</summary>
            <pre>{JSON.stringify(searchOptsState, null, 2)}</pre>
          </details>
        )}
      </Form>
    </div>
  )
}
