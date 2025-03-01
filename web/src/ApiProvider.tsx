import { FC, ReactNode, useMemo } from 'react'
import { ConfigurationParameters } from './api'
import { Api, ApiContext } from '@src/ApiContext.ts'

export const ApiProvider: FC<{ children: ReactNode; apiConfig?: ConfigurationParameters }> = ({
  children,
  apiConfig,
}) => {
  const api = useMemo(() => {
    return new Api(apiConfig)
  }, [apiConfig])
  return <ApiContext.Provider value={api}>{children}</ApiContext.Provider>
}
