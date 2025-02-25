import { DeviceSearchBody, DeviceSearchResults, DeviceTypeSearchBody, DeviceTypeSearchResults } from '@src/api'
import { assert, onTestFinished } from 'vitest'
import { isEqual } from 'lodash'
import { APIErrorModal } from '@src/api/models/APIErrorModal.ts'

export interface ExpectedRequest<Req, Res> {
  url: URL
  method: string
  requestBody?: Req
  headers?: Headers
  seenCount?: number
  body?: Res | APIErrorModal | Response | Error
  status?: number
  callback?: (url: URL, init: RequestInit, log: string[]) => Promise<Res | Response | APIErrorModal | Error>
}

export type ExpectedDeviceTypeSearch = ExpectedRequest<DeviceTypeSearchBody, DeviceTypeSearchResults>

export type ExpectedDeviceSearch = ExpectedRequest<DeviceSearchBody, DeviceSearchResults>

async function getBodyAsJson<T>(init: RequestInit): Promise<T | null> {
  if (!init.body) {
    return null
  } else {
    return new Response(init.body).json()
  }
}

function headersMatch(expected: Headers, actual: Headers): boolean {
  for (const [key, value] of expected.entries()) {
    if (actual.get(key) !== value) {
      return false
    }
  }
  return true
}

export class FakeFetch {
  private readonly unexpectedCalls: string[] = []
  private readonly expectedRequests: ExpectedRequest<any, any>[] = []
  // private readonly expectedDeviceTypeSearches: ExpectedDeviceTypeSearch[] = []
  // private readonly expectedDeviceSearches: ExpectedDeviceSearch[] = []

  private async findMatchingRequest(
    url: URL,
    init: RequestInit,
    log: string[],
  ): Promise<ExpectedRequest<any, any> | null> {
    const body = await getBodyAsJson(init)
    for (const v of this.expectedRequests) {
      if (v.url.href !== url.href) {
        log.push('url mismatch ' + v.url.href + ' ' + url.href)
        continue
      }
      if (v.method !== init.method) {
        log.push('method mismatch ' + v.method + ' ' + init.method)
        continue
      }
      if (v.headers && !headersMatch(v.headers, new Headers(init.headers))) {
        log.push('headers mismatch')
        continue
      }
      if (v.requestBody && !isEqual(v.requestBody, body)) {
        log.push('request body mismatch ' + JSON.stringify(v.requestBody) + ' ' + JSON.stringify(body))
        continue
      }
      v.seenCount = (v.seenCount ?? 0) + 1
      return v
    }
    return null
  }

  public expectDeviceTypeSearch(expected: Omit<ExpectedDeviceTypeSearch, 'url' | 'method'>) {
    this.expectedRequests.push({
      url: new URL('http://localhost/api/v1/deviceType/search'),
      method: 'POST',
      ...expected,
    })
  }

  public expectDeviceSearch(expected: Omit<ExpectedDeviceSearch, 'url' | 'method'>) {
    this.expectedRequests.push({
      url: new URL('http://localhost/api/v1/device/search'),
      method: 'POST',
      ...expected,
    })
  }

  private async response<Req, Res>(
    res: ExpectedRequest<Req, Res>,
    url: URL,
    init: RequestInit,
    log: string[],
  ): Promise<Response> {
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

  get fetch(): (input: RequestInfo, init?: RequestInit) => Promise<Response> {
    onTestFinished(() => {
      if (this.unexpectedCalls.length > 0) {
        assert.fail(this.unexpectedCalls.join('\n'))
      }
      for (const v of this.expectedRequests) {
        if (v.seenCount === undefined) {
          assert.fail('expected request not called')
        }
      }
    })
    return async (input, init) => {
      // console.trace('fetch', input, init)
      const url = new URL(input as string)
      const log: string[] = []
      let response: Response | null = null
      const res = await this.findMatchingRequest(url, init as RequestInit, log)
      if (res) {
        response = await this.response(res, url, init as RequestInit, log)
      }
      if (response) {
        return response
      }
      console.error('unexpected fetch', input, init, log)
      this.unexpectedCalls.push(init?.method + ' ' + (input as string) + ': ' + log.join(', '))
      return new Response(
        JSON.stringify({
          errors: [
            {
              code: 'not_found',
              message: 'not found',
            },
          ],
          log: log,
        }),
        {
          status: 404,
          headers: { 'Content-Type': 'application/json' },
        },
      )
    }
  }
}
