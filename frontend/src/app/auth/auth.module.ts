import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { TuiFieldErrorPipeModule, TuiInputModule, TuiInputPasswordModule } from '@taiga-ui/kit';
import { TuiErrorModule, TuiLoaderModule } from '@taiga-ui/core';

import { LoginPageComponent } from './login-page/login-page.component';
import { RegisterPageComponent } from './register-page/register-page.component';
import { LayoutsModule } from '../layouts/layouts.module';
import { AuthRoutingModule } from './auth-routing.module';
import { AuthService } from './services/auth.service'

@NgModule({
  declarations: [
    LoginPageComponent,
    RegisterPageComponent,
  ],
  imports: [
    CommonModule,
    LayoutsModule,
    AuthRoutingModule,
    TuiLoaderModule,
    TuiInputModule,
    TuiFieldErrorPipeModule,
    TuiInputPasswordModule,
    TuiErrorModule,
    ReactiveFormsModule,
    FormsModule,
  ],
  providers: [
    AuthService,
  ]
})
export class AuthModule { }
