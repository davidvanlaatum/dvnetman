import { createContext, FC, ReactNode, useContext, useMemo } from 'react'
import {
  Configuration,
  DeviceApi,
  Middleware,
  APIErrorModal,
  APIErrorModalFromJSON,
  ResponseContext,
  ManufacturerApi,
  StatsApi,
  UserApi,
  ConfigurationParameters,
} from './api'

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
      const json = await context.response.json()
      switch (context.response.status) {
        case 404:
          throw new NotFoundError(context.response, APIErrorModalFromJSON(json))
        default:
          throw new APIError(context.response, APIErrorModalFromJSON(json))
      }
    }
    return context.response
  }
}

export class Api {
  deviceApi: DeviceApi
  manufacturerApi: ManufacturerApi
  userApi: UserApi
  statsApi: StatsApi

  constructor(apiConfig?: ConfigurationParameters) {
    const config = new Configuration({
      basePath: window.location.origin,
      ...apiConfig,
      middleware: [new ErrorTransformer(), ...(apiConfig?.middleware ?? [])],
    })
    this.deviceApi = new DeviceApi(config)
    this.manufacturerApi = new ManufacturerApi(config)
    this.userApi = new UserApi(config)
    this.statsApi = new StatsApi(config)
  }
}

const ApiContext = createContext<Api | null>(null)

export const ApiProvider: FC<{ children: ReactNode; apiConfig?: ConfigurationParameters }> = ({
  children,
  apiConfig,
}) => {
  const api = useMemo(() => {
    return new Api(apiConfig)
  }, [apiConfig])
  return <ApiContext.Provider value={api}>{children}</ApiContext.Provider>
}

// eslint-disable-next-line react-refresh/only-export-components
export const useApi = () => {
  const context = useContext(ApiContext)
  if (!context) {
    throw new Error('useApi must be used within a ApiProvider')
  }
  return context
}
