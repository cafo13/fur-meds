import { Routes } from '@angular/router';
import { TabsPage } from './tabs.page';

export const routes: Routes = [
  {
    path: 'tabs',
    component: TabsPage,
    children: [
      {
        path: 'pets',
        loadComponent: () =>
          import('../pets/pets.page').then((m) => m.PetsPage),
      },
      {
        path: 'account',
        loadComponent: () =>
          import('../account/account.page').then((m) => m.AccountPage),
      },
      {
        path: 'settings',
        loadComponent: () =>
          import('../settings/settings.page').then((m) => m.SettingsPage),
      },
      {
        path: '',
        redirectTo: '/tabs/pets',
        pathMatch: 'full',
      },
    ],
  },
  {
    path: '',
    redirectTo: '/tabs/pets',
    pathMatch: 'full',
  },
];
