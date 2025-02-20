import { createContext, useContext, useMemo } from 'react'
import { Configuration, DeviceApi, Middleware, APIErrorModal, APIErrorModalFromJSON, ResponseContext } from './api'

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

export class ErrorTransformer implements Middleware {
  async post(context: ResponseContext): Promise<Response | void> {
    if (context.response.status >= 400) {
      switch (context.response.status) {
        case 404:
          throw new NotFoundError(context.response, APIErrorModalFromJSON(await context.response.json()))
        default:
          throw new APIError(context.response, APIErrorModalFromJSON(await context.response.json()))
      }
    }
    return context.response
  }
}

export class Api {
  deviceApi: DeviceApi

  constructor() {
    const config = new Configuration({
      basePath: window.location.origin,
      middleware: [new ErrorTransformer()],
    })
    this.deviceApi = new DeviceApi(config)
  }
}

const ApiContext = createContext<Api | null>(null)

export const ApiProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const api = useMemo(() => {
    return new Api()
  }, [])
  return <ApiContext.Provider value={api}>{children}</ApiContext.Provider>
}

export const useApi = () => {
  const context = useContext(ApiContext)
  if (!context) {
    throw new Error('useDeviceApi must be used within a DeviceApiProvider')
  }
  return context
}
