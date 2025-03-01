import './App.scss'
import { BrowserRouter as Router, Route, Routes } from 'react-router'
import { Nav, Navbar, NavDropdown } from 'react-bootstrap'
import { lazy, Suspense } from 'react'
import { ApiProvider } from './ApiContext.tsx'
import { UserProvider } from './UserContext.tsx'
import CurrentUser from '@src/components/CurrentUser.tsx'

const About = lazy(() => import('./About.tsx'))
const Swagger = lazy(() => import('./Swagger.tsx'))
const DeviceAdd = lazy(() => import('./pages/device/DeviceAdd.tsx'))
const DeviceDetail = lazy(() => import('./pages/device/DeviceDetail.tsx'))
const DeviceSearch = lazy(() => import('./pages/device/DeviceSearch.tsx'))
const Index = lazy(() => import('./Index.tsx'))

function App() {
  const basePath = import.meta.env.BASE_URL
  return (
    <ApiProvider>
      <UserProvider>
        <Router>
          <Navbar expand="lg" className="bg-body-tertiary">
            <Navbar.Brand href={`${basePath}/`}>DVNetMan</Navbar.Brand>
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
              </Nav>
            </Navbar.Collapse>
            <Navbar.Collapse className={'justify-content-end me-3'}>
              <NavDropdown title={<CurrentUser display="icon" />} id="user-dropdown" align={'end'}>
                <NavDropdown.Header>
                  <CurrentUser display="name" />
                </NavDropdown.Header>
                <NavDropdown.Item href={`${basePath}/profile`}>Profile</NavDropdown.Item>
                <NavDropdown.Item href={`/auth/logout`}>Logout</NavDropdown.Item>
              </NavDropdown>
            </Navbar.Collapse>
          </Navbar>
          <Suspense fallback={<div>Loading...</div>}>
            <Routes>
              <Route path={`${basePath}/`} Component={Index} />
              <Route path={`${basePath}/about`} Component={About} />
              <Route path={`${basePath}/swagger`} Component={Swagger} />
              <Route path={`${basePath}/device/add`} Component={DeviceAdd} />
              <Route path={`${basePath}/device/search`} Component={DeviceSearch} />
              <Route path={`${basePath}/device/:uuid`} Component={DeviceDetail} />
              <Route path="*" element={<h1>Not Found</h1>} />
            </Routes>
          </Suspense>
        </Router>
      </UserProvider>
    </ApiProvider>
  )
}

export default App
