import { NgDompurifySanitizer } from "@tinkoff/ng-dompurify";
import { TuiRootModule, TuiDialogModule, TuiAlertModule, TUI_SANITIZER } from "@taiga-ui/core";
import { BrowserAnimationsModule } from "@angular/platform-browser/animations";
import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import {HTTP_INTERCEPTORS, HttpClientModule} from '@angular/common/http';


import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { PageNotFoundComponent } from './components/page-not-found/page-not-found.component';
import { LayoutsModule } from './layouts/layouts.module';
import {TuiBlockStatusModule} from "@taiga-ui/layout";
import {AuthInterceptor} from "./auth/interceptors/auth.interceptor";


@NgModule({
  declarations: [
    AppComponent,
    PageNotFoundComponent,
  ],
    imports: [
      BrowserModule,
      AppRoutingModule,
      BrowserAnimationsModule,
      TuiRootModule,
      TuiDialogModule,
      TuiAlertModule,
      LayoutsModule,
      TuiBlockStatusModule,
      HttpClientModule,
    ],
  providers: [
    {
      provide: TUI_SANITIZER,
      useClass: NgDompurifySanitizer
    },
    {
      provide: HTTP_INTERCEPTORS,
      useClass: AuthInterceptor,
      multi: true,
    },
  ],
  bootstrap: [AppComponent]
})
export class AppModule { }
