import { h, FunctionComponent } from "preact";

interface InputCardPropTypes {
  value: string;
  placeholder: string;
  onSubmit: () => void;
}

export const InputCard: FunctionComponent<InputCardPropTypes> = (props) => (
  <div class="flex w-1/2 mx-auto items-center mt-5 shadow-lg border rounded-lg">
    <input
      type="text"
      name="todo"
      id="todo-input"
      placeholder={props.placeholder || "Enter new todo..."}
      class="focus:outline-none flex-1 p-3 text-md rounded-l-lg"
    />
    <button
      class="h-100 py-3 px-5 text-white font-semibold bg-green-400 hover:bg-green-500 rounded-r-lg transition-background duration-100"
      onClick={props.onSubmit}
    >
      {props.value || "Submit"}
    </button>
  </div>
);
