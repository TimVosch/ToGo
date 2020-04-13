import { h, FunctionalComponent } from "preact";
import logo from "./logo.svg";

export const Navigation: FunctionalComponent = () => (
  <nav class="flex items-stretch bg-teal-500 justify-between">
    {/* Branding */}
    <div class="flex items-center font-semibold text-xl p-1 tracking-tight text-white">
      <img src={logo} alt="Logo" class="w-10 h-10 mr-2" />
      ToGO Planning
    </div>
    {/* Nav items */}
    <div class="flex text-white">
      <a class="flex items-center pl-2 pr-3 hover:bg-teal-600" href="/">
        <span>Home</span>
      </a>
      <a class="flex items-center pl-2 pr-3 hover:bg-teal-600" href="/login">
        <span>Login</span>
      </a>
      <a
        class="flex items-center pl-2 pr-3 hover:bg-teal-600"
        href="https://github.com/timvosch/togo"
      >
        <span>Github</span>
      </a>
    </div>
  </nav>
);
