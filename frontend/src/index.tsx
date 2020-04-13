import { h, render } from "preact";
import "./index.css";
import { Router } from "./router";
import { LoginService } from "./services/login/login.service";
import { TodoService } from "./services/todo/todo.service";

/**
 * Initializes services
 */
function bootstrap(): void {
  new LoginService();
  new TodoService();
}

const root = document.querySelector("#container");
if (!root) {
  // eslint-disable-next-line
  console.error("Missing #container");
} else {
  bootstrap();
  render(<Router />, root);
}

if (module && module.hot) {
  module.hot.accept();

  module.hot.addStatusHandler((status) => {
    // eslint-disable-next-line
    if (status === "prepare") console.clear();
  });
}
