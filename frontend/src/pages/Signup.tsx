import { A, useNavigate } from "@solidjs/router";
import { createSignal, type Component } from "solid-js";
import { createUser } from "../api/user";
import { APIError } from "../api";
import toast from "solid-toast";
import { setSessionCookie } from "../lib/cookie";

export default function SignUpPage() { 
  const[name, setName] = createSignal("");
  const [email, setEmail] = createSignal("");
  const [password, setPassword] = createSignal("");

  const navigate = useNavigate();
  const [isLoading, setIsLoading] = createSignal(false);

  const handleSubmit = async (e: Event) => {
    e.preventDefault();
    setIsLoading(true);
    try {
      const { jwt_token: token } = await createUser({
        email: email(),
        password: password(),
        name: name(),
      });
      
      setSessionCookie(token);

      toast.success("Successfully created user " + name() + "!");
      
      navigate("/");
    } catch (error) {
      if (error instanceof APIError) {
        toast.error(error.message);
      } else {
        toast.error("An unknown error occurred");
      }
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div class="flex justify-start items-center flex-col gap-y-4 pt-20 h-screen">
      <A class="btn" href="/">
        home
      </A>
      Sign Up
      <form
        onSubmit={handleSubmit}
        class="flex flex-col gap-y-3 border rounded-md p-4 shadow-sm"
      >
        <input
          type="text"
          placeholder="Name"
          onInput={(e) => setName(e.currentTarget.value)}
          class="border-b focus:outline-none"
          min={2}
          required
        />
        <input
          type="email"
          placeholder="Email"
          onInput={(e) => setEmail(e.currentTarget.value)}
          class="border-b focus:outline-none"
          required
        />
        <input
          type="password"
          placeholder="Password"
          onInput={(e) => setPassword(e.currentTarget.value)}
          class="border-b focus:outline-none"
          required
        />
        <button class="btn" type="submit" disabled={isLoading()}>
          {isLoading() ? "Loading..." : "Sign Up"}
        </button>
      </form>
      <A href="/signin" class="btn">
        Already have an account?{" "}
      </A>
    </div>
  );
};