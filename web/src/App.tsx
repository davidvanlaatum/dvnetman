import './App.css'
import { BrowserRouter as Router, Route, Routes } from 'react-router'
import { Nav, Navbar, NavDropdown } from 'react-bootstrap'
import { lazy, Suspense } from 'react'
import { ApiProvider } from './ApiContext.tsx'

const About = lazy(() => import('./About.tsx'))
const Swagger = lazy(() => import('./Swagger.tsx'))
const DeviceAdd = lazy(() => import('./device/DeviceAdd.tsx'))
const DeviceDetail = lazy(() => import('./device/DeviceDetail.tsx'))
const DeviceSearch = lazy(() => import('./device/DeviceSearch.tsx'))

function App() {
  const basePath = import.meta.env.BASE_URL
  return (
    <ApiProvider>
      <Router>
        <Navbar expand="lg" className="bg-body-tertiary">
          <Navbar.Brand href="#home">DVNetMan</Navbar.Brand>
          <Navbar.Toggle aria-controls="basic-navbar-nav" />
          <Navbar.Collapse id="basic-navbar-nav">
            <Nav className="me-auto">
              <Nav.Link href={`${basePath}/`}>Home</Nav.Link>
              <NavDropdown title="Organisation" id="organisation-dropdown">
                <NavDropdown.Header>Site</NavDropdown.Header>
                <NavDropdown.Item href={`${basePath}/site`}>List</NavDropdown.Item>
                <NavDropdown.Item href={`${basePath}/site/add`}>Add</NavDropdown.Item>
                <NavDropdown.Header>Location</NavDropdown.Header>
                <NavDropdown.Item href={`${basePath}/location`}>List</NavDropdown.Item>
                <NavDropdown.Item href={`${basePath}/location/add`}>Add</NavDropdown.Item>
                <NavDropdown.Header>Rack</NavDropdown.Header>
                <NavDropdown.Item href={`${basePath}/rack`}>List</NavDropdown.Item>
                <NavDropdown.Item href={`${basePath}/rack/add`}>Add</NavDropdown.Item>
              </NavDropdown>
              <NavDropdown title="Device" id="device-dropdown">
                <NavDropdown.Item href={`${basePath}/device/search`}>List</NavDropdown.Item>
                <NavDropdown.Item href={`${basePath}/device/add`}>Add</NavDropdown.Item>
              </NavDropdown>
              <NavDropdown title="IPAM" id="ipam-dropdown">
                <NavDropdown.Header>Prefix</NavDropdown.Header>
                <NavDropdown.Item href={`${basePath}/prefix`}>List</NavDropdown.Item>
                <NavDropdown.Item href={`${basePath}/prefix/add`}>Add</NavDropdown.Item>
              </NavDropdown>
              <Nav.Link href={`${basePath}/swagger`}>API</Nav.Link>
              <Nav.Link href={`/auth/test`}>Login</Nav.Link>
            </Nav>
          </Navbar.Collapse>
        </Navbar>
        <Suspense fallback={<div>Loading...</div>}>
          <Routes>
            <Route path={`${basePath}/`} Component={About} />
            <Route path={`${basePath}/about`} Component={About} />
            <Route path={`${basePath}/swagger`} Component={Swagger} />
            <Route path={`${basePath}/device/add`} Component={DeviceAdd} />
            <Route path={`${basePath}/device/search`} Component={DeviceSearch} />
            <Route path={`${basePath}/device/:uuid`} Component={DeviceDetail} />
            <Route path="*" element={<h1>Not Found</h1>} />
          </Routes>
        </Suspense>
      </Router>
    </ApiProvider>
  )
}

export default App
