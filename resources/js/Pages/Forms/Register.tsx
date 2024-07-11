import React from "react";

const Register: React.FC = () => {
  return (
    <div className="container w-1/3 mx-auto">
      <h1 className="text-3xl text-center">Register</h1>
      <form action="/register" method="POST">
        <div className="mt-2">
          <label htmlFor="first_name" className="label-primary">
            First Name
          </label>
          <input
            id="first_name"
            name="first_name"
            type="text"
            className="input"
          />
        </div>

        <div className="mt-2">
          <label htmlFor="last_name" className="label-primary">
            Last Name
          </label>
          <input
            id="last_name"
            name="last_name"
            type="text"
            className="input"
          />
        </div>

        <div className="mt-2">
          <label htmlFor="email" className="label-primary">
            Email
          </label>
          <input id="email" name="email" type="email" className="input" />
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
          />
        </div>

        <div className="mt-2">
          <label htmlFor="confirm_password" className="label-primary">
            Confirm Password
          </label>
          <input
            id="confirm_password"
            name="confirm_password"
            type="password"
            className="input"
          />
        </div>

        <div className="mt-2">
          <label htmlFor="date_of_birth" className="label-primary">
            Date Of Birth
          </label>
          <input
            id="date_of_birth"
            name="date_of_birth"
            type="date"
            className="input"
          />
        </div>

        <div className="mt-2">
          <div className="flex items-center">
            <input id="male" name="gender" type="radio" value="male" />
            <label htmlFor="male" className="mx-2">
              Male
            </label>

            <input id="female" name="gender" type="radio" value="female" />
            <label htmlFor="female" className="mx-2">
              Female
            </label>

            <input id="other" name="gender" type="radio" value="other" />
            <label htmlFor="other" className="mx-2">
              Other
            </label>
          </div>
        </div>

        <div className="mt-2">
          <label htmlFor="favorite_language" className="label-primary">
            Favorite Language
          </label>
          <select
            id="favorite_language"
            name="favorite_language"
            className="input"
          >
            <option value="go" className="label-primary">
              Go
            </option>

            <option value="php" className="label-primary">
              PHP
            </option>

            <option value="java_script" className="label-primary">
              JavaScript
            </option>
          </select>
        </div>

        <div className="mt-2">
          <label htmlFor="propic" className="label-primary">
            Propic
          </label>
          <input id="propic" name="propic" type="file" />
        </div>

        <div className="mt-2">
          <label htmlFor="bio" className="label-primary">
            Bio
          </label>
          <textarea id="bio" name="bio" className="input"></textarea>
        </div>

        <div className="mt-2">
          <div className="flex items-center">
            <input
              id="i_agree_with_privacy_policy"
              name="i_agree_with_privacy_policy"
              type="checkbox"
            />
            <label htmlFor="i_agree_with_privacy_policy" className="mx-2">
              I Agree With Privacy Policy
            </label>
          </div>
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

export default Register;
