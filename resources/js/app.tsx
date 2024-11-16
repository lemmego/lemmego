import { createInertiaApp } from "@inertiajs/react";
import { createRoot } from "react-dom/client";
import React from "react";
import "../css/app.css";


type InertiaSetupArgs = {
  el: HTMLElement;
  App: React.ComponentType<any>;
  props: Record<string, unknown>;
};

createInertiaApp({
  resolve: (name: string) => {
    const pages = import.meta.glob("./Pages/**/*.tsx", { eager: true });
    return pages[`./Pages/${name}.tsx`];
  },
  setup({ el, App, props }: InertiaSetupArgs) {
    createRoot(el).render(<App {...props} />);
  },
});
