import { h, Component, ComponentChild } from "preact";
import { Navigation } from "../components/navigation";
import { InputCard } from "../components/input-card";
import { TodoService } from "../services/todo/todo.service";
import { TodoCard } from "../components/todo-card";
import { Todo } from "../services/todo/todo.model";

interface HomePageState {
  todos: Todo[];
}

export class HomePage extends Component<{}, HomePageState> {
  private readonly todos = new TodoService();

  constructor() {
    super();
    this.state = {
      todos: [],
    };
  }

  componentDidMount(): void {
    this.fetchTodos();
  }

  addTodo(): void {
    // this.todos.
  }

  async fetchTodos(): Promise<void> {
    const todos = await this.todos.getTodos();
    this.setState({
      todos,
    });
  }

  render(): ComponentChild {
    const { todos } = this.state;
    const todoCards = todos.map((todo) => (
      <TodoCard key={todo.id} todo={todo} />
    ));
    return (
      <div>
        <Navigation />
        <InputCard
          placeholder="Add new todo..."
          value="Add"
          onSubmit={this.addTodo.bind(this)}
        />
        {todoCards}
      </div>
    );
  }
}
