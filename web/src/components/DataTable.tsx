import { Dropdown, FormCheck, Spinner, Table } from 'react-bootstrap'
import 'bootstrap-icons/font/bootstrap-icons.css'
import { FC, ReactNode, useEffect, useState } from 'react'
import './DataTable.scss'

export interface DataTableColumnProps {
  id: string
  label: string
  sortable?: boolean
  render?: (row: any) => ReactNode
  sort?: (a: any, b: any, dir: boolean) => 0 | 1 | -1
  getDisplayValue?: (row: any) => string
  width?: string
}

export interface DataTableRow {
  id: string

  [key: string]: any
}

export interface DataTableProps {
  columns: DataTableColumnProps[]
  data: DataTableRow[]
  onSort?: (column: string | undefined, dir: boolean) => void
  selectable?: boolean
  selected?: string[]
  onSelect?: (selected: string[]) => void
  loading?: boolean
  renderOnNoData?: ReactNode
}

function toString(value: unknown): string {
  if (value === null || value === undefined) {
    return ''
  }
  switch (typeof value) {
    case 'bigint':
      return value.toString()
    case 'string':
      return value
    case 'number':
      return value.toString()
    case 'boolean':
      return value ? 'true' : 'false'
    case 'object':
      return JSON.stringify(value)
    default:
      throw new Error(`Unknown type ${typeof value}`)
  }
}

function getColumnValueFunc(column: DataTableColumnProps): (row: DataTableRow) => string {
  if (column.getDisplayValue) {
    return column.getDisplayValue
  }
  return (row: DataTableRow) => toString(row[column.id])
}

function getSortFunc(column: DataTableColumnProps, sortDir: boolean): (a: DataTableRow, b: DataTableRow) => number {
  const getValue = getColumnValueFunc(column)
  let sortFunc = (a: DataTableRow, b: DataTableRow) => {
    const aValue = getValue(a)
    const bValue = getValue(b)
    if (aValue < bValue) {
      return sortDir ? -1 : 1
    }
    if (aValue > bValue) {
      return sortDir ? 1 : -1
    }
    return 0
  }
  if (column.sort) {
    const f = column.sort
    sortFunc = (a, b) => f(a, b, sortDir)
  }
  return sortFunc
}

export const DataTable: FC<DataTableProps> = ({
  columns,
  data,
  onSort,
  selectable,
  selected,
  onSelect,
  loading,
  renderOnNoData,
}) => {
  const [sortColumn, setSortColumn] = useState<string | null>(null)
  const [sortDir, setSortDir] = useState<boolean>(true)
  const [renderData, setRenderData] = useState<DataTableRow[]>(data)

  useEffect(() => {
    const sortColumnObj = columns.find((column) => column.id == sortColumn)
    if (sortColumn && sortColumnObj) {
      const dataCopy = [...data]
      dataCopy.sort(getSortFunc(sortColumnObj, sortDir))
      setRenderData(dataCopy)
    } else {
      setRenderData(data)
    }
  }, [columns, data, sortColumn, sortDir])

  function onSelected(checked: boolean, row: DataTableRow) {
    if (onSelect) {
      if (checked) {
        onSelect([...(selected ?? []), row.id])
      } else {
        onSelect((selected ?? []).filter((id) => id != row.id))
      }
    }
  }

  function onSelectAll() {
    if (onSelect) {
      if (isAllSelected()) {
        onSelect([])
      } else {
        onSelect(data.map((row) => row.id))
      }
    }
  }

  function isAllSelected() {
    return data.length > 0 && selected?.length == data.length
  }

  function onSortClick(id: string) {
    if (sortColumn == id) {
      if (sortDir) {
        setSortDir(false)
      } else {
        setSortColumn(null)
        setSortDir(true)
      }
    } else {
      setSortColumn(id)
      setSortDir(true)
    }
    if (onSort) {
      onSort(id, sortDir)
    }
  }

  function emptyBody(body: ReactNode) {
    return (
      <tbody>
        <tr>
          <td colSpan={columns.length}>{body}</td>
        </tr>
      </tbody>
    )
  }

  function renderBody(): ReactNode {
    if (loading) {
      return emptyBody(
        <>
          <Spinner size="sm" />
          Loading...
        </>,
      )
    }
    if (renderData.length == 0 && renderOnNoData) {
      return emptyBody(renderOnNoData)
    }
    return (
      <tbody>
        {renderData.map((row) => (
          <tr key={row.id}>{columns.map((column, index) => renderBodyCell(row, column, index))}</tr>
        ))}
      </tbody>
    )
  }

  function renderBodyCell(row: DataTableRow, column: DataTableColumnProps, index: number) {
    let value: ReactNode
    if (column.render) {
      value = column.render(row)
    } else {
      value = getColumnValueFunc(column)(row)
    }
    return (
      <td key={column.id} className={'data-table-column'}>
        <div>
          {selectable && index == 0 && (
            <FormCheck
              className={'float-start'}
              checked={selected?.includes(row.id)}
              onChange={(e) => {
                onSelected(e.target.checked, row)
              }}
            />
          )}
          {value}
        </div>
      </td>
    )
  }

  return (
    <Table striped bordered hover>
      <thead>
        <tr>
          {columns.map((column, index) => {
            const sorted = column.id == sortColumn
            let ariaLabel: string | undefined
            if (sorted) {
              ariaLabel = sortDir ? 'sorted ascending' : 'sorted descending'
            }
            return (
              <th key={column.id} style={{ width: column.width }} className={'data-table-column'}>
                <div>
                  {selectable && index == 0 && <FormCheck onChange={onSelectAll} checked={isAllSelected()} />}
                  {!column.sortable && <div>column.label</div>}
                  {column.sortable && (
                    <button
                      className={['data-table-sorter', 'fill', [sorted ? 'data-table-sorted' : null]].join(' ')}
                      onClick={() => {
                        onSortClick(column.id)
                      }}
                    >
                      {column.label}
                      <i
                        className={'bi bi-sort-alpha-' + (sortDir || !sorted ? 'up' : 'down') + ' sort-dir'}
                        aria-label={ariaLabel}
                      />
                    </button>
                  )}
                  <div className={'float-end'}>
                    <Dropdown>
                      <Dropdown.Toggle variant="success" as={'span'}></Dropdown.Toggle>
                      <Dropdown.Menu>
                        <Dropdown.Item href="#/action-1">Action</Dropdown.Item>
                        <Dropdown.Item href="#/action-2">Another action</Dropdown.Item>
                        <Dropdown.Item href="#/action-3">Something else</Dropdown.Item>
                      </Dropdown.Menu>
                    </Dropdown>
                  </div>
                </div>
              </th>
            )
          })}
        </tr>
      </thead>
      {renderBody()}
    </Table>
  )
}

export default DataTable
