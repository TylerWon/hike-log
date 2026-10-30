import "../../assets/styles/button.css";

import { useMutation } from "@tanstack/react-query";
import { useState } from "react";
import { type NavigateFunction, useNavigate } from "react-router";

import { login } from "../../network/api/auth/auth";
import { AuthError } from "../../network/api/auth/errors";
import TextField from "../../components/form/TextField";
import GenericError from "../../components/GenericError";
import PendingButton from "../../components/PendingButton";
import { type LoginFormData, LoginFormDataSchema } from "../../schemas/forms/login";

export interface LoginProps {
  // Dependency injection for testing purposes
  navigateFn?: NavigateFunction;
}

export default function Login({ navigateFn }: LoginProps) {
  let navigate = useNavigate();
  navigate = navigateFn ?? navigate;

  const [password, setPassword] = useState<string>("");
  const [passwordError, setPasswordError] = useState<string>("");

  const onSubmit = async (formData: LoginFormData) => {
    try {
      await login(formData);
    } catch (e) {
      // Propogate non-auth failures to cause mutation to fail
      if (!(e instanceof AuthError)) {
        throw e;
      }

      // Display error on the password field when auth fails
      setPasswordError("Incorrect password");
      return;
    }

    navigate("/");
  };

  const loginMutation = useMutation({
    mutationFn: onSubmit, // called when mutate() is invoked for this mutation
    onError: (error) => console.error("Failed to login: ", error),
  });

  const loginPending = loginMutation.isPending;
  const loginError = loginMutation.isError;

  const validate = (): LoginFormData | null => {
    const login = { password };

    const result = LoginFormDataSchema.safeParse(login);
    if (!result.success) {
      setPasswordError(result.error.issues[0]?.message ?? "");
      return null;
    }

    return result.data;
  };

  const handleSubmit = (e: React.SubmitEvent) => {
    e.preventDefault();

    const formData = validate();
    if (!formData) {
      return;
    }

    setPasswordError("");
    loginMutation.mutate(formData);
  };

  return (
    <main className="relative flex min-h-screen items-center justify-center overflow-hidden bg-forest-950 px-6 py-16 text-cream-100">
      <div className="w-full flex flex-col gap-8 max-w-sm rounded-lg border border-forest-800 bg-forest-900 p-7">
        <div>
          <svg
            aria-hidden="true"
            className="mb-5 text-sage-500"
            fill="none"
            height="40"
            stroke="currentColor"
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth="1.25"
            viewBox="0 0 40 40"
            width="40"
          >
            <path d="M4 31 15 13l6 9 5-7 10 16H4Z" />
            <path d="m11 31 8-12" />
          </svg>
          <p className="mb-2 font-mono text-[10px] uppercase tracking-[0.2em] text-sage-500">Hike Log</p>
          <h1 className="font-serif text-3xl font-semibold tracking-tight">Admin login</h1>
        </div>

        <form className="flex flex-col gap-5" onSubmit={handleSubmit}>
          <TextField
            error={passwordError}
            id="password"
            inputProps={{
              onChange: (e) => setPassword(e.target.value),
              placeholder: "Enter password",
              type: "password",
              value: password,
            }}
            label="Password"
          />

          <PendingButton
            buttonProps={{
              className: "primary-button-default w-full px-4 py-2.5",
              disabled: loginPending,
              type: "submit",
            }}
            pending={loginPending}
            pendingText="Logging in..."
            text="Login"
          />
        </form>

        {loginError && <GenericError />}
      </div>
    </main>
  );
}
