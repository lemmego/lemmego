import React from "react";

const Register: React.FC = () => {
	return (
		<div className="container w-1/3 mx-auto">
			<h1 className="text-3xl text-center">Register</h1>
			<form action="/register" method="POST">
				
					<div className="mt-2">
						<label htmlFor="first_name" className="label-primary">First Name</label>
                        <input id="first_name" name="first_name" type="text" className="input"/>
					</div>
				
					<div className="mt-2">
						<label htmlFor="last_name" className="label-primary">Last Name</label>
                        <input id="last_name" name="last_name" type="text" className="input"/>
					</div>
				
					<div className="mt-2">
						<label htmlFor="email" className="label-primary">Email</label>
                        <input id="email" name="email" type="email" className="input"/>
					</div>
				
					<div className="mt-2">
						<label htmlFor="password" className="label-primary">Password</label>
                        <input id="password" name="password" type="password" className="input"/>
					</div>
				
					<div className="mt-2">
						<label htmlFor="password_confirmation" className="label-primary">Password Confirmation</label>
                        <input id="password_confirmation" name="password_confirmation" type="password" className="input"/>
					</div>
				
					<div className="mt-2">
						<label htmlFor="org_username" className="label-primary">Org Username</label>
                        <input id="org_username" name="org_username" type="text" className="input"/>
					</div>
				
					<div className="mt-2">
						<label htmlFor="org_name" className="label-primary">Org Name</label>
                        <input id="org_name" name="org_name" type="text" className="input"/>
					</div>
				
				<div>
					<button type="submit" className="mt-4 btn-primary">Submit</button>
				</div>
			</form>
		</div>
	);
};

export default Register;
