import { Routes } from '@angular/router';

export const routes: Routes = [
  {
    path: '',
    component: LayoutPageComponent,
    canActivate: [AuthGuard],
    canActivateChild: [AuthGuard],
    children: [
      {
        path: '',
        loadChildren: () => import('./quiz/quiz.module').then(m => m.QuizModule),
      }
    ]
  },
  {
    path: '',
    canActivate: [LoginGuard],
    canActivateChild: [LoginGuard],
    loadChildren: () => import('./auth/auth.module').then(m => m.AuthModule),
  },
  {
    path: '**',
    component: NotFoundPageComponent,
  },
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutes { }
