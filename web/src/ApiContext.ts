import {
  APIErrorModal,
  APIErrorModalFromJSON,
  Configuration,
  ConfigurationParameters,
  DeviceApi,
  ErrorContext,
  FetchParams,
  ManufacturerApi,
  Middleware,
  RequestContext,
  ResponseContext,
  StatsApi,
  UserApi,
} from '@src/api'
import { createContext, useContext } from 'react'

export class APIError extends Error {
  response: Response
  error: APIErrorModal

  constructor(response: Response, error: APIErrorModal) {
    super(response.statusText)
    this.response = response
    this.error = error
  }

  get message() {
    return this.error.errors?.map((e) => e.message).join(', ') ?? this.response.statusText
  }
}

export class NotFoundError extends APIError {}

export class UnauthorizedError extends APIError {}

export class ForbiddenError extends APIError {}

export class ConflictError extends APIError {}

export class ErrorTransformer implements Middleware {
  private readonly onUnauthorizedCallbacks: (() => void)[]

  constructor(onUnauthorizedCallbacks: (() => void)[]) {
    this.onUnauthorizedCallbacks = onUnauthorizedCallbacks
  }

  async post(context: ResponseContext): Promise<Response> {
    if (context.response.status >= 400) {
      const json = (await context.response.json()) as object
      switch (context.response.status) {
        case 401:
          this.onUnauthorizedCallbacks.forEach((fn) => {
            fn()
          })
          return Promise.reject(new UnauthorizedError(context.response, APIErrorModalFromJSON(json)))
        case 403:
          throw new ForbiddenError(context.response, APIErrorModalFromJSON(json))
        case 404:
          throw new NotFoundError(context.response, APIErrorModalFromJSON(json))
        case 409:
          throw new ConflictError(context.response, APIErrorModalFromJSON(json))
        default:
          throw new APIError(context.response, APIErrorModalFromJSON(json))
      }
    }
    return context.response
  }
}

export class OneAtATime implements Middleware {
  private requestInProgress = false
  private readonly queue: (() => void)[] = []

  async pre(context: RequestContext): Promise<FetchParams> {
    if (this.requestInProgress) {
      return new Promise<FetchParams>((resolve) => {
        this.queue.push(() => {
          resolve(this.pre(context))
        })
      })
    }
    this.requestInProgress = true
    return { url: context.url, init: context.init }
  }

  // eslint-disable-next-line @typescript-eslint/no-invalid-void-type
  post(context: ResponseContext): Promise<Response | void> {
    this.complete()
    return Promise.resolve(context.response)
  }

  // eslint-disable-next-line @typescript-eslint/no-invalid-void-type
  onError(context: ErrorContext): Promise<Response | void> {
    this.complete()
    return Promise.resolve(context.response)
  }

  private complete() {
    this.requestInProgress = false
    const next = this.queue.shift()
    if (next) {
      next()
    }
  }
}

export class Api {
  deviceApi: DeviceApi
  manufacturerApi: ManufacturerApi
  userApi: UserApi
  statsApi: StatsApi
  private readonly onUnauthorizedCallbacks: (() => void)[] = []

  constructor(apiConfig?: ConfigurationParameters) {
    const config = new Configuration({
      basePath: window.location.origin,
      ...apiConfig,
      middleware: [
        new OneAtATime(),
        new ErrorTransformer(this.onUnauthorizedCallbacks),
        ...(apiConfig?.middleware ?? []),
      ],
    })
    this.deviceApi = new DeviceApi(config)
    this.manufacturerApi = new ManufacturerApi(config)
    this.userApi = new UserApi(config)
    this.statsApi = new StatsApi(config)
  }

  public onUnauthorized(fn: () => void): () => void {
    this.onUnauthorizedCallbacks.push(fn)
    return () => this.onUnauthorizedCallbacks.splice(this.onUnauthorizedCallbacks.indexOf(fn), 1)
  }
}

export const ApiContext = createContext<Api | null>(null)

export const useApi = () => {
  const context = useContext(ApiContext)
  if (!context) {
    throw new Error('useApi must be used within a ApiProvider')
  }
  return context
}
