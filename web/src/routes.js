import React from 'react'
import {Navigate} from 'react-router-dom'
import NavLayout from 'src/layouts/NavLayout'
import MainLayout from 'src/layouts/MainLayout'
import SimpleLayout from 'src/layouts/SimpleLayout'
import NoneLayout from 'src/layouts/NoneLayout'
import DashboardView from 'src/views/index/DashboardView'
import CurrentView from 'src/views/index/CurrentView'
import EnvironmentView from 'src/views/environment/EnvironmentView'
import EnvironmentSpaceView from 'src/views/environment/EnvironmentSpaceView'
import ProjectView from 'src/views/project/ProjectView'
import DeployView from 'src/views/project/DeployView'
import AccountView from 'src/views/account/AccountView'
import RoleView from 'src/views/account/RoleView'
import ToolView from 'src/views/tool/ToolView'
import DocumentView from 'src/views/document/DocumentView'
import XtermView from 'src/views/xterm/XtermView'
import LoginView from 'src/views/login/LoginView'
import NotFoundView from 'src/views/errors/NotFoundView'

const routes = [
  {
    path: 'app',
    element: <NavLayout/>,
    children: [
      {path: 'dashboard', element: <DashboardView/>},
      {path: 'current', element: <CurrentView/>},
      {path: 'environment', element: <EnvironmentView/>},
      {path: 'spaces', element: <EnvironmentSpaceView/>},
      {path: 'projects', element: <ProjectView/>},
      {path: 'templates', element: <DashboardView/>},
      {path: 'deploys', element: <DeployView/>},
      {path: 'plans', element: <DashboardView/>},
      {path: 'configs', element: <DashboardView/>},
      {path: 'accounts', element: <AccountView/>},
      {path: 'roles', element: <RoleView/>},
      {path: 'tools', element: <ToolView/>},
      {path: '*', element: <Navigate to="/404"/>}
    ]
  },
  {
    path: 'app',
    element: <NoneLayout/>,
    children: [
      {path: 'terminal', element: <XtermView/>},
      {path: '*', element: <Navigate to="/404"/>}
    ]
  },
  {
    path: 'app',
    element: <SimpleLayout/>,
    children: [
      {path: 'document/:app/:key/:space', element: <DocumentView/>},
      {path: '*', element: <Navigate to="/404"/>}
    ]
  },
  {
    path: '/',
    element: <MainLayout/>,
    children: [
      {path: '404', element: <NotFoundView/>},
      {path: 'login', element: <LoginView/>},
      {path: '/', element: <Navigate to="/app/dashboard"/>},
      {path: '*', element: <Navigate to="/404"/>}
    ]
  }
]

export {routes}
