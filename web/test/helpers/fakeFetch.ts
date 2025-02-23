import { DeviceTypeSearchResults } from '@src/api'
import { assert, onTestFinished } from 'vitest'
import { isEqual } from 'lodash'
import { APIErrorModal } from '@src/api/models/APIErrorModal.ts'

export interface ExpectedResponse<T> {
  body?: T | APIErrorModal | Response | Error
  status?: number
  callback?: (url: URL, init: RequestInit, log: string[]) => Promise<T | Response | APIErrorModal | Error>
}

export interface ExpectedDeviceTypeSearch extends ExpectedResponse<DeviceTypeSearchResults> {
  modelRegex?: string[]
  ids?: string[]
  fields?: string[]
  seenCount?: number
}

export class FakeFetch {
  private readonly unexpectedCalls: string[] = []
  private readonly expectedDeviceTypeSearches: ExpectedDeviceTypeSearch[] = []

  public expectDeviceTypeSearch(expected: ExpectedDeviceTypeSearch) {
    this.expectedDeviceTypeSearches.push(expected)
  }

  private async response<T>(res: ExpectedResponse<T>, url: URL, init: RequestInit, log: string[]): Promise<Response> {
    const x = { body: res.body, status: res.status }
    if (res.callback) {
      const r = await res.callback(url, init, log)
      if (r instanceof Response) {
        return r
      }
      x.body = r
    }
    return new Response(JSON.stringify(x.body), {
      status: x.status ?? 200,
      headers: { 'Content-Type': 'application/json' },
    })
  }

  private async deviceType(url: URL, init: RequestInit, log: string[]): Promise<Response | null> {
    if (url.pathname == '/api/v1/deviceType' && init?.method == 'GET') {
      const modelRegex = url.searchParams.getAll('modelRegex')
      const fields = url.searchParams.getAll('fields')
      for (const v of this.expectedDeviceTypeSearches) {
        if (v.modelRegex && !isEqual(v.modelRegex, modelRegex)) {
          log.push('modelRegex mismatch')
          continue
        }
        if (v.fields && !isEqual(v.fields, fields)) {
          log.push('fields mismatch')
          continue
        }
        if (v.ids && !isEqual(v.ids, url.searchParams.getAll('ids'))) {
          log.push('ids mismatch')
          continue
        }
        v.seenCount = (v.seenCount ?? 0) + 1
        return this.response(v, url, init, log)
      }
      log.push('unexpected device type search')
    }
    return null
  }

  get fetch(): (input: RequestInfo, init?: RequestInit) => Promise<Response> {
    onTestFinished(() => {
      if (this.unexpectedCalls.length > 0) {
        assert.fail(this.unexpectedCalls.join('\n'))
      }
      for (const v of this.expectedDeviceTypeSearches) {
        if (v.seenCount === undefined) {
          assert.fail('expected device type search not called')
        }
      }
    })
    return async (input, init) => {
      // console.trace('fetch', input, init)
      const url = new URL(input as string)
      const log: string[] = []
      let response: Response | null = null
      if (url.pathname.startsWith('/api/v1/deviceType')) {
        response = await this.deviceType(url, init as RequestInit, log)
      } else {
        log.push('unknown path')
      }
      if (response) {
        return response
      }
      this.unexpectedCalls.push(init?.method + ' ' + (input as string) + ': ' + log.join(', '))
      return new Response(
        JSON.stringify({
          errors: [
            {
              code: 'not_found',
              message: 'not found',
            },
          ],
        }),
        {
          status: 404,
          headers: { 'Content-Type': 'application/json' },
        },
      )
    }
  }
}
