import { useCallback, useEffect, useRef, useState } from 'react'
import { Accordion, Breadcrumb } from 'react-bootstrap'
import { useApi } from '@src/ApiContext.tsx'
import { Device, DeviceSearchBody, ListDevicesRequest } from '@src/api'
import { DataTable, DataTableColumnProps, DataTableRow } from '@src/components/DataTable.tsx'
import { useSearchParams } from 'react-router'
import { URLSearchParamsEqual } from '@src/utils/urlsearchparams.ts'
import { DeviceSearchFilters } from '@src/components/device/DeviceSearchFilters.tsx'

function searchOptsToParams(opts: ListDevicesRequest): URLSearchParams {
  const params = new URLSearchParams()
  if (opts.deviceSearchBody?.deviceType) {
    opts.deviceSearchBody.deviceType.forEach((v) => params.append('deviceTypeId', v))
  }
  if (opts.deviceSearchBody?.nameRegex) {
    params.set('nameRegex', opts.deviceSearchBody.nameRegex)
  }
  return params
}

function paramsToSearchOpts(params: URLSearchParams): ListDevicesRequest {
  const opts: ListDevicesRequest = {}
  const searchOpts: DeviceSearchBody = {}
  if (params.has('nameRegex')) {
    searchOpts.nameRegex = params.get('nameRegex') as string
  }
  if (params.has('deviceTypeId')) {
    searchOpts.deviceType = params.getAll('deviceTypeId')
  }
  opts.deviceSearchBody = searchOpts
  return opts
}

function DeviceSearch() {
  const api = useApi()
  const [params, setParams] = useSearchParams()
  const loading = useRef(false)
  const [devices, setDevices] = useState<Device[] | null>(null)
  const [selected, setSelected] = useState<string[]>([])

  const performSearch = useCallback(
    (opts: ListDevicesRequest) => {
      const newParams = searchOptsToParams(opts)
      if (!URLSearchParamsEqual(newParams, params)) {
        setParams(newParams)
        return
      }
      loading.current = true
      api.deviceApi.listDevices(opts).then((response) => {
        setDevices(response.items)
        loading.current = false
      })
    },
    [api.deviceApi, params, setParams],
  )

  useEffect(() => {
    performSearch(paramsToSearchOpts(params))
  }, [params, performSearch])

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
            <DeviceSearchFilters
              onSearch={(opts) => performSearch({ deviceSearchBody: opts })}
              searchOpts={paramsToSearchOpts(params)?.deviceSearchBody}
            />
          </Accordion.Body>
        </Accordion.Item>
      </Accordion>
      {((devices?.length ?? 0) > 0 && (
        <DataTable
          columns={columns}
          data={(devices as DataTableRow[]) || []}
          selectable={true}
          onSelect={(selected) => setSelected(selected)}
          selected={selected}
          loading={loading.current}
        />
      )) || <div>No devices found</div>}
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
