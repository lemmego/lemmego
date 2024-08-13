import React from "react";
import { useForm, usePage } from "@inertiajs/react";

const Login: React.FC = () => {
  const { errors, input, message } = usePage().props;
  const { data, setData, post, progress } = useForm({
    email: "",
    password: "",
    org_username: "",
  });

  function handleInput(e: React.ChangeEvent<HTMLInputElement>) {
    if (e.target.files && e.target.files.length > 0) {
      setData(e.target.name, e.target.files[0]);
      return;
    }

    setData(e.target.name, e.target.value);
  }

  function submit(e: React.FormEvent<HTMLFormElement>) {
    e.preventDefault();
    post("/login");
  }

  return (
    <div className="w-1/3 mx-auto">
      <h1 className="text-3xl text-center">Login</h1>
      {message && <p className="text-blue-500 text-center">{message}</p>}
      <form onSubmit={submit}>
        <div className="mt-2">
          <label htmlFor="org_username" className="label-primary">
            Org Username
          </label>
          <input
            id="org_username"
            name="org_username"
            type="text"
            className="input"
            value={data.org_username}
            onChange={handleInput}
          />
          {errors.org_username && (
            <p className="text-xs text-red-500">
              {errors.org_username.join(", ")}
            </p>
          )}
        </div>

        <div className="mt-2">
          <label htmlFor="email" className="label-primary">
            Email
          </label>
          <input
            id="email"
            name="email"
            type="email"
            className="input"
            value={data.email}
            onChange={handleInput}
          />
          {errors.email && (
            <p className="text-xs text-red-500">{errors.email.join(", ")}</p>
          )}
        </div>

        <div className="mt-2">
          <label htmlFor="password" className="label-primary">
            Password
          </label>
          <input
            id="password"
            name="password"
            type="password"
            className="input"
            value={data.password}
            onChange={handleInput}
          />
          {errors.password && (
            <p className="text-xs text-red-500">{errors.password.join(", ")}</p>
          )}
        </div>

        <div>
          <button type="submit" className="mt-4 btn-primary">
            Submit
          </button>
        </div>
      </form>
    </div>
  );
};

export default Login;
