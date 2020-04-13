import { Todo } from "./todo.model";
import { LoginService } from "../login/login.service";

/**
 *
 */
export class TodoService {
  private readonly login = new LoginService();
  private readonly URL: string;
  constructor(baseURL = "/api", resourceURI = "/todos") {
    this.URL = baseURL + resourceURI;
  }

  /**
   * Fetch all todos for the current user
   */
  async getTodos(): Promise<Todo[]> {
    const res = await fetch(this.URL, {
      headers: {
        Authorization: "Bearer " + this.login.getToken(),
      },
    });
    const body = await res.json();
    return body.data;
  }
}
