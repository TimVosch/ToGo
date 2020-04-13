import { AuthenticationError } from "./authentication.error";

export class LoginService {
  private static instance: LoginService | null = null;
  private readonly URL: string;
  private token: string | null = null;
  private payload: Record<string, any> | null = null;

  constructor(baseURL = "/api", resourceURI = "/auth") {
    if (LoginService.instance !== null) {
      return LoginService.instance;
    }
    this.URL = baseURL + resourceURI;
    LoginService.instance = this;
    this.recoverSession();
  }

  /**
   *
   */
  isLoggedIn(): boolean {
    return typeof this.token === "string";
  }

  /**
   *
   */
  getToken(): string {
    return this.token || "";
  }

  /**
   *
   */
  async login(email: string, password: string): Promise<void> {
    const requestBody = {
      email,
      password,
    };
    const res = await fetch(this.URL, {
      method: "POST",
      body: JSON.stringify(requestBody),
    });

    // Throw if response statusCode is not in the http OK range (200-299)
    if (res.status === 401 || res.status === 403) {
      throw new AuthenticationError("Login failed");
    } else if (res.status < 200 || res.status >= 300) {
      throw new AuthenticationError("Server error");
    }

    const body = (await res.json()) as Record<string, any>;

    // Ensure body has token property
    if (typeof body?.data?.token !== "string") {
      throw new AuthenticationError("Server response is malformed");
    }
    this.setSession(body.data.token);
  }

  /**
   *
   */
  logout(): void {
    this.setSession(null);
  }

  /**
   *
   * @param token
   */
  private setSession(token: string | null): void {
    // Either set the token, which means we logged in
    // or unset it which means we signed out.
    if (typeof token === "string") {
      // Sign in
      const parts = token.split(".");
      this.payload = JSON.parse(atob(parts[1])) as Record<string, any>;

      // Check if token is expired
      if (this.payload.exp && this.payload.exp <= Date.now()) {
        throw new AuthenticationError("Session is expired");
      }

      localStorage.setItem("token", token);
    } else {
      // Sign out
      localStorage.removeItem("token");
    }
    this.token = token;
  }

  /**
   *
   */
  private recoverSession(): void {
    const token = localStorage.getItem("token");
    if (token === null) {
      return;
    }
    this.setSession(token);
  }
}
