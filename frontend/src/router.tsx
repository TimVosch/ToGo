import { h, Component, FunctionComponent, ComponentChild } from "preact";
import PreactRouter, { RouterOnChangeArgs, route } from "preact-router";
import { HomePage } from "./pages/home";
import { ProfilePage } from "./pages/profile";
import { LoginPage } from "./pages/login";
import { LoginService } from "./services/login/login.service";

type klass<T> = new () => T;
interface Route {
  path: string;
  component: klass<Component<any, any>> | FunctionComponent<any>;
  protected?: boolean;
}

const routes: Route[] = [
  {
    path: "/login",
    component: LoginPage,
  },
  {
    path: "/profile",
    component: ProfilePage,
  },
  {
    path: "/",
    component: HomePage,
    protected: true,
  },
];

export class Router extends Component {
  private loginService = new LoginService();

  /**
   * Checks whether the current user is allowed to
   * visit the given route.
   * @param e
   */
  checkAuthentication(e: RouterOnChangeArgs): void {
    // Only check protected if we're not logged in
    if (this.loginService.isLoggedIn()) {
      return;
    }
    // Get protected routes
    const protectedRoutes = routes
      .filter((r) => r.protected)
      .map((r) => r.path);
    // If route is protected, go back
    if (protectedRoutes.includes(e.url)) {
      alert("You need to be logged in");
      route(e.previous || "/login", true);
    }
  }

  // Called when the route is changed
  onRouteChange(e: RouterOnChangeArgs): void {
    this.checkAuthentication(e);
  }

  render(): ComponentChild {
    const routeItems = routes.map((r) => h(r.component, { path: r.path }, []));
    return (
      <PreactRouter onChange={this.onRouteChange.bind(this)}>
        {routeItems}
      </PreactRouter>
    );
  }
}
