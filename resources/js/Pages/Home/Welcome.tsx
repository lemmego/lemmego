import React from "react";

const Welcome: React.FC = (props: any) => {
  return (
    <div>
      <h1>This is the Welcome Page. This page is the welcome page {props.name}</h1>
    </div>
  );
};

export default Welcome;
