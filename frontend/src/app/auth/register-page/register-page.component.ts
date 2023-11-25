import { Component } from '@angular/core';
import { FormControl, FormGroup, Validators } from '@angular/forms';
import { Router } from '@angular/router';
import { AuthService } from '../services/auth.service';

@Component({
  selector: 'app-register-page',
  templateUrl: './register-page.component.html',
  styleUrls: ['./register-page.component.less']
})
export class RegisterPageComponent {
  error: string = '';

  username = new FormControl('', [Validators.required]);
  password = new FormControl('', [Validators.required]);

  form: FormGroup = new FormGroup({
    username: this.username,
    password: this.password,
  });

  constructor(
    private _auth: AuthService,
    private _router: Router,
  ) {}

  submit(): void {
    this.error = '';
    if (this.form.invalid) {
      return;
    }

    this._auth.register(
      this.form.value.username,
      this.form.value.password,
    ).subscribe(
      () => {
        this.form.reset();
        this._router.navigate(['/auth/login']);
      },
      () => {
        this.error = 'Внутренняя ошибка сервера. Попробуйте позже.';
      },
    );
  }
}
