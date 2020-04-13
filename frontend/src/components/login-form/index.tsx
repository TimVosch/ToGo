import { h, FunctionalComponent, Component, ComponentChild } from "preact";

interface InputGroupProps {
  name?: string;
  label?: string;
  placeholder?: string;
  type?: string;
  class?: string;
}

const TextInput: FunctionalComponent<InputGroupProps> = ({
  name,
  type,
  label,
  placeholder,
  ...props
}) => {
  const classes = "flex flex-col " + props.class;
  return (
    <div class={classes}>
      <label class="mb-1 px-1 text-gray-500 uppercase font-semibold tracking-wide">
        {label || "Input:"}
      </label>
      <input
        type={type || "text"}
        name={name}
        placeholder={placeholder || "Enter here..."}
        class="border-b p-1"
      />
    </div>
  );
};

interface LoginFormProps {
  onLogin?: (email: string, password: string) => void;
}

export class LoginForm extends Component<LoginFormProps> {
  onSubmit(e: Event): void {
    e.preventDefault();

    const { onLogin } = this.props;
    if (typeof onLogin !== "function") {
      return;
    }

    const emailInput = document.getElementsByName(
      "login-input-email"
    )[0] as HTMLInputElement;
    const passwordInput = document.getElementsByName(
      "login-input-password"
    )[0] as HTMLInputElement;

    onLogin(emailInput.value, passwordInput.value);
  }

  render(): ComponentChild {
    return (
      <form
        onSubmitCapture={this.onSubmit.bind(this)}
        class="flex flex-col w-1/2 mx-auto mt-5 pb-3 shadow-md rounded-lg overflow-hidden"
      >
        <div class="mb-4 border-b bg-teal-500 font-semibold text-white p-3 text-center">
          Login
        </div>
        <div class="px-3">
          <TextInput
            name="login-input-email"
            label="Your email:"
            type="email"
            placeholder="example@email.com"
            class="mb-5"
          />
          <TextInput
            name="login-input-password"
            label="Your password:"
            placeholder=""
            type="password"
            class="mb-5"
          />
          <button
            type="submit"
            class="py-2 px-3 border border-teal-600 bg-teal-500 text-white rounded-md hover:bg-teal-400 transition-bg duration-100"
          >
            Login
          </button>
        </div>
      </form>
    );
  }
}
