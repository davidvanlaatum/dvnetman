import { AsyncTypeahead, Highlighter } from 'react-bootstrap-typeahead'
import { FC, useCallback, useEffect, useRef, useState } from 'react'
import { TypeaheadMenuProps } from 'react-bootstrap-typeahead/types/components/TypeaheadMenu/TypeaheadMenu'
import { DeviceTypeResult, DeviceTypeSearchResults, ListDeviceTypesRequest } from '@src/api'
import { UseAsyncProps } from 'react-bootstrap-typeahead/types/behaviors/async'
import { isEqual } from 'lodash'

import 'react-bootstrap-typeahead/css/Typeahead.css'
import 'react-bootstrap-typeahead/css/Typeahead.bs5.css'
import '../../assets/TypeaheadOverride.scss'
import { useApi } from '@src/ApiContext.ts'

export interface TypeToFindDeviceTypeProps {
  id: string
  onSelect?: (deviceTypeId: string[]) => void
  selected?: string[]
  asyncProps?: Partial<UseAsyncProps>
}

interface CurrentlyLoading {
  opts: ListDeviceTypesRequest
  promise: Promise<DeviceTypeSearchResults>
  abort: AbortController
}

export const TypeToFindDeviceType: FC<TypeToFindDeviceTypeProps> = ({ id, onSelect, selected, asyncProps }) => {
  const api = useApi()
  const [loading, setLoading] = useState(false)
  const [options, setOptions] = useState<DeviceTypeResult[]>([])
  const [currentSelection, setCurrentSelection] = useState<DeviceTypeResult[]>()
  const currentlyLoading = useRef<CurrentlyLoading>()

  function onSearch(q: string) {
    void loadSearchResults({
      deviceTypeSearchBody: {
        modelRegex: q,
      },
    })
  }

  function onChange(options: DeviceTypeResult[]) {
    setCurrentSelection(options)
    if (onSelect) {
      onSelect(options.map((o) => o.id))
    }
  }

  const loadSearchResults = useCallback(
    async (opts: ListDeviceTypesRequest) => {
      if (!opts.deviceTypeSearchBody) opts.deviceTypeSearchBody = {}
      opts.deviceTypeSearchBody.fields = ['model', 'manufacturer']
      const abort = new AbortController()
      if (currentlyLoading.current?.opts && isEqual(currentlyLoading.current.opts, opts)) {
        return currentlyLoading.current.promise
      } else if (currentlyLoading.current?.opts) {
        currentlyLoading.current.abort.abort()
      }
      setLoading(true)
      try {
        currentlyLoading.current = {
          opts,
          abort,
          promise: api.deviceApi.listDeviceTypes(opts, { signal: abort.signal }),
        }
        const response = await currentlyLoading.current.promise
        setOptions(response.items)
        return response
      } catch (e) {
        console.error(e)
        return null
      } finally {
        setLoading(false)
      }
    },
    [api.deviceApi],
  )

  useEffect(() => {
    setCurrentSelection((current) => {
      if (isEqual(current?.map((v) => v.id).toSorted(), selected?.toSorted())) {
        return current
      }
      const newValue =
        selected?.map((v) => {
          return { id: v, version: 0, model: 'loading' }
        }) ?? []
      if (newValue.length > 0) {
        void loadSearchResults({ deviceTypeSearchBody: { ids: selected } }).then((results) => {
          setCurrentSelection(results?.items)
        })
      }
      return newValue
    })
  }, [selected, api.deviceApi, loadSearchResults, options])

  function renderMenuItemChildren(option: DeviceTypeResult, props: TypeaheadMenuProps) {
    return (
      <>
        <span className="manufacturer">{option.manufacturer?.displayName} </span>
        <Highlighter
          /* eslint-disable-next-line react/prop-types  */
          search={props.text}
        >
          {option.model ?? ''}
        </Highlighter>
      </>
    )
  }

  return (
    <span>
      <AsyncTypeahead
        id={id + 'Typeahead'}
        allowNew={false}
        isLoading={loading}
        onSearch={onSearch}
        options={options}
        useCache={true}
        onChange={(v) => {
          onChange(v as DeviceTypeResult[])
        }}
        labelKey={'model'}
        filterBy={() => true}
        selected={currentSelection ?? []}
        highlightOnlyResult={true}
        inputProps={{ id }}
        renderMenuItemChildren={(option, opts) => renderMenuItemChildren(option as DeviceTypeResult, opts)}
        {...asyncProps}
      />
    </span>
  )
}

export default TypeToFindDeviceType
