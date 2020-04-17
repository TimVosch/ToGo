import { h, Component, ComponentChild } from "preact";

interface InputCardPropTypes {
  value: string;
  placeholder: string;
  onSubmit: (title: string) => void | Promise<void>;
}

export class InputCard extends Component<InputCardPropTypes> {
  onSubmit(e: Event): void {
    e.preventDefault();

    const { onSubmit } = this.props;
    if (typeof onSubmit !== "function") {
      return;
    }

    const todoInput = document.getElementsByName(
      "todo-input-title"
    )[0] as HTMLInputElement;
    onSubmit(todoInput.value);
  }

  render(): ComponentChild {
    const { placeholder, value } = this.props;
    return (
      <form
        onSubmitCapture={this.onSubmit?.bind(this)}
        class="flex w-1/2 mx-auto items-center mt-5 shadow-lg border rounded-lg"
      >
        <input
          type="text"
          name="todo-input-title"
          id="todo-input-title"
          placeholder={placeholder || "Enter new todo..."}
          class="focus:outline-none flex-1 p-3 text-md rounded-l-lg"
        />
        <button
          class="h-100 py-3 px-5 text-white font-semibold bg-green-400 hover:bg-green-500 rounded-r-lg transition-background duration-100"
          type="submit"
        >
          {value || "Submit"}
        </button>
      </form>
    );
  }
}
