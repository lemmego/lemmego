import React, {useState} from "react";
import { usePage, router } from "@inertiajs/react";

const OauthClient: React.FC = () => {
	const {errors, input} = usePage().props
	const [values, setValues] = useState({
		name: "",
		redirect_uri: "",
	})
	function handleChange(e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) {
		const key = e.target.id;
		const value = e.target.value
		setValues(values => ({
			...values,
			[key]: value,
		}))
	}
	function handleSubmit(e: React.FormEvent<HTMLFormElement>) {
		e.preventDefault()
		router.post("/oauth/clients", values)
	}
	return (
		<div className="w-1/3 mx-auto">
			<h1 className="text-3xl text-center">Oauth Client</h1>
			<form onSubmit={handleSubmit}>
				
					<div className="mt-2">
						<label htmlFor="name" className="label-primary">Name</label>
                        <input id="name" name="name" type="text" className="input" onChange={handleChange} value={values.name} />
						{errors.name && <p className="text-xs text-red-500">{errors.name.join(', ')}</p>}
					</div>
				
					<div className="mt-2">
						<label htmlFor="redirect_uri" className="label-primary">Redirect Uri</label>
						<textarea id="redirect_uri" name="redirect_uri" className="input" onChange={handleChange} value={values.redirect_uri}></textarea>
						{errors.redirect_uri && <p className="text-xs text-red-500">{errors.redirect_uri.join(', ')}</p>}
					</div>
				
				<div>
					<button type="submit" className="mt-4 btn-primary">Submit</button>
				</div>
			</form>
		</div>
	);
};

export default OauthClient;
