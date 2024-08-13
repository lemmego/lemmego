import React from "react";
import { usePage } from "@inertiajs/react";

const OauthClient: React.FC = () => {
  const { errors } = usePage().props;
  return (
    <div className="w-1/3 mx-auto">
      <h1 className="text-3xl text-center">Authorize Application Request</h1>
      <p className="text-center text-gray-500 text-sm">
        An application is trying to access your account
      </p>
      <section className="flex justify-center space-x-3">
        <form action="/oauth/authorize/allow" method="POST">
          <div>
            <button type="submit" className="mt-4 btn-success">
              Allow
            </button>
          </div>
        </form>
        <form action="/oauth/authorize/deny" method="POST">
          <div>
            <button type="submit" className="mt-4 btn btn-danger">
              Deny
            </button>
          </div>
        </form>
      </section>
    </div>
  );
};

export default OauthClient;
