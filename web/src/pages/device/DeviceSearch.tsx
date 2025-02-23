import { useCallback, useEffect, useRef, useState } from 'react'
import { Accordion, Breadcrumb, Spinner } from 'react-bootstrap'
import { useApi } from '@src/ApiContext.tsx'
import { Device, ListDevicesRequest } from '@src/api'
import { DataTable, DataTableColumnProps, DataTableRow } from '@src/components/DataTable.tsx'
import { useSearchParams } from 'react-router'
import { URLSearchParamsEqual } from '@src/utils/urlsearchparams.ts'
import { DeviceSearchFilters } from '@src/components/device/DeviceSearchFilters.tsx'
import { isEqual } from 'lodash'

function searchOptsToParams(opts: ListDevicesRequest): URLSearchParams {
  const params = new URLSearchParams()
  if (opts.deviceType) {
    opts.deviceType.forEach((v) => params.append('deviceTypeId', v))
  }
  return params
}

function DeviceSearch() {
  const api = useApi()
  const [params, setParams] = useSearchParams()
  const loading = useRef(false)
  const [devices, setDevices] = useState<Device[] | null>(null)
  const [selected, setSelected] = useState<string[]>([])
  const [searchOpts, setSearchOpts] = useState<ListDevicesRequest>()

  const performSearch = useCallback(
    (opts: ListDevicesRequest) => {
      setParams((current) => {
        const newParams = searchOptsToParams(opts)
        if (URLSearchParamsEqual(newParams, current)) {
          loading.current = true
          api.deviceApi.listDevices(opts).then((response) => {
            setDevices(response.items)
            loading.current = false
          })
          return current
        }
        return newParams
      })
    },
    [api.deviceApi, setParams],
  )

  useEffect(() => {
    setSearchOpts((current) => {
      const newOpts: ListDevicesRequest = {}
      if (params.has('name')) {
        newOpts.nameRegex = params.get('name') as string
      }
      if (params.has('deviceTypeId')) {
        newOpts.deviceType = params.getAll('deviceTypeId')
      }
      if (isEqual(newOpts, current)) {
        return current
      }
      performSearch(newOpts)
      return newOpts
    })
  }, [params, performSearch])

  if (loading.current) {
    return (
      <div>
        <Spinner />
      </div>
    )
  }
  const basePath = import.meta.env.BASE_URL
  const columns: DataTableColumnProps[] = [
    {
      id: 'name',
      label: 'Name',
      sortable: true,
      render: (device: Device) => <a href={`${basePath}/device/${device.id}`}>{device.name}</a>,
    },
    { id: 'status', label: 'Status', sortable: true },
    {
      id: 'deviceType',
      label: 'Type',
      sortable: true,
      getDisplayValue: (device: Device) => device.deviceType?.displayName ?? '',
    },
  ]

  return (
    <div>
      <Breadcrumb>
        <Breadcrumb.Item href={`${basePath}/`}>Home</Breadcrumb.Item>
        <Breadcrumb.Item active>Search</Breadcrumb.Item>
      </Breadcrumb>
      <h2>Search</h2>
      <Accordion flush={true}>
        <Accordion.Item eventKey="0">
          <Accordion.Header as={'div'}>Filters</Accordion.Header>
          <Accordion.Body>
            <DeviceSearchFilters onSearch={performSearch} searchOpts={searchOpts} />
          </Accordion.Body>
        </Accordion.Item>
      </Accordion>
      <DataTable
        columns={columns}
        data={(devices as DataTableRow[]) || []}
        selectable={true}
        onSelect={(selected) => {
          setSelected(selected)
          setSearchOpts({ ...searchOpts, deviceType: selected })
        }}
        selected={selected}
      />
      {import.meta.env.DEV && (
        <details>
          <summary>Debug</summary>
          <pre>{JSON.stringify(devices, null, 2)}</pre>
        </details>
      )}
    </div>
  )
}

export default DeviceSearch
