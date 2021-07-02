import React from 'react';
import {Navigate} from 'react-router-dom';
import NavLayout from 'src/layouts/NavLayout';
import MainLayout from 'src/layouts/MainLayout';
import DashboardView from 'src/views/index/DashboardView';
import CurrentView from 'src/views/index/CurrentView';
import EnvironmentView from 'src/views/environment/EnvironmentView';
import EnvironmentSpaceView from 'src/views/environment/EnvironmentSpaceView';
import ProjectView from 'src/views/project/ProjectView';
import DeployView from 'src/views/project/DeployView';
import DocumentView from 'src/views/document/DocumentView';
import XtermView from 'src/views/xterm/XtermView';
import NotFoundView from 'src/views/errors/NotFoundView';

import {
  Apps as AppsIcon,
  Dashboard as DashboardIcon,
  GroupAdd as RoleIcon,
  InsertChart as DeployIcon,
  Kitchen as EnvIcon,
  People as UserIcon,
} from '@material-ui/icons';

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
      {path: 'configs', element: <DashboardView/>},
      {path: 'accounts', element: <DashboardView/>},
      {path: 'roles', element: <DashboardView/>},
      {path: 'terminal', element: <XtermView/>},
      {path: 'document/:app/:key/:space', element: <DocumentView/>},
      {path: '*', element: <Navigate to="/404"/>}
    ]
  },
  {
    path: '/',
    element: <MainLayout/>,
    children: [
      {path: '404', element: <NotFoundView/>},
      {path: '/', element: <Navigate to="/app/dashboard"/>},
      {path: '*', element: <Navigate to="/404"/>}
    ]
  }
];

const routeItems = [
  {
    href: '/app/dashboard',
    icon: DashboardIcon,
    title: '面板'
  },
  {
    href: '/app/environment',
    icon: EnvIcon,
    title: '环境'
  },
  {
    href: '/app/projects',
    icon: AppsIcon,
    title: '项目'
  },
  {
    href: '/app/deploys',
    icon: DeployIcon,
    title: '部署'
  },
  {
    href: '/app/roles',
    icon: RoleIcon,
    title: '角色'
  },
  {
    href: '/app/accounts',
    icon: UserIcon,
    title: '用户'
  },
];

export {routes, routeItems};
