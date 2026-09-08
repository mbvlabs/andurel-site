import { jsxs, jsx, Fragment } from "react/jsx-runtime";
import { useForm, Link, router, Head, createInertiaApp } from "@inertiajs/react";
import { XIcon, SearchIcon, CheckIcon, AlignLeftIcon, PanelLeftIcon, ChevronDownIcon } from "lucide-react";
import * as React from "react";
import { useState, useEffect } from "react";
import { cn } from "cn";
import { Button as Button$1 } from "@base-ui/react/button";
import { cva } from "class-variance-authority";
import { Command as Command$1 } from "cmdk";
import { Dialog as Dialog$1 } from "@base-ui/react/dialog";
import { Menu } from "@base-ui/react/menu";
import { mergeProps } from "@base-ui/react/merge-props";
import { useRender } from "@base-ui/react/use-render";
import { Tooltip as Tooltip$1 } from "@base-ui/react/tooltip";
import createServer from "@inertiajs/react/server";
import ReactDOMServer from "react-dom/server";
const routes = {
  confirmationCreate: () => "/users/confirmation",
  confirmationNew: () => "/users/confirmation/new",
  documentationShow: (version, slug) => `/docs/${version}/${slug}`,
  homePage: () => "/",
  passwordCreate: () => "/users/password",
  passwordEdit: (token) => `/users/password/${token}/edit`,
  passwordNew: () => "/users/password/new",
  passwordUpdate: () => "/users/password",
  registrationCreate: () => "/users",
  registrationNew: () => "/users/sign-up",
  sessionCreate: () => "/users/sign-in",
  sessionDestroy: () => "/users/sign-out",
  sessionNew: () => "/users/sign-in"
};
function Layout({ children }) {
  return /* @__PURE__ */ jsxs("main", { className: "relative flex min-h-screen flex-col overflow-hidden bg-[#090b0d] text-[#e4dfd2]", children: [
    /* @__PURE__ */ jsx(
      "div",
      {
        className: "pointer-events-none absolute inset-0 opacity-60",
        style: {
          backgroundImage: "radial-gradient(circle at 12% 18%, #f2ead8 0 1px, transparent 1.5px), radial-gradient(circle at 82% 22%, #aaa393 0 1px, transparent 1.5px), radial-gradient(circle at 67% 72%, #f2ead8 0 1px, transparent 1.5px), radial-gradient(circle at 24% 83%, #8f8a7d 0 1px, transparent 1.5px)"
        }
      }
    ),
    /* @__PURE__ */ jsx("header", { className: "relative", children: /* @__PURE__ */ jsxs("div", { className: "mx-auto flex w-full max-w-[960px] items-center justify-between px-4 py-3", children: [
      /* @__PURE__ */ jsxs("a", { className: "inline-flex items-center gap-3 text-sm font-semibold text-[#f2ead8]", href: routes.homePage(), children: [
        /* @__PURE__ */ jsxs("span", { className: "grid size-8 grid-cols-2 gap-1 border border-[#52605c] bg-[#101414] p-1 shadow-sm shadow-black/40", children: [
          /* @__PURE__ */ jsx("span", { className: "border border-[#8df7a4]" }),
          /* @__PURE__ */ jsx("span", { className: "border border-[#52605c]" }),
          /* @__PURE__ */ jsx("span", { className: "border border-[#52605c]" }),
          /* @__PURE__ */ jsx("span", { className: "bg-[#8df7a4]" })
        ] }),
        /* @__PURE__ */ jsx("span", { children: "Andurel." })
      ] }),
      /* @__PURE__ */ jsxs("nav", { className: "flex flex-wrap items-center justify-end gap-3 text-sm", children: [
        /* @__PURE__ */ jsx("a", { className: "px-2 py-1 text-[#aaa393] transition hover:text-[#f2ead8]", href: "https://andurel.com", children: "Documentation" }),
        /* @__PURE__ */ jsx("a", { className: "px-2 py-1 text-[#aaa393] transition hover:text-[#f2ead8]", href: routes.sessionNew(), children: "Log in" }),
        /* @__PURE__ */ jsx("a", { className: "px-2 py-1 text-[#aaa393] transition hover:text-[#f2ead8]", href: routes.registrationNew(), children: "Register" })
      ] })
    ] }) }),
    /* @__PURE__ */ jsx("div", { className: "relative flex flex-1 items-center justify-center px-6 py-6", children }),
    /* @__PURE__ */ jsxs("footer", { className: "relative py-3 text-center text-sm text-[#8f8a7d]", children: [
      "© ",
      (/* @__PURE__ */ new Date()).getFullYear(),
      " andurel."
    ] })
  ] });
}
function ConfirmEmail({ errors = {} }) {
  const form = useForm({ code: "" });
  function submit(event) {
    event.preventDefault();
    form.post(routes.confirmationCreate());
  }
  return /* @__PURE__ */ jsx(Layout, { children: /* @__PURE__ */ jsxs("section", { className: "w-full max-w-md border border-[#2f3a37] bg-[#101414]/90 shadow-sm shadow-black/40", children: [
    /* @__PURE__ */ jsxs("div", { className: "p-6 pb-0", children: [
      /* @__PURE__ */ jsx("h1", { className: "text-xl font-semibold text-[#f2ead8]", children: "Verify Your Email" }),
      /* @__PURE__ */ jsx("p", { className: "mt-1 text-sm text-[#8f8a7d]", children: "Please enter the 6-digit verification code sent to your email." })
    ] }),
    /* @__PURE__ */ jsx("div", { className: "p-6", children: /* @__PURE__ */ jsxs("form", { onSubmit: (event) => submit(event.nativeEvent), className: "space-y-5", children: [
      /* @__PURE__ */ jsxs("div", { className: "space-y-1", children: [
        /* @__PURE__ */ jsx("label", { className: "text-sm font-medium text-[#c7c0ad]", htmlFor: "code", children: "Verification Code" }),
        /* @__PURE__ */ jsx("input", { id: "code", type: "text", maxLength: 6, value: form.data.code, onChange: (event) => form.setData("code", event.target.value), className: "flex h-9 w-full border border-[#2f3a37] bg-[#090c0d] px-3 py-1 text-center text-sm tracking-[0.3em] text-[#e4dfd2] shadow-inner shadow-black/35 focus:border-[#8df7a4] focus:outline-none focus:ring-2 focus:ring-[#8df7a4]/20", required: true }),
        errors.code && /* @__PURE__ */ jsx("p", { className: "text-sm font-medium text-[#ff875f]", children: errors.code })
      ] }),
      /* @__PURE__ */ jsx("button", { type: "submit", disabled: form.processing, className: "inline-flex w-full items-center justify-center bg-[#ff6b1a] px-4 py-2 text-sm font-medium text-[#130f0b] shadow-sm shadow-black/40 hover:bg-[#ff8748] disabled:opacity-60", children: form.processing ? "Loading" : "Verify Email" })
    ] }) })
  ] }) });
}
const __vite_glob_0_0 = /* @__PURE__ */ Object.freeze(/* @__PURE__ */ Object.defineProperty({
  __proto__: null,
  default: ConfirmEmail
}, Symbol.toStringTag, { value: "Module" }));
function Login({ errors = {} }) {
  const form = useForm({ email: "", password: "" });
  function submit(event) {
    event.preventDefault();
    form.post(routes.sessionCreate());
  }
  return /* @__PURE__ */ jsx(Layout, { children: /* @__PURE__ */ jsxs("section", { className: "w-full max-w-md border border-[#2f3a37] bg-[#101414]/90 shadow-sm shadow-black/40", children: [
    /* @__PURE__ */ jsxs("div", { className: "p-6 pb-0", children: [
      /* @__PURE__ */ jsx("h1", { className: "text-xl font-semibold text-[#f2ead8]", children: "Login to your account" }),
      /* @__PURE__ */ jsx("p", { className: "mt-1 text-sm text-[#8f8a7d]", children: "Enter your details below to login to your account" })
    ] }),
    /* @__PURE__ */ jsxs("div", { className: "p-6", children: [
      /* @__PURE__ */ jsxs("form", { onSubmit: (event) => submit(event.nativeEvent), className: "space-y-5", children: [
        /* @__PURE__ */ jsxs("div", { className: "space-y-1", children: [
          /* @__PURE__ */ jsx("label", { className: "text-sm font-medium text-[#c7c0ad]", htmlFor: "email", children: "Email" }),
          /* @__PURE__ */ jsx("input", { id: "email", type: "email", value: form.data.email, onChange: (event) => form.setData("email", event.target.value), className: "flex h-9 w-full border border-[#2f3a37] bg-[#090c0d] px-3 py-1 text-sm text-[#e4dfd2] shadow-inner shadow-black/35 focus:border-[#8df7a4] focus:outline-none focus:ring-2 focus:ring-[#8df7a4]/20", required: true }),
          errors.email && /* @__PURE__ */ jsx("p", { className: "text-sm font-medium text-[#ff875f]", children: errors.email })
        ] }),
        /* @__PURE__ */ jsxs("div", { className: "space-y-1", children: [
          /* @__PURE__ */ jsx("label", { className: "text-sm font-medium text-[#c7c0ad]", htmlFor: "password", children: "Password" }),
          /* @__PURE__ */ jsx("input", { id: "password", type: "password", value: form.data.password, onChange: (event) => form.setData("password", event.target.value), className: "flex h-9 w-full border border-[#2f3a37] bg-[#090c0d] px-3 py-1 text-sm text-[#e4dfd2] shadow-inner shadow-black/35 focus:border-[#8df7a4] focus:outline-none focus:ring-2 focus:ring-[#8df7a4]/20", required: true }),
          errors.password && /* @__PURE__ */ jsx("p", { className: "text-sm font-medium text-[#ff875f]", children: errors.password })
        ] }),
        /* @__PURE__ */ jsx("p", { className: "text-right text-sm", children: /* @__PURE__ */ jsx(Link, { className: "text-[#d7d0bf] hover:text-[#f2ead8] hover:underline", href: routes.passwordNew(), children: "Forgot your password?" }) }),
        /* @__PURE__ */ jsx("button", { type: "submit", disabled: form.processing, className: "inline-flex w-full items-center justify-center bg-[#ff6b1a] px-4 py-2 text-sm font-medium text-[#130f0b] shadow-sm shadow-black/40 hover:bg-[#ff8748] disabled:opacity-60", children: form.processing ? "Loading" : "Login" })
      ] }),
      /* @__PURE__ */ jsxs("p", { className: "mt-6 text-center text-sm text-[#8f8a7d]", children: [
        "Don't have an account? ",
        /* @__PURE__ */ jsx(Link, { className: "text-[#d7d0bf] hover:text-[#f2ead8] hover:underline", href: routes.registrationNew(), children: "Sign up" })
      ] })
    ] })
  ] }) });
}
const __vite_glob_0_1 = /* @__PURE__ */ Object.freeze(/* @__PURE__ */ Object.defineProperty({
  __proto__: null,
  default: Login
}, Symbol.toStringTag, { value: "Module" }));
function Registration({ errors = {} }) {
  const form = useForm({ email: "", password: "", confirmPassword: "" });
  function submit(event) {
    event.preventDefault();
    form.post(routes.registrationCreate());
  }
  return /* @__PURE__ */ jsx(Layout, { children: /* @__PURE__ */ jsxs("section", { className: "w-full max-w-md border border-[#2f3a37] bg-[#101414]/90 shadow-sm shadow-black/40", children: [
    /* @__PURE__ */ jsxs("div", { className: "p-6 pb-0", children: [
      /* @__PURE__ */ jsx("h1", { className: "text-xl font-semibold text-[#f2ead8]", children: "Create an account" }),
      /* @__PURE__ */ jsx("p", { className: "mt-1 text-sm text-[#8f8a7d]", children: "Enter your details below to create your account" })
    ] }),
    /* @__PURE__ */ jsxs("div", { className: "p-6", children: [
      /* @__PURE__ */ jsxs("form", { onSubmit: (event) => submit(event.nativeEvent), className: "space-y-5", children: [
        /* @__PURE__ */ jsxs("div", { className: "space-y-1", children: [
          /* @__PURE__ */ jsx("label", { className: "text-sm font-medium text-[#c7c0ad]", htmlFor: "email", children: "Email" }),
          /* @__PURE__ */ jsx("input", { id: "email", type: "email", value: form.data.email, onChange: (event) => form.setData("email", event.target.value), className: "flex h-9 w-full border border-[#2f3a37] bg-[#090c0d] px-3 py-1 text-sm text-[#e4dfd2] shadow-inner shadow-black/35 focus:border-[#8df7a4] focus:outline-none focus:ring-2 focus:ring-[#8df7a4]/20", required: true }),
          errors.email && /* @__PURE__ */ jsx("p", { className: "text-sm font-medium text-[#ff875f]", children: errors.email })
        ] }),
        /* @__PURE__ */ jsxs("div", { className: "space-y-1", children: [
          /* @__PURE__ */ jsx("label", { className: "text-sm font-medium text-[#c7c0ad]", htmlFor: "password", children: "Password" }),
          /* @__PURE__ */ jsx("input", { id: "password", type: "password", value: form.data.password, onChange: (event) => form.setData("password", event.target.value), className: "flex h-9 w-full border border-[#2f3a37] bg-[#090c0d] px-3 py-1 text-sm text-[#e4dfd2] shadow-inner shadow-black/35 focus:border-[#8df7a4] focus:outline-none focus:ring-2 focus:ring-[#8df7a4]/20", required: true }),
          errors.password && /* @__PURE__ */ jsx("p", { className: "text-sm font-medium text-[#ff875f]", children: errors.password })
        ] }),
        /* @__PURE__ */ jsxs("div", { className: "space-y-1", children: [
          /* @__PURE__ */ jsx("label", { className: "text-sm font-medium text-[#c7c0ad]", htmlFor: "confirmPassword", children: "Confirm Password" }),
          /* @__PURE__ */ jsx("input", { id: "confirmPassword", type: "password", value: form.data.confirmPassword, onChange: (event) => form.setData("confirmPassword", event.target.value), className: "flex h-9 w-full border border-[#2f3a37] bg-[#090c0d] px-3 py-1 text-sm text-[#e4dfd2] shadow-inner shadow-black/35 focus:border-[#8df7a4] focus:outline-none focus:ring-2 focus:ring-[#8df7a4]/20", required: true }),
          errors.confirmPassword && /* @__PURE__ */ jsx("p", { className: "text-sm font-medium text-[#ff875f]", children: errors.confirmPassword })
        ] }),
        /* @__PURE__ */ jsx("button", { type: "submit", disabled: form.processing, className: "inline-flex w-full items-center justify-center bg-[#ff6b1a] px-4 py-2 text-sm font-medium text-[#130f0b] shadow-sm shadow-black/40 hover:bg-[#ff8748] disabled:opacity-60", children: form.processing ? "Loading" : "Sign Up" })
      ] }),
      /* @__PURE__ */ jsxs("p", { className: "mt-6 text-center text-sm text-[#8f8a7d]", children: [
        "Already have an account? ",
        /* @__PURE__ */ jsx(Link, { className: "text-[#d7d0bf] hover:text-[#f2ead8] hover:underline", href: routes.sessionNew(), children: "Login" })
      ] })
    ] })
  ] }) });
}
const __vite_glob_0_2 = /* @__PURE__ */ Object.freeze(/* @__PURE__ */ Object.defineProperty({
  __proto__: null,
  default: Registration
}, Symbol.toStringTag, { value: "Module" }));
function ResetPassword({ token, errors = {} }) {
  const form = useForm({ resetPasswordToken: token, password: "", confirmPassword: "" });
  function submit(event) {
    event.preventDefault();
    form.put(routes.passwordUpdate());
  }
  return /* @__PURE__ */ jsx(Layout, { children: /* @__PURE__ */ jsxs("section", { className: "w-full max-w-md border border-[#2f3a37] bg-[#101414]/90 shadow-sm shadow-black/40", children: [
    /* @__PURE__ */ jsxs("div", { className: "p-6 pb-0", children: [
      /* @__PURE__ */ jsx("h1", { className: "text-xl font-semibold text-[#f2ead8]", children: "Reset Your Password" }),
      /* @__PURE__ */ jsx("p", { className: "mt-1 text-sm text-[#8f8a7d]", children: "Enter your new password below." })
    ] }),
    /* @__PURE__ */ jsx("div", { className: "p-6", children: /* @__PURE__ */ jsxs("form", { onSubmit: (event) => submit(event.nativeEvent), className: "space-y-5", children: [
      errors.resetPasswordToken && /* @__PURE__ */ jsx("p", { className: "text-sm font-medium text-[#ff875f]", children: errors.resetPasswordToken }),
      /* @__PURE__ */ jsxs("div", { className: "space-y-1", children: [
        /* @__PURE__ */ jsx("label", { className: "text-sm font-medium text-[#c7c0ad]", htmlFor: "password", children: "New Password" }),
        /* @__PURE__ */ jsx("input", { id: "password", type: "password", value: form.data.password, onChange: (event) => form.setData("password", event.target.value), className: "flex h-9 w-full border border-[#2f3a37] bg-[#090c0d] px-3 py-1 text-sm text-[#e4dfd2] shadow-inner shadow-black/35 focus:border-[#8df7a4] focus:outline-none focus:ring-2 focus:ring-[#8df7a4]/20", required: true }),
        errors.password && /* @__PURE__ */ jsx("p", { className: "text-sm font-medium text-[#ff875f]", children: errors.password })
      ] }),
      /* @__PURE__ */ jsxs("div", { className: "space-y-1", children: [
        /* @__PURE__ */ jsx("label", { className: "text-sm font-medium text-[#c7c0ad]", htmlFor: "confirmPassword", children: "Confirm New Password" }),
        /* @__PURE__ */ jsx("input", { id: "confirmPassword", type: "password", value: form.data.confirmPassword, onChange: (event) => form.setData("confirmPassword", event.target.value), className: "flex h-9 w-full border border-[#2f3a37] bg-[#090c0d] px-3 py-1 text-sm text-[#e4dfd2] shadow-inner shadow-black/35 focus:border-[#8df7a4] focus:outline-none focus:ring-2 focus:ring-[#8df7a4]/20", required: true }),
        errors.confirmPassword && /* @__PURE__ */ jsx("p", { className: "text-sm font-medium text-[#ff875f]", children: errors.confirmPassword })
      ] }),
      /* @__PURE__ */ jsx("button", { type: "submit", disabled: form.processing, className: "inline-flex w-full items-center justify-center bg-[#ff6b1a] px-4 py-2 text-sm font-medium text-[#130f0b] shadow-sm shadow-black/40 hover:bg-[#ff8748] disabled:opacity-60", children: form.processing ? "Loading" : "Reset Password" })
    ] }) })
  ] }) });
}
const __vite_glob_0_3 = /* @__PURE__ */ Object.freeze(/* @__PURE__ */ Object.defineProperty({
  __proto__: null,
  default: ResetPassword
}, Symbol.toStringTag, { value: "Module" }));
function ResetPasswordRequest({ errors = {} }) {
  const form = useForm({ email: "" });
  function submit(event) {
    event.preventDefault();
    form.post(routes.passwordCreate());
  }
  return /* @__PURE__ */ jsx(Layout, { children: /* @__PURE__ */ jsxs("section", { className: "w-full max-w-md border border-[#2f3a37] bg-[#101414]/90 shadow-sm shadow-black/40", children: [
    /* @__PURE__ */ jsxs("div", { className: "p-6 pb-0", children: [
      /* @__PURE__ */ jsx("h1", { className: "text-xl font-semibold text-[#f2ead8]", children: "Reset Password" }),
      /* @__PURE__ */ jsx("p", { className: "mt-1 text-sm text-[#8f8a7d]", children: "Enter your email address and we'll send you a code to reset your password." })
    ] }),
    /* @__PURE__ */ jsxs("div", { className: "p-6", children: [
      /* @__PURE__ */ jsxs("form", { onSubmit: (event) => submit(event.nativeEvent), className: "space-y-5", children: [
        /* @__PURE__ */ jsxs("div", { className: "space-y-1", children: [
          /* @__PURE__ */ jsx("label", { className: "text-sm font-medium text-[#c7c0ad]", htmlFor: "email", children: "Email" }),
          /* @__PURE__ */ jsx("input", { id: "email", type: "email", value: form.data.email, onChange: (event) => form.setData("email", event.target.value), className: "flex h-9 w-full border border-[#2f3a37] bg-[#090c0d] px-3 py-1 text-sm text-[#e4dfd2] shadow-inner shadow-black/35 focus:border-[#8df7a4] focus:outline-none focus:ring-2 focus:ring-[#8df7a4]/20", required: true }),
          errors.email && /* @__PURE__ */ jsx("p", { className: "text-sm font-medium text-[#ff875f]", children: errors.email })
        ] }),
        /* @__PURE__ */ jsx("button", { type: "submit", disabled: form.processing, className: "inline-flex w-full items-center justify-center bg-[#ff6b1a] px-4 py-2 text-sm font-medium text-[#130f0b] shadow-sm shadow-black/40 hover:bg-[#ff8748] disabled:opacity-60", children: form.processing ? "Loading" : "Send Reset Code" })
      ] }),
      /* @__PURE__ */ jsxs("p", { className: "mt-6 text-center text-sm text-[#8f8a7d]", children: [
        "Remember your password? ",
        /* @__PURE__ */ jsx(Link, { className: "text-[#d7d0bf] hover:text-[#f2ead8] hover:underline", href: routes.sessionNew(), children: "Login" })
      ] })
    ] })
  ] }) });
}
const __vite_glob_0_4 = /* @__PURE__ */ Object.freeze(/* @__PURE__ */ Object.defineProperty({
  __proto__: null,
  default: ResetPasswordRequest
}, Symbol.toStringTag, { value: "Module" }));
const buttonVariants = cva(
  "group/button inline-flex shrink-0 items-center justify-center rounded-none border border-transparent bg-clip-padding text-xs font-medium whitespace-nowrap transition-all outline-none select-none focus-visible:border-ring focus-visible:ring-1 focus-visible:ring-ring/50 active:not-aria-[haspopup]:translate-y-px disabled:pointer-events-none disabled:opacity-50 aria-invalid:border-destructive aria-invalid:ring-1 aria-invalid:ring-destructive/20 dark:aria-invalid:border-destructive/50 dark:aria-invalid:ring-destructive/40 [&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*='size-'])]:size-4",
  {
    variants: {
      variant: {
        default: "bg-primary text-primary-foreground hover:bg-primary/80",
        outline: "border-border bg-background hover:bg-muted hover:text-foreground aria-expanded:bg-muted aria-expanded:text-foreground dark:border-input dark:bg-input/30 dark:hover:bg-input/50",
        secondary: "bg-secondary text-secondary-foreground hover:bg-[color-mix(in_oklch,var(--secondary),var(--foreground)_5%)] aria-expanded:bg-secondary aria-expanded:text-secondary-foreground",
        ghost: "hover:bg-muted hover:text-foreground aria-expanded:bg-muted aria-expanded:text-foreground dark:hover:bg-muted/50",
        destructive: "bg-destructive/10 text-destructive hover:bg-destructive/20 focus-visible:border-destructive/40 focus-visible:ring-destructive/20 dark:bg-destructive/20 dark:hover:bg-destructive/30 dark:focus-visible:ring-destructive/40",
        link: "text-primary underline-offset-4 hover:underline"
      },
      size: {
        default: "h-8 gap-1.5 px-2.5 has-data-[icon=inline-end]:pr-2 has-data-[icon=inline-start]:pl-2",
        xs: "h-6 gap-1 rounded-none px-2 text-xs has-data-[icon=inline-end]:pr-1.5 has-data-[icon=inline-start]:pl-1.5 [&_svg:not([class*='size-'])]:size-3",
        sm: "h-7 gap-1 rounded-none px-2.5 has-data-[icon=inline-end]:pr-1.5 has-data-[icon=inline-start]:pl-1.5 [&_svg:not([class*='size-'])]:size-3.5",
        lg: "h-9 gap-1.5 px-2.5 has-data-[icon=inline-end]:pr-2 has-data-[icon=inline-start]:pl-2",
        icon: "size-8",
        "icon-xs": "size-6 rounded-none [&_svg:not([class*='size-'])]:size-3",
        "icon-sm": "size-7 rounded-none",
        "icon-lg": "size-9"
      }
    },
    defaultVariants: {
      variant: "default",
      size: "default"
    }
  }
);
function Button({
  className,
  variant = "default",
  size = "default",
  ...props
}) {
  return /* @__PURE__ */ jsx(
    Button$1,
    {
      "data-slot": "button",
      className: cn(buttonVariants({ variant, size, className })),
      ...props
    }
  );
}
function Dialog({ ...props }) {
  return /* @__PURE__ */ jsx(Dialog$1.Root, { "data-slot": "dialog", ...props });
}
function DialogPortal({ ...props }) {
  return /* @__PURE__ */ jsx(Dialog$1.Portal, { "data-slot": "dialog-portal", ...props });
}
function DialogOverlay({
  className,
  ...props
}) {
  return /* @__PURE__ */ jsx(
    Dialog$1.Backdrop,
    {
      "data-slot": "dialog-overlay",
      className: cn(
        "fixed inset-0 isolate z-50 bg-black/10 duration-100 supports-backdrop-filter:backdrop-blur-xs data-open:animate-in data-open:fade-in-0 data-closed:animate-out data-closed:fade-out-0",
        className
      ),
      ...props
    }
  );
}
function DialogContent({
  className,
  children,
  showCloseButton = true,
  ...props
}) {
  return /* @__PURE__ */ jsxs(DialogPortal, { children: [
    /* @__PURE__ */ jsx(DialogOverlay, {}),
    /* @__PURE__ */ jsxs(
      Dialog$1.Popup,
      {
        "data-slot": "dialog-content",
        className: cn(
          "fixed top-1/2 left-1/2 z-50 grid w-full max-w-[calc(100%-2rem)] -translate-x-1/2 -translate-y-1/2 gap-4 rounded-none bg-popover p-4 text-xs/relaxed text-popover-foreground ring-1 ring-foreground/10 duration-100 outline-none sm:max-w-sm data-open:animate-in data-open:fade-in-0 data-open:zoom-in-95 data-closed:animate-out data-closed:fade-out-0 data-closed:zoom-out-95",
          className
        ),
        ...props,
        children: [
          children,
          showCloseButton && /* @__PURE__ */ jsxs(
            Dialog$1.Close,
            {
              "data-slot": "dialog-close",
              render: /* @__PURE__ */ jsx(
                Button,
                {
                  variant: "ghost",
                  className: "absolute top-2 right-2",
                  size: "icon-sm"
                }
              ),
              children: [
                /* @__PURE__ */ jsx(
                  XIcon,
                  {}
                ),
                /* @__PURE__ */ jsx("span", { className: "sr-only", children: "Close" })
              ]
            }
          )
        ]
      }
    )
  ] });
}
function DialogHeader({ className, ...props }) {
  return /* @__PURE__ */ jsx(
    "div",
    {
      "data-slot": "dialog-header",
      className: cn("flex flex-col gap-1 text-left", className),
      ...props
    }
  );
}
function DialogTitle({ className, ...props }) {
  return /* @__PURE__ */ jsx(
    Dialog$1.Title,
    {
      "data-slot": "dialog-title",
      className: cn("font-heading text-sm font-medium", className),
      ...props
    }
  );
}
function DialogDescription({
  className,
  ...props
}) {
  return /* @__PURE__ */ jsx(
    Dialog$1.Description,
    {
      "data-slot": "dialog-description",
      className: cn(
        "text-xs/relaxed text-muted-foreground *:[a]:underline *:[a]:underline-offset-3 *:[a]:hover:text-foreground",
        className
      ),
      ...props
    }
  );
}
function InputGroup({ className, ...props }) {
  return /* @__PURE__ */ jsx(
    "div",
    {
      "data-slot": "input-group",
      role: "group",
      className: cn(
        "group/input-group relative flex h-8 w-full min-w-0 items-center rounded-none border border-input transition-colors outline-none in-data-[slot=combobox-content]:focus-within:border-inherit in-data-[slot=combobox-content]:focus-within:ring-0 has-disabled:bg-input/50 has-disabled:opacity-50 has-[[data-slot=input-group-control]:focus-visible]:border-ring has-[[data-slot=input-group-control]:focus-visible]:ring-1 has-[[data-slot=input-group-control]:focus-visible]:ring-ring/50 has-[[data-slot][aria-invalid=true]]:border-destructive has-[[data-slot][aria-invalid=true]]:ring-1 has-[[data-slot][aria-invalid=true]]:ring-destructive/20 has-[>[data-align=block-end]]:h-auto has-[>[data-align=block-end]]:flex-col has-[>[data-align=block-start]]:h-auto has-[>[data-align=block-start]]:flex-col has-[>textarea]:h-auto dark:bg-input/30 dark:has-disabled:bg-input/80 dark:has-[[data-slot][aria-invalid=true]]:ring-destructive/40 has-[>[data-align=block-end]]:[&>input]:pt-3 has-[>[data-align=block-start]]:[&>input]:pb-3 has-[>[data-align=inline-end]]:[&>input]:pr-1.5 has-[>[data-align=inline-start]]:[&>input]:pl-1.5",
        className
      ),
      ...props
    }
  );
}
const inputGroupAddonVariants = cva(
  "flex h-auto cursor-text items-center justify-center gap-2 py-1.5 text-xs font-medium text-muted-foreground select-none group-data-[disabled=true]/input-group:opacity-50 [&>kbd]:rounded-none [&>svg:not([class*='size-'])]:size-4",
  {
    variants: {
      align: {
        "inline-start": "order-first pl-2 has-[>button]:ml-[-0.3rem] has-[>kbd]:ml-[-0.15rem]",
        "inline-end": "order-last pr-2 has-[>button]:mr-[-0.3rem] has-[>kbd]:mr-[-0.15rem]",
        "block-start": "order-first w-full justify-start px-2.5 pt-2 group-has-[>input]/input-group:pt-2 [.border-b]:pb-2",
        "block-end": "order-last w-full justify-start px-2.5 pb-2 group-has-[>input]/input-group:pb-2 [.border-t]:pt-2"
      }
    },
    defaultVariants: {
      align: "inline-start"
    }
  }
);
function InputGroupAddon({
  className,
  align = "inline-start",
  ...props
}) {
  return /* @__PURE__ */ jsx(
    "div",
    {
      role: "group",
      "data-slot": "input-group-addon",
      "data-align": align,
      className: cn(inputGroupAddonVariants({ align }), className),
      onClick: (e) => {
        if (e.target.closest("button")) {
          return;
        }
        e.currentTarget.parentElement?.querySelector("input")?.focus();
      },
      ...props
    }
  );
}
cva(
  "flex items-center gap-2 text-xs shadow-none",
  {
    variants: {
      size: {
        xs: "h-6 gap-1 rounded-none px-1.5 [&>svg:not([class*='size-'])]:size-3.5",
        sm: "gap-1",
        "icon-xs": "size-6 rounded-none p-0 has-[>svg]:p-0",
        "icon-sm": "size-7 p-0 has-[>svg]:p-0"
      }
    },
    defaultVariants: {
      size: "xs"
    }
  }
);
function Command({
  className,
  ...props
}) {
  return /* @__PURE__ */ jsx(
    Command$1,
    {
      "data-slot": "command",
      className: cn(
        "flex size-full flex-col overflow-hidden rounded-none bg-popover text-popover-foreground",
        className
      ),
      ...props
    }
  );
}
function CommandDialog({
  title = "Command Palette",
  description = "Search for a command to run...",
  children,
  className,
  showCloseButton = false,
  ...props
}) {
  return /* @__PURE__ */ jsxs(Dialog, { ...props, children: [
    /* @__PURE__ */ jsxs(DialogHeader, { className: "sr-only", children: [
      /* @__PURE__ */ jsx(DialogTitle, { children: title }),
      /* @__PURE__ */ jsx(DialogDescription, { children: description })
    ] }),
    /* @__PURE__ */ jsx(
      DialogContent,
      {
        className: cn(
          "top-1/3 translate-y-0 overflow-hidden rounded-none p-0",
          className
        ),
        showCloseButton,
        children
      }
    )
  ] });
}
function CommandInput({
  className,
  ...props
}) {
  return /* @__PURE__ */ jsx("div", { "data-slot": "command-input-wrapper", className: "border-b pb-0", children: /* @__PURE__ */ jsxs(InputGroup, { className: "h-8 border-none border-input/30 bg-input/30 shadow-none! *:data-[slot=input-group-addon]:pl-2!", children: [
    /* @__PURE__ */ jsx(
      Command$1.Input,
      {
        "data-slot": "command-input",
        className: cn(
          "w-full text-xs outline-hidden disabled:cursor-not-allowed disabled:opacity-50",
          className
        ),
        ...props
      }
    ),
    /* @__PURE__ */ jsx(InputGroupAddon, { children: /* @__PURE__ */ jsx(SearchIcon, { className: "size-4 shrink-0 opacity-50" }) })
  ] }) });
}
function CommandList({
  className,
  ...props
}) {
  return /* @__PURE__ */ jsx(
    Command$1.List,
    {
      "data-slot": "command-list",
      className: cn(
        "no-scrollbar max-h-72 scroll-py-0 overflow-x-hidden overflow-y-auto outline-none",
        className
      ),
      ...props
    }
  );
}
function CommandEmpty({
  className,
  ...props
}) {
  return /* @__PURE__ */ jsx(
    Command$1.Empty,
    {
      "data-slot": "command-empty",
      className: cn("py-6 text-center text-xs", className),
      ...props
    }
  );
}
function CommandGroup({
  className,
  ...props
}) {
  return /* @__PURE__ */ jsx(
    Command$1.Group,
    {
      "data-slot": "command-group",
      className: cn(
        "overflow-hidden text-foreground **:[[cmdk-group-heading]]:px-2 **:[[cmdk-group-heading]]:py-1.5 **:[[cmdk-group-heading]]:text-xs **:[[cmdk-group-heading]]:text-muted-foreground",
        className
      ),
      ...props
    }
  );
}
function CommandItem({
  className,
  children,
  ...props
}) {
  return /* @__PURE__ */ jsxs(
    Command$1.Item,
    {
      "data-slot": "command-item",
      className: cn(
        "group/command-item relative flex cursor-default items-center gap-2 rounded-none px-2 py-2 text-xs outline-hidden select-none in-data-[slot=dialog-content]:rounded-none! data-[disabled=true]:pointer-events-none data-[disabled=true]:opacity-50 data-selected:bg-muted data-selected:text-foreground [&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*='size-'])]:size-4 data-selected:*:[svg]:text-foreground",
        className
      ),
      ...props,
      children: [
        children,
        /* @__PURE__ */ jsx(CheckIcon, { className: "ml-auto opacity-0 group-has-data-[slot=command-shortcut]/command-item:hidden group-data-[checked=true]/command-item:opacity-100" })
      ]
    }
  );
}
function DocsSearch({
  versions,
  currentVersion,
  className
}) {
  const [open, setOpen] = useState(false);
  const catalog = versions.find((version) => version.name === currentVersion) ?? versions[0];
  useEffect(() => {
    function onKeyDown(event) {
      if ((event.metaKey || event.ctrlKey) && event.key.toLowerCase() === "k") {
        event.preventDefault();
        setOpen((current) => !current);
      }
    }
    window.addEventListener("keydown", onKeyDown);
    return () => window.removeEventListener("keydown", onKeyDown);
  }, []);
  function visit(url) {
    setOpen(false);
    router.visit(url);
  }
  return /* @__PURE__ */ jsxs(Fragment, { children: [
    /* @__PURE__ */ jsxs(
      Button,
      {
        type: "button",
        variant: "outline",
        className: cn("h-8 w-full justify-start text-muted-foreground", className),
        onClick: () => setOpen(true),
        children: [
          /* @__PURE__ */ jsx(SearchIcon, {}),
          /* @__PURE__ */ jsx("span", { className: "flex-1 truncate text-left", children: "Search documentation..." }),
          /* @__PURE__ */ jsx("kbd", { className: "pointer-events-none hidden h-5 items-center gap-1 border border-border bg-muted px-1.5 font-mono text-[0.65rem] text-muted-foreground sm:inline-flex", children: "⌘K" })
        ]
      }
    ),
    /* @__PURE__ */ jsx(
      CommandDialog,
      {
        open,
        onOpenChange: setOpen,
        title: "Search documentation",
        description: "Find a page in the Andurel docs.",
        children: /* @__PURE__ */ jsxs(Command, { children: [
          /* @__PURE__ */ jsx(CommandInput, { placeholder: "Search documentation..." }),
          /* @__PURE__ */ jsxs(CommandList, { children: [
            /* @__PURE__ */ jsx(CommandEmpty, { children: "No documentation found." }),
            catalog?.sections.map((section) => /* @__PURE__ */ jsx(CommandGroup, { heading: section.title, children: section.pages.map((page) => /* @__PURE__ */ jsx(
              CommandItem,
              {
                value: `${section.title} ${page.title} ${page.description}`,
                onSelect: () => visit(page.url),
                children: /* @__PURE__ */ jsxs("span", { className: "flex min-w-0 flex-col", children: [
                  /* @__PURE__ */ jsx("span", { children: page.title }),
                  /* @__PURE__ */ jsx("span", { className: "truncate text-muted-foreground", children: page.description })
                ] })
              },
              page.url
            )) }, section.title))
          ] })
        ] })
      }
    )
  ] });
}
function collectHeadings(rootId) {
  const root = document.getElementById(rootId);
  if (!root) {
    return [];
  }
  return [...root.querySelectorAll("h1, h2, h3")].map((element) => {
    const title = element.textContent?.trim() ?? "";
    if (!title) {
      return null;
    }
    if (!element.id) {
      element.id = title.toLowerCase().replace(/[^a-z0-9]+/g, "-").replace(/^-|-$/g, "");
    }
    return {
      id: element.id,
      title,
      level: Number(element.tagName[1])
    };
  }).filter((heading) => heading !== null);
}
function preferSectionHeadings(headings) {
  const sections = headings.filter((heading) => heading.level >= 2);
  return sections.length > 0 ? sections : headings;
}
function DocsToc({ rootId, pageKey }) {
  const [headings, setHeadings] = useState([]);
  const [activeId, setActiveId] = useState("");
  useEffect(() => {
    const nextHeadings = preferSectionHeadings(collectHeadings(rootId));
    setHeadings(nextHeadings);
    setActiveId(nextHeadings[0]?.id ?? "");
    if (nextHeadings.length === 0) {
      return;
    }
    const observed = nextHeadings.map((heading) => document.getElementById(heading.id)).filter((element) => element !== null);
    const observer = new IntersectionObserver(
      (entries) => {
        const visible = entries.filter((entry) => entry.isIntersecting).sort((left, right) => right.intersectionRatio - left.intersectionRatio);
        if (visible[0]?.target.id) {
          setActiveId(visible[0].target.id);
        }
      },
      {
        rootMargin: "-80px 0px -60% 0px",
        threshold: [0, 1]
      }
    );
    observed.forEach((element) => observer.observe(element));
    return () => observer.disconnect();
  }, [pageKey, rootId]);
  if (headings.length === 0) {
    return null;
  }
  return /* @__PURE__ */ jsxs("nav", { "aria-label": "On this page", className: "sticky top-14", children: [
    /* @__PURE__ */ jsxs("p", { className: "mb-3 flex items-center gap-2 text-xs font-semibold text-foreground", children: [
      /* @__PURE__ */ jsx(AlignLeftIcon, { className: "size-3.5 text-muted-foreground" }),
      "On this page"
    ] }),
    /* @__PURE__ */ jsx("ul", { className: "border-l border-sidebar-border", children: headings.map((heading) => /* @__PURE__ */ jsx("li", { children: /* @__PURE__ */ jsx(
      "a",
      {
        href: `#${heading.id}`,
        className: cn(
          "-ml-px block border-l-2 py-1 text-xs leading-5 transition-colors",
          heading.level > 2 ? "pl-6" : "pl-3",
          activeId === heading.id ? "border-accent text-foreground" : "border-transparent text-muted-foreground hover:text-foreground"
        ),
        children: heading.title
      }
    ) }, heading.id)) })
  ] });
}
function DropdownMenu({ ...props }) {
  return /* @__PURE__ */ jsx(Menu.Root, { "data-slot": "dropdown-menu", ...props });
}
function DropdownMenuTrigger({ ...props }) {
  return /* @__PURE__ */ jsx(Menu.Trigger, { "data-slot": "dropdown-menu-trigger", ...props });
}
function DropdownMenuContent({
  align = "start",
  alignOffset = 0,
  side = "bottom",
  sideOffset = 4,
  className,
  ...props
}) {
  return /* @__PURE__ */ jsx(Menu.Portal, { children: /* @__PURE__ */ jsx(
    Menu.Positioner,
    {
      className: "isolate z-50 outline-none",
      align,
      alignOffset,
      side,
      sideOffset,
      children: /* @__PURE__ */ jsx(
        Menu.Popup,
        {
          "data-slot": "dropdown-menu-content",
          className: cn("z-50 max-h-(--available-height) w-(--anchor-width) min-w-32 origin-(--transform-origin) overflow-x-hidden overflow-y-auto rounded-none bg-popover text-popover-foreground shadow-md ring-1 ring-foreground/10 duration-100 outline-none data-[side=bottom]:slide-in-from-top-2 data-[side=inline-end]:slide-in-from-left-2 data-[side=inline-start]:slide-in-from-right-2 data-[side=left]:slide-in-from-right-2 data-[side=right]:slide-in-from-left-2 data-[side=top]:slide-in-from-bottom-2 data-open:animate-in data-open:fade-in-0 data-open:zoom-in-95 data-closed:animate-out data-closed:overflow-hidden data-closed:fade-out-0 data-closed:zoom-out-95", className),
          ...props
        }
      )
    }
  ) });
}
function DropdownMenuGroup({ ...props }) {
  return /* @__PURE__ */ jsx(Menu.Group, { "data-slot": "dropdown-menu-group", ...props });
}
function DropdownMenuLabel({
  className,
  inset,
  ...props
}) {
  return /* @__PURE__ */ jsx(
    Menu.GroupLabel,
    {
      "data-slot": "dropdown-menu-label",
      "data-inset": inset,
      className: cn(
        "px-2 py-2 text-xs text-muted-foreground data-inset:pl-7",
        className
      ),
      ...props
    }
  );
}
function DropdownMenuItem({
  className,
  inset,
  variant = "default",
  ...props
}) {
  return /* @__PURE__ */ jsx(
    Menu.Item,
    {
      "data-slot": "dropdown-menu-item",
      "data-inset": inset,
      "data-variant": variant,
      className: cn(
        "group/dropdown-menu-item relative flex cursor-default items-center gap-2 rounded-none px-2 py-2 text-xs outline-hidden select-none focus:bg-accent focus:text-accent-foreground not-data-[variant=destructive]:focus:**:text-accent-foreground data-inset:pl-7 data-[variant=destructive]:text-destructive data-[variant=destructive]:focus:bg-destructive/10 data-[variant=destructive]:focus:text-destructive dark:data-[variant=destructive]:focus:bg-destructive/20 data-disabled:pointer-events-none data-disabled:opacity-50 [&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*='size-'])]:size-4 data-[variant=destructive]:*:[svg]:text-destructive",
        className
      ),
      ...props
    }
  );
}
const MOBILE_BREAKPOINT = 768;
function useIsMobile() {
  const [isMobile, setIsMobile] = React.useState(void 0);
  React.useEffect(() => {
    const mql = window.matchMedia(`(max-width: ${MOBILE_BREAKPOINT - 1}px)`);
    const onChange = () => {
      setIsMobile(window.innerWidth < MOBILE_BREAKPOINT);
    };
    mql.addEventListener("change", onChange);
    setIsMobile(window.innerWidth < MOBILE_BREAKPOINT);
    return () => mql.removeEventListener("change", onChange);
  }, []);
  return !!isMobile;
}
function Sheet({ ...props }) {
  return /* @__PURE__ */ jsx(Dialog$1.Root, { "data-slot": "sheet", ...props });
}
function SheetPortal({ ...props }) {
  return /* @__PURE__ */ jsx(Dialog$1.Portal, { "data-slot": "sheet-portal", ...props });
}
function SheetOverlay({ className, ...props }) {
  return /* @__PURE__ */ jsx(
    Dialog$1.Backdrop,
    {
      "data-slot": "sheet-overlay",
      className: cn(
        "fixed inset-0 z-50 bg-black/10 text-xs/relaxed transition-opacity duration-150 data-ending-style:opacity-0 data-starting-style:opacity-0 supports-backdrop-filter:backdrop-blur-xs",
        className
      ),
      ...props
    }
  );
}
function SheetContent({
  className,
  children,
  side = "right",
  showCloseButton = true,
  ...props
}) {
  return /* @__PURE__ */ jsxs(SheetPortal, { children: [
    /* @__PURE__ */ jsx(SheetOverlay, {}),
    /* @__PURE__ */ jsxs(
      Dialog$1.Popup,
      {
        "data-slot": "sheet-content",
        "data-side": side,
        className: cn(
          "fixed z-50 flex flex-col bg-popover bg-clip-padding text-xs/relaxed text-popover-foreground shadow-lg transition duration-200 ease-in-out data-ending-style:opacity-0 data-starting-style:opacity-0 data-[side=bottom]:inset-x-0 data-[side=bottom]:bottom-0 data-[side=bottom]:h-auto data-[side=bottom]:border-t data-[side=bottom]:data-ending-style:translate-y-[2.5rem] data-[side=bottom]:data-starting-style:translate-y-[2.5rem] data-[side=left]:inset-y-0 data-[side=left]:left-0 data-[side=left]:h-full data-[side=left]:w-3/4 data-[side=left]:border-r data-[side=left]:data-ending-style:translate-x-[-2.5rem] data-[side=left]:data-starting-style:translate-x-[-2.5rem] data-[side=right]:inset-y-0 data-[side=right]:right-0 data-[side=right]:h-full data-[side=right]:w-3/4 data-[side=right]:border-l data-[side=right]:data-ending-style:translate-x-[2.5rem] data-[side=right]:data-starting-style:translate-x-[2.5rem] data-[side=top]:inset-x-0 data-[side=top]:top-0 data-[side=top]:h-auto data-[side=top]:border-b data-[side=top]:data-ending-style:translate-y-[-2.5rem] data-[side=top]:data-starting-style:translate-y-[-2.5rem] data-[side=left]:sm:max-w-sm data-[side=right]:sm:max-w-sm",
          className
        ),
        ...props,
        children: [
          children,
          showCloseButton && /* @__PURE__ */ jsxs(
            Dialog$1.Close,
            {
              "data-slot": "sheet-close",
              render: /* @__PURE__ */ jsx(
                Button,
                {
                  variant: "ghost",
                  className: "absolute top-3 right-3",
                  size: "icon-sm"
                }
              ),
              children: [
                /* @__PURE__ */ jsx(
                  XIcon,
                  {}
                ),
                /* @__PURE__ */ jsx("span", { className: "sr-only", children: "Close" })
              ]
            }
          )
        ]
      }
    )
  ] });
}
function SheetHeader({ className, ...props }) {
  return /* @__PURE__ */ jsx(
    "div",
    {
      "data-slot": "sheet-header",
      className: cn("flex flex-col gap-0.5 p-4", className),
      ...props
    }
  );
}
function SheetTitle({ className, ...props }) {
  return /* @__PURE__ */ jsx(
    Dialog$1.Title,
    {
      "data-slot": "sheet-title",
      className: cn(
        "font-heading text-sm font-medium text-foreground",
        className
      ),
      ...props
    }
  );
}
function SheetDescription({
  className,
  ...props
}) {
  return /* @__PURE__ */ jsx(
    Dialog$1.Description,
    {
      "data-slot": "sheet-description",
      className: cn("text-xs/relaxed text-muted-foreground", className),
      ...props
    }
  );
}
function Tooltip({ ...props }) {
  return /* @__PURE__ */ jsx(Tooltip$1.Root, { "data-slot": "tooltip", ...props });
}
function TooltipTrigger({ ...props }) {
  return /* @__PURE__ */ jsx(Tooltip$1.Trigger, { "data-slot": "tooltip-trigger", ...props });
}
function TooltipContent({
  className,
  side = "top",
  sideOffset = 4,
  align = "center",
  alignOffset = 0,
  children,
  ...props
}) {
  return /* @__PURE__ */ jsx(Tooltip$1.Portal, { children: /* @__PURE__ */ jsx(
    Tooltip$1.Positioner,
    {
      align,
      alignOffset,
      side,
      sideOffset,
      className: "isolate z-50",
      children: /* @__PURE__ */ jsxs(
        Tooltip$1.Popup,
        {
          "data-slot": "tooltip-content",
          className: cn(
            "z-50 inline-flex w-fit max-w-xs origin-(--transform-origin) items-center gap-1.5 rounded-none bg-foreground px-3 py-1.5 text-xs text-background has-data-[slot=kbd]:pr-1.5 data-[side=bottom]:slide-in-from-top-2 data-[side=inline-end]:slide-in-from-left-2 data-[side=inline-start]:slide-in-from-right-2 data-[side=left]:slide-in-from-right-2 data-[side=right]:slide-in-from-left-2 data-[side=top]:slide-in-from-bottom-2 **:data-[slot=kbd]:relative **:data-[slot=kbd]:isolate **:data-[slot=kbd]:z-50 **:data-[slot=kbd]:rounded-none data-[state=delayed-open]:animate-in data-[state=delayed-open]:fade-in-0 data-[state=delayed-open]:zoom-in-95 data-open:animate-in data-open:fade-in-0 data-open:zoom-in-95 data-closed:animate-out data-closed:fade-out-0 data-closed:zoom-out-95",
            className
          ),
          ...props,
          children: [
            children,
            /* @__PURE__ */ jsx(Tooltip$1.Arrow, { className: "z-50 size-2.5 translate-y-[calc(-50%-2px)] rotate-45 rounded-none bg-foreground fill-foreground data-[side=bottom]:top-1 data-[side=inline-end]:top-1/2! data-[side=inline-end]:-left-1 data-[side=inline-end]:-translate-y-1/2 data-[side=inline-start]:top-1/2! data-[side=inline-start]:-right-1 data-[side=inline-start]:-translate-y-1/2 data-[side=left]:top-1/2! data-[side=left]:-right-1 data-[side=left]:-translate-y-1/2 data-[side=right]:top-1/2! data-[side=right]:-left-1 data-[side=right]:-translate-y-1/2 data-[side=top]:-bottom-2.5" })
          ]
        }
      )
    }
  ) });
}
const SIDEBAR_COOKIE_NAME = "sidebar_state";
const SIDEBAR_COOKIE_MAX_AGE = 60 * 60 * 24 * 7;
const SIDEBAR_WIDTH = "16rem";
const SIDEBAR_WIDTH_MOBILE = "18rem";
const SIDEBAR_WIDTH_ICON = "3rem";
const SIDEBAR_KEYBOARD_SHORTCUT = "b";
const SidebarContext = React.createContext(null);
function useSidebar() {
  const context = React.useContext(SidebarContext);
  if (!context) {
    throw new Error("useSidebar must be used within a SidebarProvider.");
  }
  return context;
}
function SidebarProvider({
  defaultOpen = true,
  open: openProp,
  onOpenChange: setOpenProp,
  className,
  style,
  children,
  ...props
}) {
  const isMobile = useIsMobile();
  const [openMobile, setOpenMobile] = React.useState(false);
  const [_open, _setOpen] = React.useState(defaultOpen);
  const open = openProp ?? _open;
  const setOpen = React.useCallback(
    (value) => {
      const openState = typeof value === "function" ? value(open) : value;
      if (setOpenProp) {
        setOpenProp(openState);
      } else {
        _setOpen(openState);
      }
      document.cookie = `${SIDEBAR_COOKIE_NAME}=${openState}; path=/; max-age=${SIDEBAR_COOKIE_MAX_AGE}`;
    },
    [setOpenProp, open]
  );
  const toggleSidebar = React.useCallback(() => {
    return isMobile ? setOpenMobile((open2) => !open2) : setOpen((open2) => !open2);
  }, [isMobile, setOpen, setOpenMobile]);
  React.useEffect(() => {
    const handleKeyDown = (event) => {
      if (event.key === SIDEBAR_KEYBOARD_SHORTCUT && (event.metaKey || event.ctrlKey)) {
        event.preventDefault();
        toggleSidebar();
      }
    };
    window.addEventListener("keydown", handleKeyDown);
    return () => window.removeEventListener("keydown", handleKeyDown);
  }, [toggleSidebar]);
  const state = open ? "expanded" : "collapsed";
  const contextValue = React.useMemo(
    () => ({
      state,
      open,
      setOpen,
      isMobile,
      openMobile,
      setOpenMobile,
      toggleSidebar
    }),
    [state, open, setOpen, isMobile, openMobile, setOpenMobile, toggleSidebar]
  );
  return /* @__PURE__ */ jsx(SidebarContext.Provider, { value: contextValue, children: /* @__PURE__ */ jsx(
    "div",
    {
      "data-slot": "sidebar-wrapper",
      style: {
        "--sidebar-width": SIDEBAR_WIDTH,
        "--sidebar-width-icon": SIDEBAR_WIDTH_ICON,
        ...style
      },
      className: cn(
        "group/sidebar-wrapper flex min-h-svh w-full has-data-[variant=inset]:bg-sidebar",
        className
      ),
      ...props,
      children
    }
  ) });
}
function Sidebar({
  side = "left",
  variant = "sidebar",
  collapsible = "offcanvas",
  className,
  children,
  dir,
  ...props
}) {
  const { isMobile, state, openMobile, setOpenMobile } = useSidebar();
  if (collapsible === "none") {
    return /* @__PURE__ */ jsx(
      "div",
      {
        "data-slot": "sidebar",
        className: cn(
          "flex h-full w-(--sidebar-width) flex-col bg-sidebar text-sidebar-foreground",
          className
        ),
        ...props,
        children
      }
    );
  }
  if (isMobile) {
    return /* @__PURE__ */ jsx(Sheet, { open: openMobile, onOpenChange: setOpenMobile, ...props, children: /* @__PURE__ */ jsxs(
      SheetContent,
      {
        dir,
        "data-sidebar": "sidebar",
        "data-slot": "sidebar",
        "data-mobile": "true",
        className: "w-(--sidebar-width) bg-sidebar p-0 text-sidebar-foreground [&>button]:hidden",
        style: {
          "--sidebar-width": SIDEBAR_WIDTH_MOBILE
        },
        side,
        children: [
          /* @__PURE__ */ jsxs(SheetHeader, { className: "sr-only", children: [
            /* @__PURE__ */ jsx(SheetTitle, { children: "Sidebar" }),
            /* @__PURE__ */ jsx(SheetDescription, { children: "Displays the mobile sidebar." })
          ] }),
          /* @__PURE__ */ jsx("div", { className: "flex h-full w-full flex-col", children })
        ]
      }
    ) });
  }
  return /* @__PURE__ */ jsxs(
    "div",
    {
      className: "group peer hidden text-sidebar-foreground md:block",
      "data-state": state,
      "data-collapsible": state === "collapsed" ? collapsible : "",
      "data-variant": variant,
      "data-side": side,
      "data-slot": "sidebar",
      children: [
        /* @__PURE__ */ jsx(
          "div",
          {
            "data-slot": "sidebar-gap",
            className: cn(
              "relative w-(--sidebar-width) bg-transparent transition-[width] duration-200 ease-linear",
              "group-data-[collapsible=offcanvas]:w-0",
              "group-data-[side=right]:rotate-180",
              variant === "floating" || variant === "inset" ? "group-data-[collapsible=icon]:w-[calc(var(--sidebar-width-icon)+(--spacing(4)))]" : "group-data-[collapsible=icon]:w-(--sidebar-width-icon)"
            )
          }
        ),
        /* @__PURE__ */ jsx(
          "div",
          {
            "data-slot": "sidebar-container",
            "data-side": side,
            className: cn(
              "fixed inset-y-0 z-10 hidden h-svh w-(--sidebar-width) transition-[left,right,width] duration-200 ease-linear data-[side=left]:left-0 data-[side=left]:group-data-[collapsible=offcanvas]:left-[calc(var(--sidebar-width)*-1)] data-[side=right]:right-0 data-[side=right]:group-data-[collapsible=offcanvas]:right-[calc(var(--sidebar-width)*-1)] md:flex",
              // Adjust the padding for floating and inset variants.
              variant === "floating" || variant === "inset" ? "p-2 group-data-[collapsible=icon]:w-[calc(var(--sidebar-width-icon)+(--spacing(4))+2px)]" : "group-data-[collapsible=icon]:w-(--sidebar-width-icon) group-data-[side=left]:border-r group-data-[side=right]:border-l",
              className
            ),
            ...props,
            children: /* @__PURE__ */ jsx(
              "div",
              {
                "data-sidebar": "sidebar",
                "data-slot": "sidebar-inner",
                className: "flex size-full flex-col bg-sidebar group-data-[variant=floating]:rounded-none group-data-[variant=floating]:shadow-sm group-data-[variant=floating]:ring-1 group-data-[variant=floating]:ring-sidebar-border",
                children
              }
            )
          }
        )
      ]
    }
  );
}
function SidebarTrigger({
  className,
  onClick,
  ...props
}) {
  const { toggleSidebar } = useSidebar();
  return /* @__PURE__ */ jsxs(
    Button,
    {
      "data-sidebar": "trigger",
      "data-slot": "sidebar-trigger",
      variant: "ghost",
      size: "icon-sm",
      className: cn(className),
      onClick: (event) => {
        onClick?.(event);
        toggleSidebar();
      },
      ...props,
      children: [
        /* @__PURE__ */ jsx(PanelLeftIcon, {}),
        /* @__PURE__ */ jsx("span", { className: "sr-only", children: "Toggle Sidebar" })
      ]
    }
  );
}
function SidebarInset({ className, ...props }) {
  return /* @__PURE__ */ jsx(
    "main",
    {
      "data-slot": "sidebar-inset",
      className: cn(
        "relative flex w-full flex-1 flex-col bg-background md:peer-data-[variant=inset]:m-2 md:peer-data-[variant=inset]:ml-0 md:peer-data-[variant=inset]:rounded-none md:peer-data-[variant=inset]:shadow-sm md:peer-data-[variant=inset]:peer-data-[state=collapsed]:ml-2",
        className
      ),
      ...props
    }
  );
}
function SidebarHeader({ className, ...props }) {
  return /* @__PURE__ */ jsx(
    "div",
    {
      "data-slot": "sidebar-header",
      "data-sidebar": "header",
      className: cn("flex flex-col gap-2 p-2", className),
      ...props
    }
  );
}
function SidebarContent({ className, ...props }) {
  return /* @__PURE__ */ jsx(
    "div",
    {
      "data-slot": "sidebar-content",
      "data-sidebar": "content",
      className: cn(
        "no-scrollbar flex min-h-0 flex-1 flex-col gap-0 overflow-auto group-data-[collapsible=icon]:overflow-hidden",
        className
      ),
      ...props
    }
  );
}
function SidebarGroup({ className, ...props }) {
  return /* @__PURE__ */ jsx(
    "div",
    {
      "data-slot": "sidebar-group",
      "data-sidebar": "group",
      className: cn("relative flex w-full min-w-0 flex-col p-2", className),
      ...props
    }
  );
}
function SidebarGroupLabel({
  className,
  render,
  ...props
}) {
  return useRender({
    defaultTagName: "div",
    props: mergeProps(
      {
        className: cn(
          "flex h-8 shrink-0 items-center rounded-none px-2 text-xs text-sidebar-foreground/70 ring-sidebar-ring outline-hidden transition-[margin,opacity] duration-200 ease-linear group-data-[collapsible=icon]:-mt-8 group-data-[collapsible=icon]:opacity-0 focus-visible:ring-2 [&>svg]:size-4 [&>svg]:shrink-0",
          className
        )
      },
      props
    ),
    render,
    state: {
      slot: "sidebar-group-label",
      sidebar: "group-label"
    }
  });
}
function SidebarGroupContent({
  className,
  ...props
}) {
  return /* @__PURE__ */ jsx(
    "div",
    {
      "data-slot": "sidebar-group-content",
      "data-sidebar": "group-content",
      className: cn("w-full text-xs", className),
      ...props
    }
  );
}
function SidebarMenu({ className, ...props }) {
  return /* @__PURE__ */ jsx(
    "ul",
    {
      "data-slot": "sidebar-menu",
      "data-sidebar": "menu",
      className: cn("flex w-full min-w-0 flex-col gap-0", className),
      ...props
    }
  );
}
function SidebarMenuItem({ className, ...props }) {
  return /* @__PURE__ */ jsx(
    "li",
    {
      "data-slot": "sidebar-menu-item",
      "data-sidebar": "menu-item",
      className: cn("group/menu-item relative", className),
      ...props
    }
  );
}
const sidebarMenuButtonVariants = cva(
  "peer/menu-button group/menu-button flex w-full items-center gap-2 overflow-hidden rounded-none p-2 text-left text-xs ring-sidebar-ring outline-hidden transition-[width,height,padding] group-has-data-[sidebar=menu-action]/menu-item:pr-8 group-data-[collapsible=icon]:size-8! group-data-[collapsible=icon]:p-2! hover:bg-sidebar-accent hover:text-sidebar-accent-foreground focus-visible:ring-2 active:bg-sidebar-accent active:text-sidebar-accent-foreground disabled:pointer-events-none disabled:opacity-50 aria-disabled:pointer-events-none aria-disabled:opacity-50 data-open:hover:bg-sidebar-accent data-open:hover:text-sidebar-accent-foreground data-active:bg-sidebar-accent data-active:font-medium data-active:text-sidebar-accent-foreground [&_svg]:size-4 [&_svg]:shrink-0 [&>span:last-child]:truncate",
  {
    variants: {
      variant: {
        default: "hover:bg-sidebar-accent hover:text-sidebar-accent-foreground",
        outline: "bg-background shadow-[0_0_0_1px_var(--sidebar-border)] hover:bg-sidebar-accent hover:text-sidebar-accent-foreground hover:shadow-[0_0_0_1px_var(--sidebar-accent)]"
      },
      size: {
        default: "h-8 text-xs",
        sm: "h-7 text-xs",
        lg: "h-12 text-xs group-data-[collapsible=icon]:p-0!"
      }
    },
    defaultVariants: {
      variant: "default",
      size: "default"
    }
  }
);
function SidebarMenuButton({
  render,
  isActive = false,
  variant = "default",
  size = "default",
  tooltip,
  className,
  ...props
}) {
  const { isMobile, state } = useSidebar();
  const comp = useRender({
    defaultTagName: "button",
    props: mergeProps(
      {
        className: cn(sidebarMenuButtonVariants({ variant, size }), className)
      },
      props
    ),
    render: !tooltip ? render : /* @__PURE__ */ jsx(TooltipTrigger, { render }),
    state: {
      slot: "sidebar-menu-button",
      sidebar: "menu-button",
      size,
      active: isActive
    }
  });
  if (!tooltip) {
    return comp;
  }
  if (typeof tooltip === "string") {
    tooltip = {
      children: tooltip
    };
  }
  return /* @__PURE__ */ jsxs(Tooltip, { children: [
    comp,
    /* @__PURE__ */ jsx(
      TooltipContent,
      {
        side: "right",
        align: "center",
        hidden: state !== "collapsed" || isMobile,
        ...tooltip
      }
    )
  ] });
}
const DOC_ARTICLE_ID = "doc-article";
const docsPad = "px-14 sm:px-16";
function currentCatalog(versions, currentVersion) {
  return versions.find((version) => version.name === currentVersion) ?? versions[0];
}
function BrandMark() {
  return /* @__PURE__ */ jsxs("span", { className: "grid size-8 shrink-0 grid-cols-2 gap-1 border border-sidebar-border bg-sidebar p-1 shadow-sm shadow-black/40", children: [
    /* @__PURE__ */ jsx("span", { className: "border border-sidebar-primary" }),
    /* @__PURE__ */ jsx("span", { className: "border border-sidebar-border" }),
    /* @__PURE__ */ jsx("span", { className: "border border-sidebar-border" }),
    /* @__PURE__ */ jsx("span", { className: "bg-sidebar-primary" })
  ] });
}
function DocBrand() {
  return /* @__PURE__ */ jsxs(
    "a",
    {
      className: "flex h-full items-center gap-3 text-sm font-semibold text-sidebar-foreground",
      href: routes.homePage(),
      children: [
        /* @__PURE__ */ jsx(BrandMark, {}),
        /* @__PURE__ */ jsx("span", { children: "Andurel." })
      ]
    }
  );
}
function HeaderActions({
  versions,
  currentVersion
}) {
  return /* @__PURE__ */ jsxs("nav", { className: "flex shrink-0 items-center gap-2 text-sm", children: [
    /* @__PURE__ */ jsxs(DropdownMenu, { children: [
      /* @__PURE__ */ jsxs(DropdownMenuTrigger, { render: /* @__PURE__ */ jsx(Button, { variant: "outline", size: "sm" }), children: [
        currentVersion,
        /* @__PURE__ */ jsx(ChevronDownIcon, {})
      ] }),
      /* @__PURE__ */ jsx(DropdownMenuContent, { align: "end", children: /* @__PURE__ */ jsxs(DropdownMenuGroup, { children: [
        /* @__PURE__ */ jsx(DropdownMenuLabel, { children: "Versions" }),
        versions.map((version) => /* @__PURE__ */ jsx(DropdownMenuItem, { render: /* @__PURE__ */ jsx(Link, { href: version.url }), children: version.name }, version.name))
      ] }) })
    ] }),
    /* @__PURE__ */ jsx(
      "a",
      {
        className: "px-2 py-1 text-muted-foreground transition hover:text-foreground",
        href: "https://github.com/mbvlabs/andurel",
        children: "GitHub"
      }
    )
  ] });
}
function DocLayout({
  children,
  versions,
  currentVersion,
  currentSlug
}) {
  const catalog = currentCatalog(versions, currentVersion);
  return /* @__PURE__ */ jsxs(SidebarProvider, { children: [
    /* @__PURE__ */ jsxs(Sidebar, { className: "border-sidebar-border", children: [
      /* @__PURE__ */ jsx(SidebarHeader, { className: "h-14 flex-row items-center border-b border-sidebar-border px-4 py-0", children: /* @__PURE__ */ jsx(DocBrand, {}) }),
      /* @__PURE__ */ jsx(SidebarContent, { children: catalog?.sections.map((section) => /* @__PURE__ */ jsxs(SidebarGroup, { children: [
        /* @__PURE__ */ jsx(SidebarGroupLabel, { children: section.title }),
        /* @__PURE__ */ jsx(SidebarGroupContent, { children: /* @__PURE__ */ jsx(SidebarMenu, { children: section.pages.map((page) => /* @__PURE__ */ jsx(SidebarMenuItem, { children: /* @__PURE__ */ jsx(
          SidebarMenuButton,
          {
            isActive: currentSlug === page.slug,
            "aria-current": currentSlug === page.slug ? "page" : void 0,
            tooltip: page.title,
            render: /* @__PURE__ */ jsx(Link, { href: page.url }),
            children: /* @__PURE__ */ jsx("span", { children: page.title })
          }
        ) }, page.slug)) }) })
      ] }, section.title)) })
    ] }),
    /* @__PURE__ */ jsxs(SidebarInset, { children: [
      /* @__PURE__ */ jsx("header", { className: "sticky top-0 z-30 border-b border-sidebar-border bg-background", children: /* @__PURE__ */ jsxs("div", { className: "relative flex h-14 items-center", children: [
        /* @__PURE__ */ jsxs("div", { className: `absolute inset-y-0 left-0 z-10 flex items-center ${docsPad}`, children: [
          /* @__PURE__ */ jsx(SidebarTrigger, {}),
          /* @__PURE__ */ jsxs(
            "a",
            {
              className: "ml-3 inline-flex items-center gap-2 text-sm font-semibold md:hidden",
              href: routes.homePage(),
              children: [
                /* @__PURE__ */ jsx(BrandMark, {}),
                /* @__PURE__ */ jsx("span", { className: "sr-only", children: "Andurel Docs" })
              ]
            }
          )
        ] }),
        /* @__PURE__ */ jsx("div", { className: `flex w-full items-center justify-center ${docsPad}`, children: /* @__PURE__ */ jsx("div", { className: "w-full max-w-3xl", children: /* @__PURE__ */ jsx(DocsSearch, { versions, currentVersion }) }) }),
        /* @__PURE__ */ jsx("div", { className: `absolute inset-y-0 right-0 z-10 flex items-center ${docsPad}`, children: /* @__PURE__ */ jsx(HeaderActions, { versions, currentVersion }) })
      ] }) }),
      /* @__PURE__ */ jsxs("div", { className: `flex flex-1 py-8 ${docsPad}`, children: [
        /* @__PURE__ */ jsx("div", { className: "hidden w-64 shrink-0 xl:block", "aria-hidden": "true" }),
        /* @__PURE__ */ jsx("div", { className: "flex min-w-0 flex-1 justify-center", children: /* @__PURE__ */ jsx("div", { id: DOC_ARTICLE_ID, className: "w-full max-w-3xl", children }) }),
        /* @__PURE__ */ jsx("aside", { className: "hidden w-64 shrink-0 xl:block", children: /* @__PURE__ */ jsx(DocsToc, { rootId: DOC_ARTICLE_ID, pageKey: `${currentVersion}:${currentSlug}` }) })
      ] })
    ] })
  ] });
}
function Show({
  versions,
  currentVersion,
  currentSlug,
  currentSection,
  title,
  description
}) {
  return /* @__PURE__ */ jsxs(
    DocLayout,
    {
      versions,
      currentVersion,
      currentSlug,
      children: [
        /* @__PURE__ */ jsx(Head, { title }),
        /* @__PURE__ */ jsxs("p", { className: "mb-7 font-mono text-xs font-semibold uppercase tracking-[0.14em] text-accent", children: [
          currentSection,
          " / ",
          currentVersion
        ] }),
        /* @__PURE__ */ jsx("h1", { id: currentSlug, className: "scroll-mt-20 text-3xl font-semibold text-foreground", children: title }),
        /* @__PURE__ */ jsx("p", { className: "mt-4 text-base leading-7 text-muted-foreground", children: description })
      ]
    }
  );
}
const __vite_glob_0_5 = /* @__PURE__ */ Object.freeze(/* @__PURE__ */ Object.defineProperty({
  __proto__: null,
  default: Show
}, Symbol.toStringTag, { value: "Module" }));
function BadRequest() {
  return /* @__PURE__ */ jsx(Layout, { children: /* @__PURE__ */ jsxs("section", { className: "w-full max-w-md border border-[#2f3a37] bg-[#101414]/90 p-6 text-center shadow-sm shadow-black/40", children: [
    /* @__PURE__ */ jsx("p", { className: "text-sm font-medium uppercase tracking-wide text-[#8df7a4]", children: "400" }),
    /* @__PURE__ */ jsx("h1", { className: "mt-2 text-2xl font-semibold text-[#f2ead8]", children: "Bad request" }),
    /* @__PURE__ */ jsx("p", { className: "mt-3 text-sm leading-6 text-[#8f8a7d]", children: "The request made was invalid." })
  ] }) });
}
const __vite_glob_0_6 = /* @__PURE__ */ Object.freeze(/* @__PURE__ */ Object.defineProperty({
  __proto__: null,
  default: BadRequest
}, Symbol.toStringTag, { value: "Module" }));
function InternalError() {
  return /* @__PURE__ */ jsx(Layout, { children: /* @__PURE__ */ jsxs("section", { className: "w-full max-w-md border border-[#2f3a37] bg-[#101414]/90 p-6 text-center shadow-sm shadow-black/40", children: [
    /* @__PURE__ */ jsx("p", { className: "text-sm font-medium uppercase tracking-wide text-[#8df7a4]", children: "500" }),
    /* @__PURE__ */ jsx("h1", { className: "mt-2 text-2xl font-semibold text-[#f2ead8]", children: "Something went wrong." }),
    /* @__PURE__ */ jsx("p", { className: "mt-3 text-sm leading-6 text-[#8f8a7d]", children: "The application hit an unexpected error." })
  ] }) });
}
const __vite_glob_0_7 = /* @__PURE__ */ Object.freeze(/* @__PURE__ */ Object.defineProperty({
  __proto__: null,
  default: InternalError
}, Symbol.toStringTag, { value: "Module" }));
function NotFound() {
  return /* @__PURE__ */ jsx(Layout, { children: /* @__PURE__ */ jsxs("section", { className: "w-full max-w-md border border-[#2f3a37] bg-[#101414]/90 p-6 text-center shadow-sm shadow-black/40", children: [
    /* @__PURE__ */ jsx("p", { className: "text-sm font-medium uppercase tracking-wide text-[#8df7a4]", children: "404" }),
    /* @__PURE__ */ jsx("h1", { className: "mt-2 text-2xl font-semibold text-[#f2ead8]", children: "Not found" }),
    /* @__PURE__ */ jsx("p", { className: "mt-3 text-sm leading-6 text-[#8f8a7d]", children: "The page you are looking for could not be found." })
  ] }) });
}
const __vite_glob_0_8 = /* @__PURE__ */ Object.freeze(/* @__PURE__ */ Object.defineProperty({
  __proto__: null,
  default: NotFound
}, Symbol.toStringTag, { value: "Module" }));
const badgeVariants = cva(
  "group/badge inline-flex h-5 w-fit shrink-0 items-center justify-center gap-1 overflow-hidden rounded-none border border-transparent px-2 py-0.5 text-xs font-medium whitespace-nowrap transition-all focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50 has-data-[icon=inline-end]:pr-1.5 has-data-[icon=inline-start]:pl-1.5 aria-invalid:border-destructive aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 [&>svg]:pointer-events-none [&>svg]:size-3!",
  {
    variants: {
      variant: {
        default: "bg-primary text-primary-foreground [a]:hover:bg-primary/80",
        secondary: "bg-secondary text-secondary-foreground [a]:hover:bg-secondary/80",
        destructive: "bg-destructive/10 text-destructive focus-visible:ring-destructive/20 dark:bg-destructive/20 dark:focus-visible:ring-destructive/40 [a]:hover:bg-destructive/20",
        outline: "border-border text-foreground [a]:hover:bg-muted [a]:hover:text-muted-foreground",
        ghost: "hover:bg-muted hover:text-muted-foreground dark:hover:bg-muted/50",
        link: "text-primary underline-offset-4 hover:underline"
      }
    },
    defaultVariants: {
      variant: "default"
    }
  }
);
function Badge({
  className,
  variant = "default",
  render,
  ...props
}) {
  return useRender({
    defaultTagName: "span",
    props: mergeProps(
      {
        className: cn(badgeVariants({ variant }), className)
      },
      props
    ),
    render,
    state: {
      slot: "badge",
      variant
    }
  });
}
function Card({
  className,
  size = "default",
  ...props
}) {
  return /* @__PURE__ */ jsx(
    "div",
    {
      "data-slot": "card",
      "data-size": size,
      className: cn(
        "group/card flex flex-col gap-(--card-spacing) overflow-hidden rounded-none bg-card py-(--card-spacing) text-xs/relaxed text-card-foreground ring-1 ring-foreground/10 [--card-spacing:--spacing(4)] has-data-[slot=card-footer]:pb-0 has-[>img:first-child]:pt-0 data-[size=sm]:[--card-spacing:--spacing(3)] data-[size=sm]:has-data-[slot=card-footer]:pb-0 *:[img:first-child]:rounded-none *:[img:last-child]:rounded-none",
        className
      ),
      ...props
    }
  );
}
function CardHeader({ className, ...props }) {
  return /* @__PURE__ */ jsx(
    "div",
    {
      "data-slot": "card-header",
      className: cn(
        "group/card-header @container/card-header grid auto-rows-min items-start gap-1 rounded-none px-(--card-spacing) has-data-[slot=card-action]:grid-cols-[1fr_auto] has-data-[slot=card-description]:grid-rows-[auto_auto] [.border-b]:pb-(--card-spacing)",
        className
      ),
      ...props
    }
  );
}
function CardTitle({ className, ...props }) {
  return /* @__PURE__ */ jsx(
    "div",
    {
      "data-slot": "card-title",
      className: cn(
        "font-heading text-sm font-medium group-data-[size=sm]/card:text-sm",
        className
      ),
      ...props
    }
  );
}
function CardDescription({ className, ...props }) {
  return /* @__PURE__ */ jsx(
    "div",
    {
      "data-slot": "card-description",
      className: cn("text-xs/relaxed text-muted-foreground", className),
      ...props
    }
  );
}
function CardContent({ className, ...props }) {
  return /* @__PURE__ */ jsx(
    "div",
    {
      "data-slot": "card-content",
      className: cn("px-(--card-spacing)", className),
      ...props
    }
  );
}
const features = [
  {
    title: "Documentation",
    description: "Start with the framework guides and learn the conventions that shape an Andurel app.",
    href: "https://andurel.com"
  },
  {
    title: "Authentication",
    description: "Registration, sessions, email confirmation, and password reset are already wired.",
    href: routes.registrationNew()
  },
  {
    title: "Templ and Inertia",
    description: "Render server-side Templ pages or scaffold Inertia pages when Inertia is enabled.",
    href: routes.homePage()
  },
  {
    title: "Command line",
    description: "Generate models, factories, controllers, routes, jobs, emails, and views from the CLI.",
    href: "https://andurel.com"
  }
];
const consoleLines = [
  { label: "auth", value: "sessions, registration, reset password" },
  { label: "views", value: "Templ + Inertia scaffold online" },
  { label: "routes", value: "slugged endpoints mounted" },
  { label: "jobs", value: "queue worker ready" }
];
function Home() {
  return /* @__PURE__ */ jsx(Layout, { children: /* @__PURE__ */ jsx("div", { className: "mx-auto w-full max-w-[960px] px-4 py-4", children: /* @__PURE__ */ jsxs("section", { className: "grid items-center gap-6 py-6 lg:grid-cols-[minmax(0,1fr)_22rem]", children: [
    /* @__PURE__ */ jsxs("div", { className: "max-w-3xl", children: [
      /* @__PURE__ */ jsx("h1", { className: "text-4xl font-semibold text-card-foreground sm:text-5xl lg:text-6xl", children: "Space-grade Go, wired locally." }),
      /* @__PURE__ */ jsx("p", { className: "mt-5 max-w-2xl text-lg leading-7 text-muted-foreground", children: "Andurel has generated the core app shell: routing, controllers, validation, authentication, email, queues, Templ, and Inertia-ready frontends." }),
      /* @__PURE__ */ jsxs("div", { className: "mt-8 flex flex-wrap gap-3", children: [
        /* @__PURE__ */ jsx(Button, { render: /* @__PURE__ */ jsx("a", { href: routes.documentationShow("latest", "introduction") }), children: "Read the docs" }),
        /* @__PURE__ */ jsx(Button, { variant: "secondary", render: /* @__PURE__ */ jsx(Link, { href: routes.registrationNew() }), children: "Create ACCOUNT" })
      ] })
    ] }),
    /* @__PURE__ */ jsxs(Card, { className: "border-border shadow-sm shadow-black/40", children: [
      /* @__PURE__ */ jsx(CardHeader, { className: "border-b border-border pb-3", children: /* @__PURE__ */ jsxs(CardTitle, { className: "flex items-center justify-between text-xs uppercase text-text-muted", children: [
        /* @__PURE__ */ jsx("span", { children: "deploy console" }),
        /* @__PURE__ */ jsx(Badge, { variant: "outline", className: "border-accent text-accent", children: "ready" })
      ] }) }),
      /* @__PURE__ */ jsx(CardContent, { className: "space-y-4 font-mono text-sm", children: consoleLines.map((line) => /* @__PURE__ */ jsxs("div", { className: "flex items-start gap-3", children: [
        /* @__PURE__ */ jsx("span", { className: "mt-1 size-1.5 bg-accent" }),
        /* @__PURE__ */ jsxs("div", { children: [
          /* @__PURE__ */ jsx("p", { className: "text-accent", children: line.label }),
          /* @__PURE__ */ jsx("p", { className: "text-text-muted", children: line.value })
        ] })
      ] }, line.label)) })
    ] }),
    /* @__PURE__ */ jsx("section", { className: "grid gap-4 md:grid-cols-2 lg:col-span-2", children: features.map((feature) => /* @__PURE__ */ jsx(
      "a",
      {
        className: "group block border border-border transition hover:border-ring",
        href: feature.href,
        children: /* @__PURE__ */ jsx(Card, { className: "h-full border-0 ring-0", children: /* @__PURE__ */ jsxs(CardHeader, { children: [
          /* @__PURE__ */ jsx("div", { className: "mb-2 flex size-8 items-center justify-center border border-ring bg-background text-accent transition group-hover:border-accent", children: /* @__PURE__ */ jsx("span", { className: "size-2 bg-accent" }) }),
          /* @__PURE__ */ jsx(CardTitle, { className: "text-base", children: feature.title }),
          /* @__PURE__ */ jsx(CardDescription, { className: "text-sm leading-5", children: feature.description })
        ] }) })
      },
      feature.title
    )) })
  ] }) }) });
}
const __vite_glob_0_9 = /* @__PURE__ */ Object.freeze(/* @__PURE__ */ Object.defineProperty({
  __proto__: null,
  default: Home
}, Symbol.toStringTag, { value: "Module" }));
function isFlashMessage(value) {
  if (!value || typeof value !== "object") {
    return false;
  }
  const flash = value;
  return typeof flash.Type === "string" && typeof flash.Message === "string";
}
function pageFlashes(value) {
  if (!Array.isArray(value)) {
    return void 0;
  }
  return value.filter(isFlashMessage);
}
function FlashToasts({ initialFlashes }) {
  const [toasts, setToasts] = useState([]);
  useEffect(() => {
    let nextId = 0;
    function pushFlashes(flashes) {
      if (!flashes || flashes.length === 0) {
        return;
      }
      for (const flash of flashes) {
        const id = nextId++;
        setToasts((current) => [...current, { ...flash, id }]);
        window.setTimeout(() => {
          setToasts((current) => current.filter((toast) => toast.id !== id));
        }, 5e3);
      }
    }
    pushFlashes(initialFlashes);
    const removeListener = router.on("success", (event) => {
      pushFlashes(pageFlashes(event.detail.page.flash));
    });
    return () => {
      removeListener();
    };
  }, [initialFlashes]);
  if (toasts.length === 0) {
    return null;
  }
  return /* @__PURE__ */ jsx("div", { className: "fixed bottom-4 right-4 z-50", children: toasts.map((toast) => {
    const colorClass = toast.Type === "success" ? "border-[#8df7a4] text-[#8df7a4]" : toast.Type === "error" ? "border-[#ff875f] text-[#ff875f]" : "border-[#ff6b1a] text-[#e4dfd2]";
    return /* @__PURE__ */ jsx(
      "div",
      {
        className: `${colorClass} mb-2 border bg-[#101414] px-4 py-3 shadow-lg shadow-black/40 transition-opacity duration-300`,
        children: toast.Message
      },
      toast.id
    );
  }) });
}
const serverOptions = {
  host: process.env.INERTIA_SSR_HOST ?? "127.0.0.1",
  port: Number(process.env.INERTIA_SSR_PORT ?? "13714")
};
const renderPage = (page) => createInertiaApp({
  page,
  render: ReactDOMServer.renderToString,
  resolve: (name) => {
    const pages = /* @__PURE__ */ Object.assign({ "./Pages/Auth/ConfirmEmail.tsx": __vite_glob_0_0, "./Pages/Auth/Login.tsx": __vite_glob_0_1, "./Pages/Auth/Registration.tsx": __vite_glob_0_2, "./Pages/Auth/ResetPassword.tsx": __vite_glob_0_3, "./Pages/Auth/ResetPasswordRequest.tsx": __vite_glob_0_4, "./Pages/Documentation/Show.tsx": __vite_glob_0_5, "./Pages/Errors/BadRequest.tsx": __vite_glob_0_6, "./Pages/Errors/InternalError.tsx": __vite_glob_0_7, "./Pages/Errors/NotFound.tsx": __vite_glob_0_8, "./Pages/Home.tsx": __vite_glob_0_9 });
    return pages[`./Pages/${name}.tsx`].default;
  },
  setup: ({ App, props }) => /* @__PURE__ */ jsxs(Fragment, { children: [
    /* @__PURE__ */ jsx(App, { ...props }),
    /* @__PURE__ */ jsx(FlashToasts, { initialFlashes: pageFlashes(props.initialPage?.flash) })
  ] })
});
{
  createServer(
    renderPage,
    serverOptions
  );
}
export {
  renderPage as default
};
//# sourceMappingURL=ssr.js.map
