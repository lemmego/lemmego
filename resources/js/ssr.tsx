import { createInertiaApp } from "@inertiajs/react";
import createServer from "@inertiajs/react/server";
import ReactDOMServer from "react-dom/server";
import React from "react";

type InertiaSetupArgs = {
  App: React.ComponentType<any>;
  props: Record<string, unknown>;
};

createServer((page) =>
  createInertiaApp({
    page,
    render: ReactDOMServer.renderToString,
    resolve: (name: string) => {
      const pages = import.meta.glob("./Pages/**/*.tsx", { eager: true });
      return pages[`./Pages/${name}.tsx`];
    },
    setup: ({ App, props }: InertiaSetupArgs) => <App {...props} />,
  }),
);
