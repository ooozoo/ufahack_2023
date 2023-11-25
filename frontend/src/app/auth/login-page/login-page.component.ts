import { Component } from '@angular/core';
import {FormControl, FormGroup, Validators} from "@angular/forms";
import {AuthService} from "../services/auth.service";
import {ActivatedRoute, Router} from "@angular/router";

@Component({
  selector: 'app-login-page',
  templateUrl: './login-page.component.html',
  styleUrls: ['./login-page.component.less']
})
export class LoginPageComponent {
  error: string = '';

  username = new FormControl('', [Validators.required]);
  password = new FormControl('', [Validators.required]);

  form: FormGroup = new FormGroup({
    username: this.username,
    password: this.password,
  });

  constructor(
    private _authService: AuthService,
    private _router: Router,
    private _route: ActivatedRoute,
  ) { }

  submit(): void {
    this.error = '';
    if (this.form.invalid) {
      return;
    }

    this._authService.login(
      this.form.value.username,
      this.form.value.password
    ).subscribe(
      () => {
        this.form.reset();
        this._router.navigate(['/']);
      },
      () => {
        this.error = 'Неверное имя учетной записи или пароль.';
      },
    );
  }
}
