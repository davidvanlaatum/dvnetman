import { FC } from 'react'

interface ErrorPageProps {
  error: unknown
}

const ErrorPage: FC<ErrorPageProps> = ({ error }) => {
  if (!(error instanceof Error)) {
    return (
      <div>
        <h1>Unknown Error</h1>
        <p>{JSON.stringify(error)}</p>
      </div>
    )
  }
  return (
    <div>
      <h1>{error.constructor.name}</h1>
      <p>{error.message}</p>
    </div>
  )
}

export default ErrorPage
