import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { PageNotFoundComponent } from "./components/page-not-found/page-not-found.component";
import {TuiBlockStatusModule} from "@taiga-ui/layout";



@NgModule({
  declarations: [
    PageNotFoundComponent
  ],
  imports: [
    CommonModule,
    TuiBlockStatusModule
  ]
})
export class AppModule { }
