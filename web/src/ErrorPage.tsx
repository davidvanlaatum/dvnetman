import { FC } from 'react'

interface ErrorPageProps {
  error: Error
}

const ErrorPage: FC<ErrorPageProps> = ({ error }) => {
  return (
    <div>
      <h1>{error.constructor.name}</h1>
      <p>{error.message}</p>
    </div>
  )
}

export default ErrorPage
