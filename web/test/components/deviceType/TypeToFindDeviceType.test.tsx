import { describe, expect, it } from 'vitest'
import { render, screen, waitFor } from '@testing-library/react'
import TypeToFindDeviceType, { TypeToFindDeviceTypeProps } from '@src/components/deviceType/TypeToFindDeviceType.tsx'
import { ApiProvider } from '@src/ApiProvider.tsx'
import userEvent from '@testing-library/user-event'
import { FakeFetch } from '@test/helpers/fakeFetch.ts'
import { FC, useState } from 'react'
import { DeviceTypeSearchResults } from '@src/api'

const TestComponent: FC<Partial<TypeToFindDeviceTypeProps>> = ({ selected, ...props }) => {
  const [ids, setIds] = useState<string[]>(selected ?? [])
  return (
    <>
      <TypeToFindDeviceType
        id="deviceType"
        {...props}
        onSelect={(id) => {
          setIds(id)
        }}
        selected={ids}
      />
      <div data-testid="selected-ids">{JSON.stringify(ids)}</div>
    </>
  )
}

describe('<TypeToFindDeviceType />', () => {
  it('should render correctly', async () => {
    const fetch = new FakeFetch()
    fetch.expectDeviceTypeSearch({
      requestBody: {
        modelRegex: 'abc',
        fields: ['model', 'manufacturer'],
      },
      body: {
        items: [{ id: '1', model: 'abc', manufacturer: { id: 'm1', displayName: 'def' }, version: 1 }],
        count: 1,
        next: false,
      },
    })
    render(
      <ApiProvider apiConfig={{ basePath: 'http://localhost', fetchApi: fetch.fetch }}>
        <TestComponent asyncProps={{ delay: 0, minLength: 3 }} />
      </ApiProvider>,
    )
    const input = screen.getByRole('combobox')
    expect(input).not.toBeNull()
    await userEvent.click(input)
    await userEvent.type(input, 'abc')
    expect(screen.getByText('abc')).not.toBeNull()
    await userEvent.type(input, '{enter}')
    expect(screen.getByTestId('selected-ids').textContent).toEqual(JSON.stringify(['1']))
  })

  it('should render initial value correctly', async () => {
    const fetch = new FakeFetch()
    let resolve: (value: DeviceTypeSearchResults) => void = () => {
      throw new Error('not implemented')
    }
    const p = new Promise<DeviceTypeSearchResults>((r) => {
      resolve = r
    })
    fetch.expectDeviceTypeSearch({
      requestBody: {
        ids: ['1'],
        fields: ['model', 'manufacturer'],
      },
      callback: async () => p,
    })
    render(
      <ApiProvider apiConfig={{ basePath: 'http://localhost', fetchApi: fetch.fetch }}>
        <TestComponent asyncProps={{ delay: 0, minLength: 3 }} selected={['1']} />
      </ApiProvider>,
    )
    const input = screen.getByRole('combobox')
    expect(input).not.toBeNull()
    expect(input).toHaveValue('loading')
    resolve({
      items: [{ id: '1', model: 'abc', manufacturer: { id: 'm1', displayName: 'def' }, version: 1 }],
      count: 1,
      next: false,
    })
    await waitFor(() => {
      expect(input).not.toHaveValue('loading')
      expect(input).toHaveValue('abc')
    })
  })
})
