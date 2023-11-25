import { Component, OnInit } from '@angular/core';
import { FormControl, FormGroup, Validators } from '@angular/forms';
import { ActivatedRoute, Router } from '@angular/router';
import { AuthService } from '../services/auth.service';

@Component({
  selector: 'app-register-page',
  templateUrl: './register-page.component.html',
  styleUrls: ['./register-page.component.less']
})
export class RegisterPageComponent {
  registration: any = {
    loading: false,
  }
  error: string = '';

  username = new FormControl('', [Validators.required]);
  password = new FormControl('', [Validators.required]);

  form: FormGroup = new FormGroup({
    username: this.username,
    password: this.password,
  });

  get loader(): boolean {
    return this._loader;
  }

  set loader(value: boolean) {
    if (value) {
      this.username.disable();
      this.password.disable();
    } else {
      this.username.enable();
      this.password.enable();
    }
    this._loader = value;
  }

  private _loader = false;

  constructor(
    private _auth: AuthService,
    private _route: ActivatedRoute,
    private _router: Router,
  ) {}

  submit(): void {
    if (this.loader === true) {
      return;
    }

    this.error = '';
    if (this.form.invalid) {
      return;
    }

    this.loader = true;
    this._auth.register(
      this.form.value.username,
      this.form.value.password,
    ).subscribe(
      () => {
        this.loader = false;
        this.form.reset();
        this._router.navigate(['/']);
      },
      () => {
        this.loader = false;
        this.error = 'Внутренняя ошибка сервера. Попробуйте позже.';
      },
    );
  }
}
