import { Routes, RouterModule } from '@angular/router';
import { NgModule } from '@angular/core';

import { PageNotFoundComponent } from "./components/page-not-found/page-not-found.component";

export const routes: Routes = [
  // {
  //   path: '',
  //   component: LayoutPageComponent,
  //   canActivate: [AuthGuard],
  //   canActivateChild: [AuthGuard],
  //   children: [
  //     {
  //       path: '',
  //       loadChildren: () => import('./quiz/quiz.module').then(m => m.QuizModule),
  //     }
  //   ]
  // },
  // {
  //   path: '',
  //   canActivate: [LoginGuard],
  //   canActivateChild: [LoginGuard],
  //   loadChildren: () => import('./auth/auth.module').then(m => m.AuthModule),
  // },
  {
    path: '**',
    component: PageNotFoundComponent,
  },
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutes { }
