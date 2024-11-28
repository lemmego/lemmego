import React from "react";

const IndexReact = () => {
  return (
    <div className="relative flex min-h-screen flex-col justify-center overflow-hidden bg-gray-100 text-gray-600 dark:text-gray-400 dark:bg-gray-900 py-6 sm:py-12">
      <div className="relative bg-white dark:bg-gray-800 px-6 pt-10 pb-8 shadow-xl ring-1 ring-gray-900/5 sm:mx-auto sm:max-w-xl sm:rounded-lg sm:px-10">
        <div className="mx-auto">
          <div className="flex items-center justify-center space-x-6">
            <img
              src="https://avatars.githubusercontent.com/u/109903896?s=200&v=4"
              alt="Lemmego"
              className="w-20 h-20"
            />
          </div>
          <div className="divide-y divide-gray-300 dark:divide-gray-700">
            <div className="py-8 text-base leading-7">
              <p>
                <strong>Lemmego</strong> is a modern, full-stack web development
                framework built with Go, designed to streamline the creation of
                scalable and high-performance applications.
              </p>

              <h2 className="font-semibold mt-6 text-red-500 hover:text-red-600">
                Build Faster with Lemmego
              </h2>
              <p>
                A full-stack framework designed for developers to create
                scalable web applications effortlessly.
              </p>

              <h2 className="font-semibold mt-6 text-red-500 hover:text-red-600">
                Blazing Fast
              </h2>
              <p>Optimized for high performance and speed.</p>

              <h2 className="font-semibold mt-6 text-red-500 hover:text-red-600">
                Full-Stack Support
              </h2>
              <p>Integrated tools for both frontend and backend.</p>
              <h2 className="font-semibold mt-6 text-red-500 hover:text-red-600">
                Scalable
              </h2>
              <p>Perfect for projects of any size, big or small.</p>

              <hr className="my-5" />
              <p>
                Want more information on <strong>Lemmego</strong>
              </p>
              <a
                href="https://lemmego.github.io"
                target="_blank"
                rel="noopener noreferrer"
                className="font-semibold text-red-500 hover:text-red-600"
              >
                Read the docs →
              </a>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default IndexReact;
