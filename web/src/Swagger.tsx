import SwaggerUI from 'swagger-ui-react'
import 'swagger-ui-react/swagger-ui.css'

function Swagger() {
  return (
    <div style={{ background: 'white', paddingTop: '1px' }}>
      <SwaggerUI url="/api/openapi.yaml" deepLinking={true} />
    </div>
  )
}

export default Swagger
