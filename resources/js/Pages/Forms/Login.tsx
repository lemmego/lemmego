import React from "react";

const Login: React.FC = () => {
	return (
		<div className="container w-1/3 mx-auto">
			<h1 className="text-3xl text-center">Login</h1>
			<form action="/login" method="POST">
				
					<div className="mt-2">
						<label htmlFor="email" className="label-primary">Email</label>
                        <input id="email" name="email" type="email" className="input"/>
					</div>
				
					<div className="mt-2">
						<label htmlFor="password" className="label-primary">Password</label>
                        <input id="password" name="password" type="password" className="input"/>
					</div>
				
					<div className="mt-2">
						<label htmlFor="org_username" className="label-primary">Org Username</label>
                        <input id="org_username" name="org_username" type="text" className="input"/>
					</div>
				
				<div>
					<button type="submit" className="mt-4 btn-primary">Submit</button>
				</div>
			</form>
		</div>
	);
};

export default Login;
