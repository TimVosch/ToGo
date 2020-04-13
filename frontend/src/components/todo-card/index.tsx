import { h } from "preact";
import { FunctionComponent } from "preact";
import { Todo } from "../../services/todo/todo.model";

interface TodoCardProps {
  todo: Todo;
}

export const TodoCard: FunctionComponent<TodoCardProps> = ({ todo }) => (
  <div
    key={todo.id}
    class="flex w-1/2 mx-auto mt-5 shadow-md rounded-lg border-t border-gray-400"
  >
    <div class="p-3 text-md">{todo.title}</div>
  </div>
);
