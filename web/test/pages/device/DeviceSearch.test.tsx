import { describe, expect, it } from 'vitest'
import { render, screen, waitFor } from '@testing-library/react'
import DeviceSearch from '@src/pages/device/DeviceSearch.tsx'
import { ApiProvider } from '@src/ApiProvider.tsx'
import { FakeFetch } from '@test/helpers/fakeFetch.ts'
import { MemoryRouter } from 'react-router'

describe('DeviceSearch', () => {
  it('should render correctly - none found', async () => {
    const fetch = new FakeFetch()
    fetch.expectDeviceSearch({
      requestBody: {
        nameRegex: 'test',
      },
      body: {
        items: [],
        count: 0,
        next: false,
      },
    })
    render(
      <MemoryRouter initialEntries={['/search?nameRegex=test']}>
        <ApiProvider apiConfig={{ basePath: 'http://localhost', fetchApi: fetch.fetch }}>
          <DeviceSearch />
        </ApiProvider>
      </MemoryRouter>,
    )
    await waitFor(() => {
      expect(screen.getByText('No devices found')).not.toBeNull()
    })
  })

  it('should render correctly', async () => {
    const fetch = new FakeFetch()
    fetch.expectDeviceSearch({
      requestBody: {
        nameRegex: 'test',
      },
      body: {
        items: [
          {
            id: '1',
            name: 'test-device',
            deviceType: { id: '1', displayName: 'test-type' },
            version: 1,
          },
        ],
        count: 0,
        next: false,
      },
    })
    render(
      <MemoryRouter initialEntries={['/search?nameRegex=test']}>
        <ApiProvider apiConfig={{ basePath: 'http://localhost', fetchApi: fetch.fetch }}>
          <DeviceSearch />
        </ApiProvider>
      </MemoryRouter>,
    )
    await waitFor(() => {
      expect(screen.queryByText('test-device')).not.toBeNull()
    })
  })
})
