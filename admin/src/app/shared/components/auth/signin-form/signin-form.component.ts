import { Component } from "@angular/core";
import { LabelComponent } from "../../form/label/label.component";
import { CheckboxComponent } from "../../form/input/checkbox.component";
import { ButtonComponent } from "../../ui/button/button.component";
import { InputFieldComponent } from "../../form/input/input-field.component";
import { RouterModule, Router } from "@angular/router";
import { FormsModule } from "@angular/forms";
import { AuthService } from "../../../../core/services/auth.service";
import { AuthStateService } from "../../../../core/services/auth-state.service";
import { LoginRequest } from "../../../models/auth.model";

@Component({
  selector: "app-signin-form",
  imports: [
    LabelComponent,
    CheckboxComponent,
    ButtonComponent,
    InputFieldComponent,
    RouterModule,
    FormsModule,
  ],
  templateUrl: "./signin-form.component.html",
  styles: ``,
})
export class SigninFormComponent {
  showPassword = false;
  isChecked = false;

  email = "";
  password = "";

  isLoading = false;
  error: string | null = null;

  constructor(
    private auth: AuthService,
    private authState: AuthStateService,
    private router: Router,
  ) {}

  togglePasswordVisibility() {
    this.showPassword = !this.showPassword;
  }

  onSignIn() {
    if (!this.email || !this.password) {
      this.error = "Email and password are required";
      return;
    }

    this.error = null;
    this.isLoading = true;

    const payload: LoginRequest = {
      email: this.email,
      password: this.password,
    };

    this.auth.login(payload).subscribe({
      next: () => {
        this.isLoading = false;
        this.authState.setAuthenticated(true);
        this.router.navigate(["/"]);
      },
      error: (err) => {
        console.error("Login failed", err);
        this.isLoading = false;
        this.error =
          err?.error?.message ||
          err?.message ||
          "Login failed. Please try again.";
      },
    });
  }
}
