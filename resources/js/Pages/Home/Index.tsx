import { usePage } from "@inertiajs/react";
import React from "react";

const Index: React.FC = () => {
  const { user } = usePage().props;
  return (
    <div>
      <h1>Welcome {user.username}</h1>
    </div>
  );
};

export default Index;
