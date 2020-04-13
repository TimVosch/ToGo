import { h, Component, ComponentChild, Fragment } from "preact";
import { route } from "preact-router";
import { Navigation } from "../components/navigation";
import { LoginForm } from "../components/login-form";
import { LoginService } from "../services/login/login.service";
import { AuthenticationError } from "../services/login/authentication.error";

export class LoginPage extends Component {
  private readonly loginService = new LoginService();

  componentDidMount(): void {
    if (this.loginService.isLoggedIn()) {
      // TODO: Modal
      alert("Already logged in!");
      route("/");
    }
  }

  async login(email: string, password: string): Promise<void> {
    try {
      // Login and redirect
      await this.loginService.login(email, password);
      route("/");
    } catch (e) {
      if (e instanceof AuthenticationError) {
        alert(e);
      }
      // TODO: we dont know what happened.
    }
  }

  render(): ComponentChild {
    return (
      <Fragment>
        <Navigation />
        <LoginForm onLogin={this.login.bind(this)} />
      </Fragment>
    );
  }
}
