import React from "react";
import { useForm, usePage } from "@inertiajs/react";

const Register: React.FC = () => {
	const {errors, input} = usePage().props
  const { data, setData, post, progress } = useForm({
    first_name: "",
    last_name: "",
    logo: undefined,
    email: "",
    password: "",
    password_confirmation: "",
    org_name: "",
    org_email: "",
    org_username: "",
  });

  function handleInput(e: React.ChangeEvent<HTMLInputElement>) {
    if (e.target.files && e.target.files.length > 0) {
      setData(e.target.name, e.target.files[0]);
      return;
    }

    setData(e.target.name, e.target.value);
  }

  function submit(e: React.FormEvent<HTMLFormElement>) {
    e.preventDefault();
    post("/register");
  }

	return (
		<div className="w-1/3 mx-auto">
			<h1 className="text-3xl text-center">Register</h1>
			<form onSubmit={submit}>
				
					<div className="mt-2">
						<label htmlFor="first_name" className="label-primary">First Name</label>
                        <input id="first_name" name="first_name" type="text" className="input" value={data.first_name} onChange={handleInput}/>
						{errors.first_name && <p className="text-xs text-red-500">{errors.first_name.join(', ')}</p>}
					</div>
				
					<div className="mt-2">
						<label htmlFor="last_name" className="label-primary">Last Name</label>
                        <input id="last_name" name="last_name" type="text" className="input" value={data.last_name} onChange={handleInput}/>
						{errors.last_name && <p className="text-xs text-red-500">{errors.last_name.join(', ')}</p>}
					</div>
				
					<div className="mt-2">
						<label htmlFor="logo" className="label-primary">Logo</label>
						<input id="logo" name="logo" type="file" value={data.logo} onChange={handleInput}/>
						{errors.logo && <p className="text-xs text-red-500">{errors.logo.join(', ')}</p>}
					</div>
				
					<div className="mt-2">
						<label htmlFor="email" className="label-primary">Email</label>
                        <input id="email" name="email" type="email" className="input" value={data.email} onChange={handleInput}/>
						{errors.email && <p className="text-xs text-red-500">{errors.email.join(', ')}</p>}
					</div>
				
					<div className="mt-2">
						<label htmlFor="password" className="label-primary">Password</label>
                        <input id="password" name="password" type="password" className="input" value={data.password} onChange={handleInput}/>
						{errors.password && <p className="text-xs text-red-500">{errors.password.join(', ')}</p>}
					</div>
				
					<div className="mt-2">
						<label htmlFor="password_confirmation" className="label-primary">Password Confirmation</label>
                        <input id="password_confirmation" name="password_confirmation" type="password" className="input" value={data.password_confirmation} onChange={handleInput}/>
						{errors.password_confirmation && <p className="text-xs text-red-500">{errors.password_confirmation.join(', ')}</p>}
					</div>
				
					<div className="mt-2">
						<label htmlFor="org_name" className="label-primary">Org Name</label>
                        <input id="org_name" name="org_name" type="text" className="input" value={data.org_name} onChange={handleInput}/>
						{errors.org_name && <p className="text-xs text-red-500">{errors.org_name.join(', ')}</p>}
					</div>
				
					<div className="mt-2">
						<label htmlFor="org_email" className="label-primary">Org Email</label>
                        <input id="org_email" name="org_email" type="email" className="input" value={data.org_email} onChange={handleInput}/>
						{errors.org_email && <p className="text-xs text-red-500">{errors.org_email.join(', ')}</p>}
					</div>
				
					<div className="mt-2">
						<label htmlFor="org_username" className="label-primary">Org Username</label>
                        <input id="org_username" name="org_username" type="text" className="input" value={data.org_username} onChange={handleInput}/>
						{errors.org_username && <p className="text-xs text-red-500">{errors.org_username.join(', ')}</p>}
					</div>
				
				<div>
					<button type="submit" className="mt-4 btn-primary">Submit</button>
				</div>
			</form>
		</div>
	);
};

export default Register;
